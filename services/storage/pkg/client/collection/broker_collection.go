package collection

import "github.com/fancyinnovations/fancyspaces/storage/pkg/client"

// MessageBrokerCollection represents a collection that provides message broker functionality in the storage system.
type MessageBrokerCollection struct {
	collection
}

// NewMessageBrokerCollection creates a new MessageBrokerCollection instance for the specified database and collection name.
// It uses the provided client to interact with the storage system.
func NewMessageBrokerCollection(client *client.Client, database, name string) *MessageBrokerCollection {
	// TODO: check if database and collection exist, return error if not

	// TODO check if collection is of type MessageBroker, return error if not

	return &MessageBrokerCollection{
		collection: collection{
			database: database,
			name:     name,
			client:   client,
		},
	}
}

// Subscribe subscribes the client to a subject in the collection.
//
// The wildcards "*" and ">" can be used in the subject to match multiple subjects.
// For example, "foo.*" will match "foo.bar" and "foo.baz", while "foo.>" will match "foo.bar", "foo.baz", and any other subject that starts with "foo.".
//
// The fn parameter is a callback function that will be called whenever a message is published to the subscribed subject.
// The message will be passed as a byte slice to the callback function.
func (coll *MessageBrokerCollection) Subscribe(subject string, fn func(msg []byte)) error {
	return coll.client.BrokerSubscribe(coll.database, coll.name, subject, fn)
}

// SubscribeQueue subscribes the client to a subject in the collection with a queue group.
//
// The wildcards "*" and ">" can be used in the subject to match multiple subjects.
// For example, "foo.*" will match "foo.bar" and "foo.baz", while "foo.>" will match "foo.bar", "foo.baz", and any other subject that starts with "foo.".
//
// The queue parameter is a string that identifies the queue group for load balancing messages among multiple subscribers.
// When multiple subscribers are subscribed to the same subject and queue group, messages published to that subject will be distributed among the subscribers in a round-robin fashion.
//
// The fn parameter is a callback function that will be called whenever a message is published to the subscribed subject and queue group.
// The message will be passed as a byte slice to the callback function.
func (coll *MessageBrokerCollection) SubscribeQueue(subject, queue string, fn func(msg []byte)) error {
	return coll.client.BrokerSubscribeQueue(coll.database, coll.name, subject, queue, fn)
}

// Unsubscribe unsubscribes the client from a subject in the collection.
func (coll *MessageBrokerCollection) Unsubscribe(subject string) error {
	return coll.client.BrokerUnsubscribe(coll.database, coll.name, subject)
}

// Publish publishes a message to a subject in the collection.
func (coll *MessageBrokerCollection) Publish(subject string, msg []byte) error {
	return coll.client.BrokerPublish(coll.database, coll.name, subject, msg)
}
