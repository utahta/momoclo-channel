package model

type (
	// LineNotifyTokenClient interface
	LineNotifyTokenClient interface {
		GetToken(string) (string, error)
	}
)
