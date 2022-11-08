package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/CRaLFa/comchat"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address")
)

type commandChatClient struct {
	pb.CommandChatClient

	user string
}

func (c *commandChatClient) runChat() {
	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatalf("Failed to chat: %v", err)
	}
	defer stream.CloseSend()

	waitCh := make(chan bool)

	go c.receive(stream, waitCh)
	go c.send(stream, waitCh)
	c.sendMessage(stream, "LOGGED_IN")

	<-waitCh
}

func (c *commandChatClient) receive(stream pb.CommandChat_ChatClient, ch chan bool) {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Failed to receive message: %v", err)
			continue
		}
		log.Printf("[%s] : %s\n", msg.Author, msg.Body)
	}

	ch <- true
}

func (c *commandChatClient) send(stream pb.CommandChat_ChatClient, ch chan bool) {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Failed to read: %v", err)
		}
		trimmed := strings.TrimSpace(text)
		if len(trimmed) == 0 {
			continue
		}
		if trimmed == "exit" {
			c.sendMessage(stream, "LOGGED_OUT")
			break
		}
		c.sendMessage(stream, trimmed)
	}

	ch <- true
}

func (c *commandChatClient) sendMessage(stream pb.CommandChat_ChatClient, body string) {
	msg := &pb.ChatMessage{
		Author: c.user,
		Body:   body,
	}
	if err := stream.Send(msg); err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func getUserName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter your name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read: %v", err)
		return "(anonymous)"
	}
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return "(anonymous)"
	}
	return trimmed
}

func newClient(cc grpc.ClientConnInterface, userName string) *commandChatClient {
	return &commandChatClient{
		CommandChatClient: pb.NewCommandChatClient(cc),
		user:              userName,
	}
}

func doChat() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := newClient(conn, getUserName())
	client.runChat()
}

func main() {
	doChat()
	os.Exit(0)
}
