package hook

import (
	"github.com/utahta/momoclo-channel/validator"
)

type (
	// PersistenceBeforeSaver hook
	PersistenceBeforeSaver interface {
		BeforeSave()
	}
)

// BeforeSave hook
func BeforeSave(src interface{}) {
	if p, ok := src.(PersistenceBeforeSaver); ok {
		p.BeforeSave()
	}
}

// Validate hook
func Validate(src interface{}) error {
	return validator.Validate(src)
}
