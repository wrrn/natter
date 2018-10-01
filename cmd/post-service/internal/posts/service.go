package posts

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
	pb "github.com/wrrn/natter/pkg/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewService() pb.PostServiceServer {
	return &Service{
		mutex: &sync.RWMutex{},
		posts: dadJokes(),
	}
}

// The service the implements the post service
type Service struct {
	mutex *sync.RWMutex
	posts map[string]pb.Post
}

func (s *Service) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {

	p := pb.Post{
		Uuid: uuid.New().String(),
		Msg:  req.Msg,
	}
	s.addPost(p)
	return &pb.CreatePostResponse{Post: &p}, nil
}

func (s *Service) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	log.Println("Getting Posts")
	return &pb.ListPostsResponse{Posts: s.getPosts()}, nil
}

func (s *Service) BatchGetPosts(ctx context.Context, req *pb.BatchGetPostsRequest) (*pb.BatchGetPostsResponse, error) {
	posts := make([]*pb.Post, 0, len(req.Uuids))
	for _, id := range req.Uuids {
		post, found := s.getPost(id)
		if !found {
			return nil, status.Errorf(codes.NotFound, "post not found")
		}

		posts = append(posts, &post)
	}

	return &pb.BatchGetPostsResponse{Posts: posts}, nil
}
