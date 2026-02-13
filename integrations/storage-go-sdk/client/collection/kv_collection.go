package collection

import (
	"time"

	"github.com/fancyinnovations/fancyspaces/storage-sdk/client"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
)

// KeyValueCollection represents a key-value collection in the storage system.
type KeyValueCollection struct {
	collection
}

// NewKeyValueCollection creates a new KeyValueCollection instance for the specified database and collection name.
// It uses the provided client to interact with the storage system.
func NewKeyValueCollection(c *client.Client, database, name string) (*KeyValueCollection, error) {
	coll, err := c.DBCollectionGet(database, name)
	if err != nil {
		return nil, err
	}

	if coll.Engine != "kv" {
		return nil, client.ErrInvalidEngine
	}

	return &KeyValueCollection{
		collection: collection{
			database: database,
			name:     name,
			client:   c,
		},
	}, nil
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

func (coll *KeyValueCollection) GetTTL(key string) (time.Duration, error) {
	expiresAt, err := coll.client.KVGetTTL(coll.database, coll.name, key)
	if err != nil {
		return time.Duration(0), err
	}

	if expiresAt == 0 {
		return time.Duration(0), nil // Key does not exist or has no TTL
	}

	now := time.Now().UnixNano()
	if expiresAt <= now {
		return time.Duration(0), nil // Key has already expired
	}

	return time.Duration(expiresAt - now), nil
}

func (coll *KeyValueCollection) GetMultipleTTL(keys []string) (map[string]time.Duration, error) {
	expiresAtMap, err := coll.client.KVGetMultipleTTL(coll.database, coll.name, keys)
	if err != nil {
		return nil, err
	}

	result := make(map[string]time.Duration, len(expiresAtMap))
	now := time.Now().UnixNano()

	for key, expiresAt := range expiresAtMap {

		if expiresAt == 0 {
			result[key] = time.Duration(0) // Key does not exist or has no TTL
			continue
		}

		if expiresAt <= now {
			result[key] = time.Duration(0) // Key has already expired
			continue
		}

		result[key] = time.Duration(expiresAt - now)
	}

	return result, nil
}

func (coll *KeyValueCollection) GetAllTTL() (map[string]time.Duration, error) {
	expiresAtMap, err := coll.client.KVGetAllTTL(coll.database, coll.name)
	if err != nil {
		return nil, err
	}

	result := make(map[string]time.Duration, len(expiresAtMap))
	now := time.Now().UnixNano()

	for key, expiresAt := range expiresAtMap {
		if expiresAt == 0 {
			result[key] = time.Duration(0) // Key does not exist or has no TTL
			continue
		}

		if expiresAt <= now {
			result[key] = time.Duration(0) // Key has already expired
			continue
		}

		result[key] = time.Duration(expiresAt - now)
	}

	return result, nil
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

// SetWithTTL sets a key-value pair in the collection with a specified time-to-live (TTL) duration.
// After the TTL expires, the key will be automatically deleted from the collection.
// The value can be of any type supported by codex.ValueType.
func (coll *KeyValueCollection) SetWithTTL(key string, value any, ttl time.Duration) error {
	cval, err := codex.NewValue(value)
	if err != nil {
		return err
	}

	expiresAt := uint64(time.Now().Add(ttl).UnixNano())

	return coll.client.KVSetTTL(coll.database, coll.name, key, cval, expiresAt)
}

// SetStructWithTTL marshals the provided struct value and sets it in the collection under the specified key with a specified time-to-live (TTL) duration.
// After the TTL expires, the key will be automatically deleted from the collection.
func (coll *KeyValueCollection) SetStructWithTTL(key string, value any, ttl time.Duration) error {
	data, err := codex.Marshal(value)
	if err != nil {
		return err
	}

	return coll.SetWithTTL(key, data, ttl)
}

// Delete removes the key-value pair associated with the specified key from the collection.
func (coll *KeyValueCollection) Delete(key string) error {
	return coll.client.KVDelete(coll.database, coll.name, key)
}

// DeleteMultiple removes the key-value pairs associated with the specified keys from the collection.
func (coll *KeyValueCollection) DeleteMultiple(keys []string) error {
	return coll.client.KVDeleteMultiple(coll.database, coll.name, keys)
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

// Count returns the total number of key-value pairs present in the collection.
func (coll *KeyValueCollection) Count() (uint32, error) {
	return coll.client.KVCount(coll.database, coll.name)
}

// Size returns the total size in bytes of all key-value pairs present in the collection.
func (coll *KeyValueCollection) Size() (uint64, error) {
	return coll.client.KVSize(coll.database, coll.name)
}
