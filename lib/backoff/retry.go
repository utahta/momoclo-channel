package backoff

import "github.com/cenkalti/backoff"

// Retry wraps cenkalti/backoff.Retry
func Retry(max uint64, o backoff.Operation) error {
	return backoff.Retry(o, backoff.WithMaxTries(backoff.NewExponentialBackOff(), max))
}
