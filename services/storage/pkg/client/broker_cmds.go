package client

import (
	"encoding/binary"

	"github.com/fancyinnovations/fancyspaces/storage/pkg/codex"
	"github.com/fancyinnovations/fancyspaces/storage/pkg/protocol"
)

// BrokerSubscribe subscribes the client to a subject in the specified collection.
// The subject is a string that identifies the topic or channel to which the client wants to subscribe.
// The fn parameter is a callback function that will be called whenever a message is published to the subscribed subject.
// The message will be passed as a byte slice to the callback function.
func (c *Client) BrokerSubscribe(db, coll string, subject string, fn func(msg []byte)) error {
	totalLen := 2 + len(subject)
	payload := make([]byte, totalLen)

	// Subject
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(subject)))
	copy(payload[2:2+len(subject)], []byte(subject))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandBrokerSubscribe,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		return ErrUnexpectedStatusCode
	}

	c.brokerSubjectListenersMu.Lock()
	c.brokerSubjectListeners[db+"."+coll+"."+subject] = append(c.brokerSubjectListeners[db+"."+coll+"."+subject], fn)
	c.brokerSubjectListenersMu.Unlock()

	return nil
}

// BrokerSubscribeQueue subscribes the client to a subject in the specified collection with a queue group.
// The subject is a string that identifies the topic or channel to which the client wants to subscribe.
// The queue parameter is a string that identifies the queue group for load balancing messages among multiple subscribers.
// The fn parameter is a callback function that will be called whenever a message is published to the subscribed subject and queue group.
// The message will be passed as a byte slice to the callback function.
func (c *Client) BrokerSubscribeQueue(db, coll string, subject, queue string, fn func(msg []byte)) error {
	totalLen := 2 + len(subject) + 2 + len(queue)
	payload := make([]byte, totalLen)

	// Subject
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(subject)))
	copy(payload[2:2+len(subject)], []byte(subject))

	// Queue
	queueOffset := 2 + len(subject)
	binary.BigEndian.PutUint16(payload[queueOffset:queueOffset+2], uint16(len(queue)))
	copy(payload[queueOffset+2:queueOffset+2+len(queue)], []byte(queue))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandBrokerSubscribeQueue,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		return ErrUnexpectedStatusCode
	}
	return nil
}

// BrokerUnsubscribe unsubscribes the client from a subject in the specified collection.
// The subject is a string that identifies the topic or channel from which the client wants to unsubscribe.
// After calling this method, the client will no longer receive messages published to the specified subject.
func (c *Client) BrokerUnsubscribe(db, coll string, subject string) error {
	totalLen := 2 + len(subject)
	payload := make([]byte, totalLen)

	// Subject
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(subject)))
	copy(payload[2:2+len(subject)], []byte(subject))

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandBrokerUnsubscribe,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		return ErrUnexpectedStatusCode
	}

	c.brokerSubjectListenersMu.Lock()
	delete(c.brokerSubjectListeners, db+"."+coll+"."+subject)
	c.brokerSubjectListenersMu.Unlock()

	return nil
}

// BrokerPublish publishes a message to a subject in the specified collection.
// The subject is a string that identifies the topic or channel to which the message will be published.
// The message is a byte slice that contains the data to be sent to subscribers of the specified subject.
func (c *Client) BrokerPublish(db, coll string, subject string, msg []byte) error {
	payload := make([]byte, 2+len(subject))

	// Subject
	binary.BigEndian.PutUint16(payload[0:2], uint16(len(subject)))
	copy(payload[2:2+len(subject)], []byte(subject))

	// Message
	payload = append(payload, codex.EncodeBinary(msg)...)

	resp, err := c.SendCmd(&protocol.Command{
		ID:             protocol.ServerCommandBrokerPublish,
		DatabaseName:   db,
		CollectionName: coll,
		Payload:        payload,
	})
	if err != nil {
		return err
	}

	if resp.Code != protocol.StatusOK {
		return ErrUnexpectedStatusCode
	}
	return nil
}
