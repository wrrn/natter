package likes

import (
	"context"
	"sort"
	"sync"
)

func newCache() *cache {
	c := &cache{
		likes:     make(map[string]int64),
		publisher: newPublisher(100),
		RWMutex:   &sync.RWMutex{},
	}
	go c.publisher.listen()
	return c
}

type cache struct {
	likes     map[string]int64
	publisher *updatesPublisher
	*sync.RWMutex
}

func (c *cache) update(ctx context.Context, uuid string, count int64) int64 {

	c.Lock()
	totalLikes := c.likes[uuid] + count
	c.likes[uuid] = totalLikes
	c.Unlock()
	c.publisher.notify(ctx, uuid, totalLikes)
	return totalLikes
}

func (c *cache) onUpdate(ctx context.Context, onUpdate func(uuid string, totalLikes int64) error) error {
	return c.publisher.subscribe(ctx, onUpdate)
}

func (c *cache) getTopLikes(limit int32) likesCountList {
	likes := make([]likesCount, 0, limit)
	c.RLock()
	for postID, numLikes := range c.likes {
		likes = append(likes, likesCount{postID: postID, numLikes: numLikes})
	}
	c.RUnlock()

	sort.Slice(likes, func(i, j int) bool {
		// Do this in reverse so that order is in descending
		return likes[i].numLikes > likes[j].numLikes
	})
	return likes[:limit]
}

type likesCountList []likesCount

func (ll likesCountList) IDs() []string {
	ids := make([]string, 0, len(ll))
	for _, l := range ll {
		ids = append(ids, l.postID)
	}
	return ids
}
