package dao

import (
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/utahta/momoclo-channel/testutil"
	"google.golang.org/appengine/aetest"
)

type TestEntity struct {
	ID        string `datastore:"-" goon:"id" validate:"required"`
	Name      string `validate:"required"`
	CreatedAt time.Time
}

func (e *TestEntity) BeforeSave() {
	e.CreatedAt = time.Now()
}

func TestDatastoreHandler_Put(t *testing.T) {
	ctx, done, err := testutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	h := NewDatastoreHandler(ctx)
	if err := h.Put(&TestEntity{ID: "plan-1", Name: "ageshio-kyukou"}); err != nil {
		t.Fatal(err)
	}

	e := TestEntity{ID: "plan-1"}
	if err := h.Get(&e); err != nil {
		t.Fatal(err)
	}

	if e.Name != "ageshio-kyukou" {
		t.Errorf("Expected ageshio-kyukou, got %v", e.Name)
	}

	if e.CreatedAt.IsZero() {
		t.Errorf("Expected set createdAt, got %v", e.CreatedAt)
	}

	err = h.Put(&TestEntity{ID: "plan-2"})
	if errs, ok := err.(validator.ValidationErrors); !ok {
		t.Errorf("Expected validation errors, got %v", errs)
	}
}

func TestDatastoreHandler_PutMulti(t *testing.T) {
	ctx, done, err := testutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	h := NewDatastoreHandler(ctx)

	es := []*TestEntity{
		{ID: "plan-1", Name: "ageshio-kyukou"},
		{ID: "plan-2", Name: "taroimo"},
	}
	if err := h.PutMulti(es); err != nil {
		t.Fatal(err)
	}

	es = []*TestEntity{
		{ID: "plan-1"},
		{ID: "plan-2"},
	}
	if err := h.GetMulti(es); err != nil {
		t.Fatal(err)
	}

	if es[0].Name != "ageshio-kyukou" || es[1].Name != "taroimo" {
		t.Errorf("Expected ageshio-kyukou and taroimo, got 0:%v, 1:%v", es[0].Name, es[1].Name)
	}

	if es[0].CreatedAt.IsZero() || es[1].CreatedAt.IsZero() {
		t.Errorf("Expected set createdAt, got 0:%v, 1:%v", es[0].CreatedAt, es[1].CreatedAt)
	}
}
