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

	req, err := t.newPOSTTask(task)
	if err != nil {
		return errors.Wrap(err, errTag)
	}

	if _, err := taskqueue.Add(t.ctx, req, task.QueueName); err != nil {
		return errors.Wrap(err, errTag)
	}
	return nil
}

func (t *taskQueue) PushMulti(tasks []event.Task) error {
	const errTag = "taskQueue.PushMulti failed"

	if len(tasks) == 0 {
		return nil
	}

	reqsMap := map[string][]*taskqueue.Task{}
	for _, task := range tasks {
		req, err := t.newPOSTTask(task)
		if err != nil {
			return errors.Wrap(err, errTag)
		}
		reqsMap[task.QueueName] = append(reqsMap[task.QueueName], req)
	}

	// see: https://cloud.google.com/appengine/quotas?hl=en#Task_Queue
	const maxTaskQueueNum = 100
	for queueName, reqs := range reqsMap {
		for i := 0; i < len(reqs); i += maxTaskQueueNum {
			last := i + maxTaskQueueNum
			if last > len(reqs) {
				last = len(reqs)
			}

			_, err := taskqueue.AddMulti(t.ctx, reqs[i:last], queueName)
			if err != nil {
				return errors.Wrap(err, errTag)
			}
		}
	}
	return nil
}

func (t *taskQueue) newPOSTTask(task event.Task) (*taskqueue.Task, error) {
	v, err := task.Params()
	if err != nil {
		return nil, errors.Wrapf(err, "taskQueue.newPOSTTask failed: task:%v", task)
	}

	req := taskqueue.NewPOSTTask(task.Path, v)
	if task.Delay > 0 {
		req.Delay = task.Delay
	}

	opts := &taskqueue.RetryOptions{}
	if task.RetryLimit > 0 {
		opts.RetryLimit = int32(task.RetryLimit)
	}
	req.RetryOptions = opts

	return req, nil
}
