package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/kqns91/kubernetes-prac/kind/sample-proto/gen/go/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type handler struct {
	pb.UnimplementedHelloWorldServer
}

func (h *handler) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	time.Sleep(15 * time.Second)
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s", in.Name)}, nil
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Received request: %v", info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		log.Printf("Error handling request: %v", err)
	}

	return resp, err
}

func main() {
	h := &handler{}
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 10 * time.Second,
		}))
	pb.RegisterHelloWorldServer(s, h)
	s.Serve(lis)
}
