package ws

import (
	"context"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slices"
)

type subscriber struct {
	conn    *Conn
	channel string
	pubsub  *redis.PubSub
}

var subscribers []*subscriber

func unsub(conn *Conn, channel string) {
	for i, s := range subscribers {
		if s == nil {
			continue
		}

		if s.channel == channel && s.conn == conn {
			s.pubsub.Unsubscribe(context.Background(), channel)
			s.pubsub.Close()

			subscribers = slices.Delete(subscribers, i, i+1)
		}
	}
}

func removeAnySubscriber(conn *Conn) {
	for i, s := range subscribers {
		if s == nil {
			continue
		}

		if s.conn != nil && s.conn == conn {
			subscribers[i] = nil
			subscribers = slices.Delete(subscribers, i, i+1)
		}
	}
}
