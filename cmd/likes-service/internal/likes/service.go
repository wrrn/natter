package likes

import (
	"context"

	pb "github.com/wrrn/natter/pkg/likes"
	pb_post "github.com/wrrn/natter/pkg/post"
)

func NewService(postsClient pb_post.PostServiceClient) pb.LikesServiceServer {
	s := &service{
		postsClient: postsClient,
		cache:       newCache(),
	}

	return s
}

type service struct {
	cache       *cache
	postsClient pb_post.PostServiceClient
}

type likesCount struct {
	postID   string
	numLikes int64
}

func (s *service) UpdateLikes(ctx context.Context, req *pb.UpdateLikesRequest) (*pb.UpdateLikesResponse, error) {
	totalLikes := s.cache.update(ctx, req.Uuid, req.Count)
	return &pb.UpdateLikesResponse{Uuid: req.Uuid, TotalLikes: totalLikes}, nil
}

func (s *service) GetTrending(ctx context.Context, req *pb.GetTrendingRequest) (*pb.GetTrendingResponse, error) {
	limit := req.Limit
	if limit > int32(len(s.cache.likes)) {
		limit = int32(len(s.cache.likes))
	}

	topLikes := s.cache.getTopLikes(limit)
	resp, err := s.postsClient.BatchGetPosts(ctx, &pb_post.BatchGetPostsRequest{Uuids: topLikes.IDs()})
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
	ctx := stream.Context()
	err := s.cache.onUpdate(ctx, func(string, int64) error {
		resp, err := s.GetTrending(ctx, req)
		if err != nil {
			return err
		}
		return stream.Send(resp)
	})

	return err
}
