package hook

import (
	"reflect"

	"github.com/utahta/momoclo-channel/domain/model"
)

// BeforeSave hook
func BeforeSave(src interface{}) {
	if p, ok := src.(model.PersistenceBeforeSaver); ok {
		p.BeforeSave()
	}
}

// BeforeSaveMulti hook
func BeforeSaveMulti(src interface{}) {
	//TODO need to prevent panic given invalid args. (e.g. validation)
	v := reflect.ValueOf(src)
	for i := 0; i < v.Len(); i++ {
		BeforeSave(v.Index(i).Interface())
	}
}
