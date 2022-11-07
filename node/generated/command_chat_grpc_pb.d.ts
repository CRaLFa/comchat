// GENERATED CODE -- DO NOT EDIT!

// package: comchat
// file: command_chat.proto

import * as command_chat_pb from "./command_chat_pb";
import * as grpc from "@grpc/grpc-js";

interface ICommandChatService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
  chat: grpc.MethodDefinition<command_chat_pb.ChatMessage, command_chat_pb.ChatMessage>;
}

export const CommandChatService: ICommandChatService;

export interface ICommandChatServer extends grpc.UntypedServiceImplementation {
  chat: grpc.handleBidiStreamingCall<command_chat_pb.ChatMessage, command_chat_pb.ChatMessage>;
}

export class CommandChatClient extends grpc.Client {
  constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
  chat(metadataOrOptions?: grpc.Metadata | grpc.CallOptions | null): grpc.ClientDuplexStream<command_chat_pb.ChatMessage, command_chat_pb.ChatMessage>;
  chat(metadata?: grpc.Metadata | null, options?: grpc.CallOptions | null): grpc.ClientDuplexStream<command_chat_pb.ChatMessage, command_chat_pb.ChatMessage>;
}
