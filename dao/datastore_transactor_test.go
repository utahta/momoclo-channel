package dao

import (
	"testing"

	"github.com/utahta/momoclo-channel/testutil"
	"google.golang.org/appengine/aetest"
)

func TestDatastoreTransactor_RunInTransaction(t *testing.T) {
	ctx, done, err := testutil.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
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
	err = tran.RunInTransaction(func(p PersistenceHandler) error {
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
