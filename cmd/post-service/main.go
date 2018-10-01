package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/wrrn/natter/cmd/post-service/internal/posts"
	pb "github.com/wrrn/natter/pkg/post"
	"google.golang.org/grpc"
)

func main() {
	var (
		port = flag.Int("port", 9000, "port the gRCP service will listen on")
	)
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPostServiceServer(grpcServer, posts.NewService())
	grpcServer.Serve(lis)
}
