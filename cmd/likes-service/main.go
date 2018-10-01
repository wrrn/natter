package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/wrrn/natter/cmd/likes-service/internal/likes"
	pb "github.com/wrrn/natter/pkg/likes"
	"github.com/wrrn/natter/pkg/post"
	"google.golang.org/grpc"
)

func main() {
	var (
		port             = flag.Int("port", 9001, "port the gRCP service will listen on")
		postsServiceAddr = flag.String("posts-addr", "127.0.0.1:9000", "address of the posts service")
	)
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	conn, err := grpc.Dial(*postsServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial the post service client: %v", err)
	}
	service := likes.NewService(post.NewPostServiceClient(conn))

	grpcServer := grpc.NewServer()
	pb.RegisterLikesServiceServer(grpcServer, service)
	grpcServer.Serve(lis)
}
