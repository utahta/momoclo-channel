package dao

import (
	"testing"

	"context"

	"github.com/utahta/momoclo-channel/testutil"
	"google.golang.org/appengine/aetest"
)

func TestDatastoreTransactor_RunInTransaction(t *testing.T) {
	ctx, done, err := testutil.NewContext(&aetest.Options{StronglyConsistentDatastore: true})
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	h := NewDatastoreHandler()
	e := &TestEntity{ID: "plan-2", Name: "taroimo"}
	if err := h.Put(ctx, e); err != nil {
		t.Fatal(err)
	}

	tran := NewDatastoreTransactor()
	err = tran.RunInTransaction(ctx, func(ctx context.Context) error {
		e.Name = "taroimo_z"
		if err := h.Put(ctx, e); err != nil {
			return err
		}

		if err := h.Get(ctx, e); err != nil {
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

	h.FlushLocalCache(ctx)
	if err := h.Get(ctx, e); err != nil {
		t.Fatal(err)
	}

	if e.Name != "taroimo_z" {
		t.Errorf("Expected taroimo_z, got %v", e.Name)
	}
}
