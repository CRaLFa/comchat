// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var command_chat_pb = require('./command_chat_pb.js');

function serialize_comchat_ChatMessage(arg) {
  if (!(arg instanceof command_chat_pb.ChatMessage)) {
    throw new Error('Expected argument of type comchat.ChatMessage');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_comchat_ChatMessage(buffer_arg) {
  return command_chat_pb.ChatMessage.deserializeBinary(new Uint8Array(buffer_arg));
}


var CommandChatService = exports.CommandChatService = {
  chat: {
    path: '/comchat.CommandChat/Chat',
    requestStream: true,
    responseStream: true,
    requestType: command_chat_pb.ChatMessage,
    responseType: command_chat_pb.ChatMessage,
    requestSerialize: serialize_comchat_ChatMessage,
    requestDeserialize: deserialize_comchat_ChatMessage,
    responseSerialize: serialize_comchat_ChatMessage,
    responseDeserialize: deserialize_comchat_ChatMessage,
  },
};

exports.CommandChatClient = grpc.makeGenericClientConstructor(CommandChatService);
