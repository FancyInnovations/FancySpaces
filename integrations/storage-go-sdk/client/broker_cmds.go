package client

import (
	"encoding/binary"

	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/codex"
	"github.com/fancyinnovations/fancyspaces/integrations/storage-go-sdk/protocol"
)

// BrokerSubscribe implements the client side of the protocol.ServerCommandBrokerSubscribe command.
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

// BrokerSubscribeQueue implements the client side of the protocol.ServerCommandBrokerSubscribeQueue command.
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

// BrokerUnsubscribe implements the client side of the protocol.ServerCommandBrokerUnsubscribe command.
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

// BrokerPublish implements the client side of the protocol.ServerCommandBrokerPublish command.
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
