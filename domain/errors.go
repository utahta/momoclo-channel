package domain

import "github.com/pkg/errors"

var (
	ErrNoSuchEntity     = errors.New("error no such entity")
	ErrInvalidSignature = errors.New("error invalid signature")
)
