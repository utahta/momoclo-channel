package event

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/domain/event"
	"google.golang.org/appengine/taskqueue"
)

type taskQueue struct {
	ctx context.Context
}

// NewTaskQueue returns event.TaskQueue that wraps appengine/taskqueue
func NewTaskQueue(ctx context.Context) event.TaskQueue {
	return &taskQueue{ctx}
}

func (t *taskQueue) Push(task event.Task) error {
	const errTag = "taskQueue.Push failed"

	v, err := task.Params()
	if err != nil {
		return errors.Wrap(err, errTag)
	}

	postTask := taskqueue.NewPOSTTask(task.Path, v)
	if task.Delay > 0 {
		postTask.Delay = task.Delay
	}

	if _, err := taskqueue.Add(t.ctx, postTask, task.QueueName); err != nil {
		return errors.Wrap(err, errTag)
	}

	return nil
}
