package api

import (
	"context"
	"net/http"

	"github.com/utahta/momoclo-channel/log"
)

// failResponse responses error
func failResponse(ctx context.Context, w http.ResponseWriter, err error, code int) {
	var message string
	if err != nil {
		message = err.Error()
	}

	log.NewAELogger().Errorf(ctx, "An error has occurred! code:%v err:%+v", code, err)
	http.Error(w, message, code)
}
