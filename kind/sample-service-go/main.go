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
	// gin でサーバー立てる
	r := gin.Default()

	// sample-service-grpc に接続
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

	g := r.Group("/go")

	// アクセスエンドポイントを作成
	g.GET("/sample", func(ctx *gin.Context) {
		response, err := c.SayHello(ctx, &pb.HelloRequest{Name: "kqns91"})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, response)
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
