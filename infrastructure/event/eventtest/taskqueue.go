package eventtest

import (
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/event"
)

// TaskQueue presents mock
type TaskQueue struct {
	log   core.Logger
	Tasks []event.Task
}

// NewTaskQueue returns no ops taskQueue instance
func NewTaskQueue(log core.Logger) *TaskQueue {
	return &TaskQueue{log, nil}
}

// Push add task
func (t *TaskQueue) Push(task event.Task) error {
	t.Tasks = append(t.Tasks, task)
	return nil
}

// PushMulti add tasks
func (t *TaskQueue) PushMulti(tasks []event.Task) error {
	t.Tasks = append(t.Tasks, tasks...)
	return nil
}
