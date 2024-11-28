package taskvault

import (
	"encoding/json"
	"errors"
	"io"
	"sync"
)

type Pair struct {
	Key string
	Value string
}

type KV struct {
	mu   sync.RWMutex
	vals map[string]string
}

func New() *KV {
	return &KV{
		vals: make(map[string]string),
	}
}

func (k *KV) Restore(r io.ReadCloser) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	defer r.Close()
	jsonData, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, &k.vals); err != nil {
		return err
	}

	return nil
}

func (k *KV) Get(key string) (string, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	return k.vals[key], nil
}

func (k *KV) Set(key, val string) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.vals[key] = val

	return nil
}

func (k *KV) SetNx(key, val string) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if _, ok := k.vals[key]; ok {
		return errors.New("key already exists")
	}

	k.vals[key] = val

	return nil
}
