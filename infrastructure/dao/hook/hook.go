package hook

import (
	"github.com/utahta/momoclo-channel/domain/core"
	"github.com/utahta/momoclo-channel/domain/model"
)

// BeforeSave hook
func BeforeSave(src interface{}) {
	if p, ok := src.(model.PersistenceBeforeSaver); ok {
		p.BeforeSave()
	}
}

// Validate hook
func Validate(src interface{}) error {
	return core.Validate(src)
}
