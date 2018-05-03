package usecase

import (
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/config"
	"github.com/utahta/momoclo-channel/linenotify"
	"github.com/utahta/momoclo-channel/log"
	"github.com/utahta/momoclo-channel/types"
)

type (
	// AddLineNotification use case
	AddLineNotification struct {
		log   log.Logger
		token linenotify.Token
		repo  types.LineNotificationRepository
	}

	// AddLineNotificationParams use case params
	AddLineNotificationParams struct {
		Code string
	}
)

// NewAddLineNotification returns AddLineNotification use case
func NewAddLineNotification(
	logger log.Logger,
	token linenotify.Token,
	repo types.LineNotificationRepository) *AddLineNotification {
	return &AddLineNotification{
		log:   logger,
		token: token,
		repo:  repo,
	}
}

// Do add line notification entity
func (use *AddLineNotification) Do(params AddLineNotificationParams) error {
	const errTag = "AddLineNotification.Do failed"

	token, err := use.token.GetAccessToken(params.Code)
	if err != nil {
		return errors.Wrap(err, errTag)
	}

	ln, err := types.NewLineNotification(config.C.LineNotify.TokenKey, token)
	if err != nil {
		return errors.Wrap(err, errTag)
	}

	if err := use.repo.Save(ln); err != nil {
		return errors.Wrap(err, errTag)
	}

	use.log.Infof("added LineNotification. id:%v", ln.ID)
	return nil
}
