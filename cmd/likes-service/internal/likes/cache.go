package likes

import "sort"

func (s *service) updateLikes(uuid string, count int64) int64 {

	s.mutex.Lock()
	totalLikes := s.postLikes[uuid] + count
	s.postLikes[uuid] = totalLikes
	s.mutex.Unlock()
	return totalLikes
}

func (s *service) getTopLikes(limit int32) []likesCount {
	likes := make([]likesCount, 0, limit)
	s.mutex.RLock()
	for postID, numLikes := range s.postLikes {
		likes = append(likes, likesCount{postID: postID, numLikes: numLikes})
	}
	s.mutex.RUnlock()

	sort.Slice(likes, func(i, j int) bool {
		// Do this in reverse so that order is in descending
		return likes[i].numLikes > likes[j].numLikes
	})
	return likes[:limit]
}
