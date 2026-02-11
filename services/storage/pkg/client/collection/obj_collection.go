package collection

import "github.com/fancyinnovations/fancyspaces/storage/pkg/client"

// ObjectCollection represents a collection that stores objects in the storage system.
type ObjectCollection struct {
	collection
}

// NewObjectCollection creates a new ObjectCollection instance for the specified database and collection name.
// It uses the provided client to interact with the storage system.
func NewObjectCollection(client *client.Client, database, name string) *ObjectCollection {
	// TODO: check if database and collection exist, return error if not

	// TODO check if collection is of type Object, return error if not

	return &ObjectCollection{
		collection: collection{
			database: database,
			name:     name,
			client:   client,
		},
	}
}

// Put stores the given binary data in the collection under the specified key.
func (coll *ObjectCollection) Put(key string, data []byte) error {
	return coll.client.ObjPut(coll.database, coll.name, key, data)
}

// Get retrieves the binary data associated with the specified key from the collection.
func (coll *ObjectCollection) Get(key string) ([]byte, error) {
	return coll.client.ObjGet(coll.database, coll.name, key)
}

// GetMetadata retrieves the metadata associated with the specified key from the collection.
func (coll *ObjectCollection) GetMetadata(key string) (*client.ObjectMetadata, error) {
	return coll.client.ObjGetMetadata(coll.database, coll.name, key)
}

// Delete removes the object associated with the specified key from the collection.
func (coll *ObjectCollection) Delete(key string) error {
	return coll.client.ObjDelete(coll.database, coll.name, key)
}

// Count returns the total number of objects stored in the collection.
func (coll *ObjectCollection) Count() (uint32, error) {
	return coll.client.ObjCount(coll.database, coll.name)
}

// Size returns the total size of all objects stored in the collection in bytes.
func (coll *ObjectCollection) Size() (uint64, error) {
	return coll.client.ObjSize(coll.database, coll.name)
}
