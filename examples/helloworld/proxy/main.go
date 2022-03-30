package main

import (
	"context"
	"log"
	"net"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)))

	log.Println("proxy start listening...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func director(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
	log.Printf("Proxy Request for method: %v\n", fullMethodName)
	clientConn, err := grpc.DialContext(
		ctx,
		"localhost:50051",
		grpc.WithCodec(proxy.Codec()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	return ctx, clientConn, err
}
