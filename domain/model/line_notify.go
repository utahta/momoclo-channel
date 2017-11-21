package model

type (
	// LineNotifyToken interface
	LineNotifyToken interface {
		GetAccessToken(string) (string, error)
	}
)
