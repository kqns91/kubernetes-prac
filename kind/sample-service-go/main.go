package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/kqns91/kubernetes-prac/kind/sample-proto/gen/go/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

func main() {
	resolver.SetDefaultScheme("dns")
	conn, err := grpc.Dial("sample-service-grpc:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloWorldClient(conn)

	r := gin.Default()
	g := r.Group("/go")

	g.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	g.GET("/sample", func(ctx *gin.Context) {
		name := ctx.heQuery("name")

		response, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, response)
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
