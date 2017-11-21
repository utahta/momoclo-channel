package model

type (
	// LineNotifyTokenGetter interface
	LineNotifyTokenGetter interface {
		Get(string) (string, error)
	}
)
