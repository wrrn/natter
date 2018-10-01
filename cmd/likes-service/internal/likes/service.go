package likes

import (
	"context"
	"log"
	"sync"

	pb "github.com/wrrn/natter/pkg/likes"
	pb_post "github.com/wrrn/natter/pkg/post"
)

func NewService(postsClient pb_post.PostServiceClient) pb.LikesServiceServer {
	s := &service{
		postLikes:   make(map[string]int64),
		publisher:   newPublisher(100),
		postsClient: postsClient,
		mutex:       &sync.RWMutex{},
	}
	go s.publisher.listen()
	return s

}

type service struct {
	postLikes   map[string]int64
	postsClient pb_post.PostServiceClient

	publisher *updatesPublisher
	mutex     *sync.RWMutex
}

type likesCount struct {
	postID   string
	numLikes int64
}

func (s *service) UpdateLikes(ctx context.Context, req *pb.UpdateLikesRequest) (*pb.UpdateLikesResponse, error) {
	log.Println("Updating the likes")
	defer log.Println("Likes Updated")
	totalLikes := s.updateLikes(req.Uuid, req.Count)
	s.publisher.notify(ctx, req.Uuid, totalLikes)
	return &pb.UpdateLikesResponse{Uuid: req.Uuid, TotalLikes: totalLikes}, nil
}

func (s *service) GetTrending(ctx context.Context, req *pb.GetTrendingRequest) (*pb.GetTrendingResponse, error) {
	limit := req.Limit
	if limit > int32(len(s.postLikes)) {
		limit = int32(len(s.postLikes))
	}

	topLikes := s.getTopLikes(limit)
	postIDs := make([]string, 0, limit)
	for _, l := range topLikes {
		postIDs = append(postIDs, l.postID)
	}

	resp, err := s.postsClient.BatchGetPosts(ctx, &pb_post.BatchGetPostsRequest{Uuids: postIDs})
	if err != nil {
		return nil, err
	}

	trendingPosts := make([]*pb.TrendingPost, 0, req.Limit)
	for i, post := range resp.Posts {
		trendingPosts = append(trendingPosts, &pb.TrendingPost{
			Post:     post,
			NumLikes: topLikes[i].numLikes,
		})
	}

	return &pb.GetTrendingResponse{Posts: trendingPosts}, nil
}

func (s *service) StreamTrending(req *pb.GetTrendingRequest, stream pb.LikesService_StreamTrendingServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	var err error

	s.publisher.subscribe(ctx, func(l likesCount) {
		var resp *pb.GetTrendingResponse
		resp, err = s.GetTrending(ctx, req)
		if err != nil {
			cancel()
		}
		stream.Send(resp)
	})

	return err
}
