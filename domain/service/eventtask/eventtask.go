package eventtask

import (
	"github.com/utahta/momoclo-channel/domain/event"
	"github.com/utahta/momoclo-channel/domain/model"
)

// NewEnqueueTweets returns enqueue tweets task
func NewEnqueueTweets(v model.FeedItem) event.Task {
	return event.Task{QueueName: "enqueue", Path: "/enqueue/tweets", Object: v}
}

// NewEnqueueLines returns enqueue lines task
func NewEnqueueLines(v model.FeedItem) event.Task {
	return event.Task{QueueName: "enqueue", Path: "/enqueue/lines", Object: v}
}

// NewTweet returns tweet task
func NewTweet(v model.TweetRequest) event.Task {
	return NewTweets([]model.TweetRequest{v})
}

// NewTweets returns tweet task
func NewTweets(v []model.TweetRequest) event.Task {
	return event.Task{QueueName: "queue-tweet", Path: "/queue/tweet", Object: v}
}

// NewLine returns line task
func NewLine(v model.LineNotifyRequest) event.Task {
	return NewLines([]model.LineNotifyRequest{v})
}

// NewLines returns line task
func NewLines(v []model.LineNotifyRequest) event.Task {
	return event.Task{QueueName: "queue-line", Path: "/queue/line", Object: v}
}
