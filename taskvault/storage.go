package taskvault

import "io"

type Storage interface {
	GetValue(key string) (string, error)
	UpdateValue(key string, value string) error
	SetValue(key string, value string) error
	DeleteValue(key string) error
	GetAllValues() ([]Pair, error)
	Shutdown() error
	Snapshot(w io.WriteCloser) error
	Restore(r io.ReadCloser) error
}
