/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	clientStreamDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqOrigin := &pb.HelloRequest{Name: *name}
	payload, err := proto.Marshal(reqOrigin)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ********* 核心代理处理逻辑 ********

	clientStream, err := grpc.NewClientStream(ctx, clientStreamDescForProxying, conn, "/helloworld.Greeter/SayHello")
	if err != nil {
		log.Fatalf("get stream error: %v", err)
	}

	// 发送消息
	req := &anypb.Any{}
	if err := proto.Unmarshal(payload, req); err != nil {
		fmt.Println(err)
		return
	}
	if err := clientStream.SendMsg(req); err != nil {
		log.Fatalf("send message error: %v", err)
	}

	// 接收消息
	reply := &anypb.Any{}
	if err := clientStream.RecvMsg(reply); err != nil {
		log.Fatalf("receive message error: %v", err)
	}
	replyPayload, err := proto.Marshal(reply)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ********* 核心代理处理逻辑 ********

	res := &pb.HelloReply{}

	if err := proto.Unmarshal(replyPayload, res); err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("Greeting: %v", res.GetMessage())
}
