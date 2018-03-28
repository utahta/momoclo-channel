package hook

import (
	"github.com/utahta/momoclo-channel/types"
	"github.com/utahta/momoclo-channel/validator"
)

// BeforeSave hook
func BeforeSave(src interface{}) {
	if p, ok := src.(types.PersistenceBeforeSaver); ok {
		p.BeforeSave()
	}
}

// Validate hook
func Validate(src interface{}) error {
	return validator.Validate(src)
}
