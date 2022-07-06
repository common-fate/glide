package local

import (
	"context"

	"github.com/hashicorp/go-memdb"
)

// Runtime is a local runtime for testing use only which executes
// grants using goroutines. It stores grants in memory.
type Runtime struct {
	db *memdb.MemDB
}

// Init initialises the runtime and sets up the in-memory storage.
func (r *Runtime) Init(ctx context.Context) error {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"grants": {
				Name: "grants",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return err
	}

	r.db = db
	return nil
}
