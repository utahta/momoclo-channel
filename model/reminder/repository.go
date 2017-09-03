package reminder

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/model"
	"google.golang.org/appengine/datastore"
)

type repository struct{}

var Repository *repository = &repository{}

func (repo *repository) GetAll(ctx context.Context) ([]*model.Reminder, error) {
	q := datastore.NewQuery("Reminder").Filter("Enabled =", true)

	var dst []*model.Reminder
	for t := q.Run(ctx); ; {
		var rd model.Reminder
		k, err := t.Next(&rd)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get all Reminders.")
		}
		rd.Id = k.IntID()
		dst = append(dst, &rd)
	}
	return dst, nil
}
