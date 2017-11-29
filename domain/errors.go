package domain

import "github.com/pkg/errors"

var (
	ErrNoSuchEntity       = errors.New("mcz: no such entity")
	ErrInvalidSignature   = errors.New("mcz: invalid signature")
	ErrInvalidAccessToken = errors.New("mcz: invalid access token")
)
