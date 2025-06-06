package taskvault

import (
	"io"

	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

type Store struct {
	db *buntdb.DB

	logger *zap.SugaredLogger
}

var _ SyncraStorage = (*Store)(nil)

func (s *Store) DeleteValue(key string) error {
	err := s.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})

	return err
}

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

func (s *Store) GetValue(key string) (string, error) {
	var value string

	err := s.db.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(key)
		if err != nil {
			return err
		}

		value = v

		return nil
	})

	return value, err
}

func (s *Store) Restore(r io.ReadCloser) error {
	return s.db.Load(r)
}

func (s *Store) SetValue(key string, value string) error {
	err := s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})

	return err
}

func (s *Store) Shutdown() error {
	return s.db.Close()
}

func (s *Store) Snapshot(w io.WriteCloser) error {
	return s.db.Save(w)
}

func (s *Store) UpdateValue(key string, value string) error {
	return s.SetValue(key, value)
}

var _ SyncraStorage = (*Store)(nil)

func NewStore(logger *zap.SugaredLogger) (*Store, error) {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}

	store := &Store{
		db:     db,
		logger: logger,
	}

	return store, nil
}
