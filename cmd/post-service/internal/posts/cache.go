package posts

import pb "github.com/wrrn/natter/pkg/post"

func (s *Service) addPost(p pb.Post) {

	s.mutex.Lock()
	s.posts[p.Uuid] = p
	s.mutex.Unlock()
}

func (s *Service) getPost(id string) (pb.Post, bool) {

	s.mutex.RLock()
	p, found := s.posts[id]
	s.mutex.RUnlock()
	return p, found
}

func (s *Service) getPosts() []*pb.Post {

	posts := make([]*pb.Post, 0, len(s.posts))
	s.mutex.RLock()
	for _, post := range s.posts {
		posts = append(posts, &post)
	}
	s.mutex.RUnlock()
	return posts
}
