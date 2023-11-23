package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/kqns91/kubernetes-prac/kind/sample-proto/gen/go/protobuf"
	"google.golang.org/grpc"
)

type handler struct {
	pb.UnimplementedHelloWorldServer
}

func (h *handler) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s", in.Name)}, nil
}

func main() {
	h := &handler{}
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloWorldServer(s, h)

	s.Serve(lis)
}
