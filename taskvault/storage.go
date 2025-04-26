package taskvault

import (
	"io"

	"github.com/hashicorp/raft"
)

type SyncraStorage interface {
	GetValue(key string) (string, error)
	UpdateValue(key string, value string) error
	SetValue(key string, value string) error
	DeleteValue(key string) error
	GetAllValues() ([]Pair, error)
	Shutdown() error
	Snapshot(w io.WriteCloser) error
	Restore(r io.ReadCloser) error
}

type RaftStore interface {
	raft.StableStore
	raft.LogStore
	Close() error
}
