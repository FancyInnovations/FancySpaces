package collection

import (
	"time"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/client"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
)

// KeyValueCollection represents a key-value collection in the storage system.
type KeyValueCollection struct {
	collection
}

// NewKeyValueCollection creates a new KeyValueCollection instance for the specified database and collection name.
// It uses the provided client to interact with the storage system.
func NewKeyValueCollection(client *client.Client, database, name string) *KeyValueCollection {
	// TODO: check if database and collection exist, return error if not

	// TODO check if collection is of type KV, return error if not

	return &KeyValueCollection{
		collection: collection{
			database: database,
			name:     name,
			client:   client,
		},
	}
}

// Get retrieves the value associated with the specified key from the collection.
func (coll *KeyValueCollection) Get(key string) (*codex.Value, error) {
	return coll.client.KVGet(coll.database, coll.name, key)
}

// GetStruct retrieves the value associated with the specified key and tries to unmarshal it into the provided destination struct.
func (coll *KeyValueCollection) GetStruct(key string, dest any) error {
	val, err := coll.Get(key)
	if err != nil {
		return err
	}

	if val.Type != codex.TypeBinary {
		return codex.ErrInvalidType
	}

	return codex.Unmarshal(val.AsBinary(), dest)
}

// GetMultiple retrieves the values associated with the specified keys from the collection and returns them as a map.
// The keys in the returned map are the keys from the collection, and the values are the corresponding codex.Value.
func (coll *KeyValueCollection) GetMultiple(keys []string) (map[string]*codex.Value, error) {
	return coll.client.KVGetMultiple(coll.database, coll.name, keys)
}

// GetMultipleStruct retrieves the values associated with the specified keys and tries to unmarshal each value into the provided destination map.
// The keys in the destination map will be the keys from the collection, and the values will be the unmarshaled struct values.
func (coll *KeyValueCollection) GetMultipleStruct(keys []string, dest map[string]any) error {
	vals, err := coll.GetMultiple(keys)
	if err != nil {
		return err
	}

	for k, v := range vals {
		if v.Type != codex.TypeBinary {
			return codex.ErrInvalidType
		}

		var out any
		if err := codex.Unmarshal(v.AsBinary(), &out); err != nil {
			return err
		}
		dest[k] = out
	}

	return nil
}

// GetAll retrieves all key-value pairs from the collection and returns them as a map.
// The keys in the returned map are the keys from the collection, and the values are the corresponding codex.Value.
func (coll *KeyValueCollection) GetAll() (map[string]*codex.Value, error) {
	return coll.client.KVGetAll(coll.database, coll.name)
}

// GetAllStruct retrieves all key-value pairs from the collection and tries to unmarshal each value into the provided destination map.
// The keys in the destination map will be the keys from the collection, and the values will be the unmarshaled struct values.
func (coll *KeyValueCollection) GetAllStruct(dest map[string]any) error {
	vals, err := coll.GetAll()
	if err != nil {
		return err
	}

	for k, v := range vals {
		if v.Type != codex.TypeBinary {
			return codex.ErrInvalidType
		}

		var out any
		if err := codex.Unmarshal(v.AsBinary(), &out); err != nil {
			return err
		}
		dest[k] = out
	}

	return nil
}

// Set sets a key-value pair in the collection.
// The value can be of any type supported by codex.ValueType.
func (coll *KeyValueCollection) Set(key string, value any) error {
	cval, err := codex.NewValue(value)
	if err != nil {
		return err
	}

	return coll.client.KVSet(coll.database, coll.name, key, cval)
}

// SetStruct marshals the provided struct value and sets it in the collection under the specified key.
func (coll *KeyValueCollection) SetStruct(key string, value any) error {
	data, err := codex.Marshal(value)
	if err != nil {
		return err
	}

	return coll.Set(key, data)
}

// SetWithTTL sets a key-value pair in the collection with a specified time-to-live (TTL) in milliseconds.
// After the TTL expires, the key will be automatically deleted from the collection.
// The value can be of any type supported by codex.ValueType.
func (coll *KeyValueCollection) SetWithTTL(key string, value any, ttlMillis uint64) error {
	cval, err := codex.NewValue(value)
	if err != nil {
		return err
	}

	ttlNanos := ttlMillis * 1_000_000 // Convert milliseconds to nanoseconds
	expiresAt := uint64(time.Now().UnixNano()) + ttlNanos

	return coll.client.KVSetTTL(coll.database, coll.name, key, cval, expiresAt)
}

// SetStructWithTTL marshals the provided struct value and sets it in the collection under the specified key with a TTL.
// After the TTL expires, the key will be automatically deleted from the collection.
func (coll *KeyValueCollection) SetStructWithTTL(key string, value any, ttlMillis uint64) error {
	data, err := codex.Marshal(value)
	if err != nil {
		return err
	}

	ttlNanos := ttlMillis * 1_000_000 // Convert milliseconds to nanoseconds
	expiresAt := uint64(time.Now().UnixNano()) + ttlNanos

	return coll.SetWithTTL(key, data, expiresAt)
}

// Delete removes the key-value pair associated with the specified key from the collection.
func (coll *KeyValueCollection) Delete(key string) error {
	return coll.client.KVDelete(coll.database, coll.name, key)
}

// DeleteAll removes all key-value pairs from the collection.
func (coll *KeyValueCollection) DeleteAll() error {
	return coll.client.KVDeleteAll(coll.database, coll.name)
}

// Exists checks if a key exists in the collection and returns true if it does, false otherwise.
func (coll *KeyValueCollection) Exists(key string) (bool, error) {
	return coll.client.KVExists(coll.database, coll.name, key)
}

// Keys retrieves a list of all keys present in the collection.
func (coll *KeyValueCollection) Keys() ([]string, error) {
	return coll.client.KVKeys(coll.database, coll.name)
}

func (coll *KeyValueCollection) Count() (uint32, error) {
	return coll.client.KVCount(coll.database, coll.name)
}
