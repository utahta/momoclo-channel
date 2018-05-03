package eventtask

import (
	"github.com/utahta/momoclo-channel/crawler"
	"github.com/utahta/momoclo-channel/event"
	"github.com/utahta/momoclo-channel/twitter"
	"github.com/utahta/momoclo-channel/types"
)

// NewEnqueueTweets returns enqueue tweets task
func NewEnqueueTweets(v crawler.FeedItem) event.Task {
	return event.Task{QueueName: "enqueue", Path: "/enqueue/tweets", Object: v}
}

// NewEnqueueLines returns enqueue lines task
func NewEnqueueLines(v crawler.FeedItem) event.Task {
	return event.Task{QueueName: "enqueue", Path: "/enqueue/lines", Object: v}
}

// NewTweet returns tweet task
func NewTweet(v twitter.TweetRequest) event.Task {
	return NewTweets([]twitter.TweetRequest{v})
}

// NewTweets returns tweet task
func NewTweets(v []twitter.TweetRequest) event.Task {
	return event.Task{QueueName: "queue-tweet", Path: "/tweet", Object: v}
}

// NewLineBroadcast returns broadcast line notification task
func NewLineBroadcast(v types.LineNotifyMessage) event.Task {
	return NewLinesBroadcast([]types.LineNotifyMessage{v})
}

// NewLinesBroadcast returns broadcast line notification task
func NewLinesBroadcast(v []types.LineNotifyMessage) event.Task {
	return event.Task{QueueName: "queue-line", Path: "/line/notify/broadcast", Object: v}
}

// NewLine returns line task
func NewLine(v types.LineNotifyRequest) event.Task {
	return event.Task{QueueName: "queue-line", Path: "/line/notify", Object: v}
}
