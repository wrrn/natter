package main

import (
	"context"
	"flag"
	"log"

	"github.com/wrrn/natter/pkg/likes"
	"google.golang.org/grpc"
)

func main() {
	var (
		likesServiceAddr = flag.String("likes-addr", "127.0.0.1:9001", "address of the posts service")
		limit            = flag.Int("limit", 10, "how many of the top posts that you want to see")
	)
	flag.Parse()

	likesConn, err := grpc.Dial(*likesServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial the like service client: %v", err)
	}
	likesClient := likes.NewLikesServiceClient(likesConn)

	stream, err := likesClient.StreamTrending(context.Background(), &likes.GetTrendingRequest{Limit: int32(*limit)})
	if err != nil {
		log.Fatalf("Failed to stream trending: %v", err)
	}

	var printer printer
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("Receiving the trending stream failed: %v", err)
			return
		}
		printer.printPosts(resp.Posts)

	}

}
