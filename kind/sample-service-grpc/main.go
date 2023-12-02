package main

import (
	"context"
	"fmt"
	"log"
	"math"
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
	if in.Name == "scale" {
		go func(ctx context.Context) {
			x := 246.0
			defer func() {
				fmt.Printf("x: %v\n", x)
			}()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				x += math.Sqrt(x)
			}
		}(ctx)
		time.Sleep(7 * time.Second)
	}

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
			MaxConnectionAge: 60 * time.Second,
		}),
	)
	pb.RegisterHelloWorldServer(s, h)
	s.Serve(lis)
}
