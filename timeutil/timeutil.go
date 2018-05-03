package timeutil

import (
	"time"
)

var (
	// Now wraps time.Now (eventually for test)
	Now = time.Now

	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
)

// JST returns a Location that uses Asia/Tokyo.
func JST() *time.Location {
	return jst
}
