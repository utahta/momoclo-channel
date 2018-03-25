package eventtask

import (
	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/event"
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
	return event.Task{QueueName: "queue-tweet", Path: "/tweet", Object: v}
}

// NewLineBroadcast returns broadcast line notification task
func NewLineBroadcast(v model.LineNotifyMessage) event.Task {
	return NewLinesBroadcast([]model.LineNotifyMessage{v})
}

// NewLinesBroadcast returns broadcast line notification task
func NewLinesBroadcast(v []model.LineNotifyMessage) event.Task {
	return event.Task{QueueName: "queue-line", Path: "/line/notify/broadcast", Object: v}
}

// NewLine returns line task
func NewLine(v model.LineNotifyRequest) event.Task {
	return event.Task{QueueName: "queue-line", Path: "/line/notify", Object: v}
}
