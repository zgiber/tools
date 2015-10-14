package memstore

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"time"
)

type Snapshot struct {
	Version uint64                 `json:"ver"`
	Values  map[string]interface{} `json:"val"`
}

type KeyValueStore interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

type PersistentKeyValueStore interface {
	KeyValueStore
	SaveSnapshot() error
	LoadSnapshot() error
}

type MemoryStore struct {
	l         sync.RWMutex
	values    map[string]interface{}
	version   uint64
	lastSaved time.Time
	fileName  string
	file      *os.File
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		values: map[string]interface{}{},
	}
}

func (ms *MemoryStore) Get(key string) (interface{}, bool) {
	ms.l.RLock()
	value, exists := ms.values[key]
	ms.l.RUnlock()
	return value, exists

}

func (ms *MemoryStore) Set(key string, value interface{}) {
	ms.l.Lock()
	if value == nil {
		delete(ms.values, key)
	} else {
		ms.values[key] = value
	}

	ms.version++
	ms.l.Unlock()
	return
}

func (ms *MemoryStore) Keys() []string {
	keys := []string{}
	ms.l.RLock()
	for key, _ := range ms.values {
		keys = append(keys, key)
	}
	ms.l.RUnlock()
	return keys
}

func (ms *MemoryStore) SaveSnapshot(w io.Writer) error {
	enc := json.NewEncoder(w)

	ms.l.RLock()
	defer ms.l.RUnlock()

	snapshot := &Snapshot{ms.version, ms.values}
	err := enc.Encode(snapshot)
	if err != nil {
		return err
	}

	ms.lastSaved = time.Now().UTC()
	return nil
}

func (ms *MemoryStore) LoadSnapshot(r io.Reader) error {

	snapshot := &Snapshot{}
	dec := json.NewDecoder(r)
	err := dec.Decode(snapshot)
	if err != nil {
		return err
	}

	ms.l.Lock()
	ms.values = snapshot.Values
	ms.version = snapshot.Version
	ms.l.Unlock()

	return nil
}

// GetValue takes dst which must pointer to a string, int, float64, bool, or struct
// and key as a string for parameters. If the MemoryStore has an entry for the
// key, and the dst type matches the stored type, then the dst's value will be
// given the value found in the MemoryStore
func (ms *MemoryStore) GetValue(dst interface{}, key string) (err error) {
	v, ok := ms.Get(key)
	if !ok {
		return fmt.Errorf("Value not found.")
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Invalid type: dst must be pointer type.")
		}
	}()

	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.Indirect(reflect.ValueOf(v))
	if srcValue.Type() == dstValue.Type() {
		dstValue.Set(srcValue)
	} else {
		return fmt.Errorf("Invalid type: expected %v, found %v.", dstValue.Type(), srcValue.Type())
	}
	return nil
}
