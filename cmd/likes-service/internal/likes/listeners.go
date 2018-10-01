package likes

import (
	"context"

	"github.com/google/uuid"
)

func newPublisher(queueSize int) *updatesPublisher {
	return &updatesPublisher{
		updateNotification:  make(chan likesCount, queueSize),
		newSubscriptions:    make(chan subscriber),
		closedSubscriptions: make(chan subscriber),
		done:                make(chan struct{}),
	}
}

type updatesPublisher struct {
	updateNotification  chan likesCount
	newSubscriptions    chan subscriber
	closedSubscriptions chan subscriber
	done                chan struct{}
}

type subscriber struct {
	uuid uuid.UUID
	c    chan likesCount
}

func (u *updatesPublisher) listen() {

	subscribers := make(map[uuid.UUID]subscriber)
	for {
		select {
		case likes := <-u.updateNotification:
			publish(subscribers, likes)
		case s := <-u.newSubscriptions:
			subscribers[s.uuid] = s
		case s := <-u.closedSubscriptions:
			delete(subscribers, s.uuid)
			close(s.c)
		case <-u.done:
			cleanup(subscribers)
			return
		}
	}
}

func publish(subscribers map[uuid.UUID]subscriber, likes likesCount) {

	for _, s := range subscribers {
		s.c <- likes
	}
}

func cleanup(subscribers map[uuid.UUID]subscriber) {

	for _, s := range subscribers {
		close(s.c)
	}
}

func (u *updatesPublisher) close() {

	close(u.done)
}

func (u *updatesPublisher) notify(ctx context.Context, postID string, totalLikes int64) bool {

	select {
	case <-ctx.Done():
		return false
	case u.updateNotification <- likesCount{postID: postID, numLikes: totalLikes}:
		return true
	}
}

func (u *updatesPublisher) subscribe(ctx context.Context, onUpdate func(l likesCount)) error {

	s := subscriber{
		uuid: uuid.New(),
		c:    make(chan likesCount),
	}

	u.newSubscriptions <- s
	for {

		select {
		case likes := <-s.c:
			onUpdate(likes)
		case <-ctx.Done():
			u.closedSubscriptions <- s
			return ctx.Err()
		}
	}
}
