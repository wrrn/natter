package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/wrrn/natter/pkg/likes"
	"github.com/wrrn/natter/pkg/post"
	"google.golang.org/grpc"
)

func main() {
	var (
		likesServiceAddr = flag.String("likes-addr", "127.0.0.1:9001", "address of the posts service")
		postsServiceAddr = flag.String("posts-addr", "127.0.0.1:9000", "address of the posts service")
	)
	flag.Parse()

	postsConn, err := grpc.Dial(*postsServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial the post service client: %v", err)
	}
	postsClient := post.NewPostServiceClient(postsConn)

	likesConn, err := grpc.Dial(*likesServiceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial the like service client: %v", err)
	}
	likesClient := likes.NewLikesServiceClient(likesConn)

	rand.Seed(time.Now().UnixNano())
	for {
		resp, err := postsClient.ListPosts(context.Background(), &post.ListPostsRequest{})
		if err != nil {
			log.Println("Error getting the posts")
		}
		postID := resp.Posts[rand.Intn(len(resp.Posts))].Uuid
		likesClient.UpdateLikes(
			context.Background(),
			&likes.UpdateLikesRequest{Uuid: postID, Count: rand.Int63n(30) - 10},
		)
		time.Sleep(time.Second / 10)
	}
}
