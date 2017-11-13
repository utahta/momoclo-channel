package event

import (
	"encoding/json"
	"net/url"
	"time"
)

// TaskQueue interface
type TaskQueue interface {
	Push(Task) error
}

// Task event
type Task struct {
	QueueName string
	Path      string
	Object    interface{}
	Payload   []byte
	Delay     time.Duration
}

// Params sets payload to url.Values
func (t *Task) Params() (url.Values, error) {
	v := url.Values{}

	if t.Object != nil {
		b, err := json.Marshal(t.Object)
		if err != nil {
			return v, err
		}
		t.Payload = b
	}

	v.Set("payload", string(t.Payload))
	return v, nil
}

// Parse parses url.Values
func (t *Task) Parse(v url.Values, o interface{}) error {
	t.Payload = []byte(v.Get("payload"))
	return json.Unmarshal(t.Payload, o)
}
