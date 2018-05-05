package eventtest

import (
	"context"

	"github.com/utahta/momoclo-channel/event"
)

// TaskQueue presents mock
type TaskQueue struct {
	Tasks []event.Task
}

// NewTaskQueue returns no ops taskQueue instance
func NewTaskQueue() *TaskQueue {
	return &TaskQueue{}
}

// Push add task
func (t *TaskQueue) Push(_ context.Context, task event.Task) error {
	t.Tasks = append(t.Tasks, task)
	return nil
}

// PushMulti add tasks
func (t *TaskQueue) PushMulti(_ context.Context, tasks []event.Task) error {
	t.Tasks = append(t.Tasks, tasks...)
	return nil
}
