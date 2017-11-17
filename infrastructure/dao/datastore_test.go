package dao

import (
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/domain/model"
	"github.com/utahta/momoclo-channel/lib/aetestutil"
	"google.golang.org/appengine/aetest"
)

type TestEntity struct {
	ID        string `datastore:"-" goon:"id"`
	Name      string
	CreatedAt time.Time
}

func (e *TestEntity) BeforeSave() {
	e.CreatedAt = time.Now()
}

func TestDatastoreHandler_Put(t *testing.T) {
	ctx, done, err := aetestutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
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
}

func TestDatastoreHandler_PutMulti(t *testing.T) {
	ctx, done, err := aetestutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
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

func TestDatastoreTransactor_RunInTransaction(t *testing.T) {
	ctx, done, err := aetestutil.NewContex(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	h := NewDatastoreHandler(ctx)
	e := &TestEntity{ID: "plan-2", Name: "taroimo"}
	if err := h.Put(e); err != nil {
		t.Fatal(err)
	}

	tran := NewDatastoreTransactor(ctx)
	err = tran.RunInTransaction(func(p model.PersistenceHandler) error {
		e.Name = "taroimo_z"
		if err := p.Put(e); err != nil {
			return err
		}

		if err := p.Get(e); err != nil {
			return err
		}

		// expected to not commit yet
		if e.Name != "taroimo" {
			t.Errorf("Expected taroimo, got %v", e.Name)
		}
		return nil
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	h.FlushLocalCache()
	if err := h.Get(e); err != nil {
		t.Fatal(err)
	}

	if e.Name != "taroimo_z" {
		t.Errorf("Expected taroimo_z, got %v", e.Name)
	}
}
