package publisher

import (
	"context"

	"github.com/google/uuid"
)

type UpdatesPublisher struct {
	updateNotification  chan Likes
	newSubscriptions    chan subscriber
	closedSubscriptions chan subscriber
	done                chan struct{}
}

type subscriber struct {
	uuid uuid.UUID
	c    chan Likes
}

func (u *UpdatesPublisher) Listen(done <-chan struct{}) {
	subscribers := make(map[uuid.UUID]subscriber)
	for {
		select {
		case likes := <-u.updateNotification:
			u.publish(likes)
		case s := <-u.newSubscriptions:
			subscribers[s.uuid] = s
		case s := <-u.closedSubsriptions:
			delete(subscribers[s.uuid])
			close(s.c)
		case <-u.done:
			cleanup(subscribers)
			return
		}
	}
}

func publish(subscribers map[uuid.UUID]subscriber, likes Likes) {

	for _, s := range subscribers {
		s.c <- likes
	}
}

func cleanup(subscribers map[uuid.UUID]subcriber) {

	for _, s := range subscribers {
		close(s.c)
	}
}

func (u *UpdatesPublisher) Close() {
	close(u.done)
}

func (u *UpdatesPublisher) Notify(ctx context.Context, postID string, totalLikes int64) bool {
	select {
	case <-ctx.Done():
		return false
	case u.updateNotification <- Likes{postID: postID, likes: totalLikes}:
		return true
	}
}

func (u *UpdatesPublisher) Subscribe(ctx context.Context, onUpdate func(l Likes)) {
	s := subscriber{
		uuid: uuid.New(),
		c:    make(chan Likes),
	}

	u.newSubscriptions <- s

	for {
		select {
		case likes := <-s.c:
			onUpdate(likes)
		case <-ctx.Done():
			u.closedSubscriptions <- s
			return
		}
	}
}
