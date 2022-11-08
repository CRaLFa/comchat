package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "github.com/CRaLFa/comchat"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type commandChatServer struct {
	pb.UnimplementedCommandChatServer

	mu       sync.Mutex
	clients  map[string]pb.CommandChat_ChatServer
	msgQueue []*pb.ChatMessage
}

func (s *commandChatServer) Chat(stream pb.CommandChat_ChatServer) error {
	errCh := make(chan error)

	go s.receive(stream, errCh)
	go s.send(stream, errCh)

	return <-errCh
}

func (s *commandChatServer) receive(stream pb.CommandChat_ChatServer, errCh chan error) {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Printf("Failed to receive message: %v", err)
			errCh <- err
			continue
		}
		if msg.Body == "LOGGED_IN" || msg.Body == "LOGGED_OUT" {
			name := msg.Author
			msg.Author = "SYSTEM"

			var format string
			s.mu.Lock()
			if msg.Body == "LOGGED_IN" {
				format = "%s has entered."
				s.clients[name] = stream
			} else {
				format = "%s has exited."
				delete(s.clients, name)
			}
			s.mu.Unlock()
			msg.Body = fmt.Sprintf(format, name)
		} else {
			log.Printf("Received message: {%v}", msg)
		}

		s.mu.Lock()
		s.msgQueue = append(s.msgQueue, msg)
		s.mu.Unlock()
	}
}

func (s *commandChatServer) send(stream pb.CommandChat_ChatServer, errCh chan error) {
	for {
		time.Sleep(100 * time.Millisecond)

		s.mu.Lock()
		msgLen := len(s.msgQueue)
		if msgLen == 0 {
			s.mu.Unlock()
			continue
		}
		mq := make([]*pb.ChatMessage, msgLen)
		copy(mq, s.msgQueue)
		s.mu.Unlock()

		for _, msg := range mq {
			for _, cs := range s.clients {
				if err := cs.Send(msg); err != nil {
					log.Printf("Failed to send message: %v", err)
					errCh <- err
					continue
				}
			}
			log.Printf("Sent message: {%v}", msg)
		}

		s.mu.Lock()
		s.msgQueue = []*pb.ChatMessage{}
		s.mu.Unlock()
	}
}

func newServer() *commandChatServer {
	return &commandChatServer{
		clients:  make(map[string]pb.CommandChat_ChatServer),
		msgQueue: []*pb.ChatMessage{},
	}
}

func main() {
	address := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server running at %s", address)

	grpcServer := grpc.NewServer()
	pb.RegisterCommandChatServer(grpcServer, newServer())
	grpcServer.Serve(listener)
}
