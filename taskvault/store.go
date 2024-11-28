package taskvault

import (
	"io"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
)

// Store is the local implementation of the Storage interface.
// It gives dkron the ability to manipulate its embedded storage
// BuntDB.
type Store struct {
	db   *buntdb.DB
	lock *sync.Mutex // for

	logger *logrus.Entry
}

// DeleteValue implements Storage.
func (s *Store) DeleteValue(key string) error {
	err := s.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})

	return err
}

// GetAllValues implements Storage.
func (s *Store) GetAllValues() ([]Pair, error) {
	var pairs []Pair

	err := s.db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(k, v string) bool {
			pairs = append(pairs, Pair{
				Key:   k,
				Value: v,
			})
			return true
		})

		return err
	})

	return pairs, err
}

// GetValue implements Storage.
func (s *Store) GetValue(key string) (string, error) {
	var value string

	err := s.db.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(key)
		if err != nil {
			return err
		}

		value = v

		s.logger.WithFields(logrus.Fields{
			"value": v,
			"key":   key,
		}).Debug("store: Retrieved value from database")

		return nil
	})

	return value, err
}

// Restore implements Storage.
func (s *Store) Restore(r io.ReadCloser) error {
	return s.db.Load(r)
}

// SetValue implements Storage.
func (s *Store) SetValue(key string, value string) error {
	err := s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})

	return err
}

// Shutdown implements Storage.
func (s *Store) Shutdown() error {
	return s.db.Close()
}

// Snapshot implements Storage.
func (s *Store) Snapshot(w io.WriteCloser) error {
	return s.db.Save(w)
}

// UpdateValue implements Storage.
func (s *Store) UpdateValue(key string, value string) error {
	return s.SetValue(key, value)
}

var _ Storage = (*Store)(nil)

type kv struct {
	Key   string
	Value string
}

func NewStore(logger *logrus.Entry) (*Store, error) {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}

	store := &Store{
		db:     db,
		lock:   &sync.Mutex{},
		logger: logger,
	}

	return store, nil
}
