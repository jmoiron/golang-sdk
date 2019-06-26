// Copyright 2019, 1533 Systems, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package driveline

import (
	"context"
	"encoding/hex"
)

var (
	recordIDHead = RecordID([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	recordIDTail = RecordID([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
)

// RecordID is an opaque value that is used to identify a record in Driveline.
type RecordID []byte

// String implements Stringer interface. It's an hex-dump of the RecordID opaque data.
func (id RecordID) String() string {
	return hex.EncodeToString(id)
}

// Record is the core data exchange structure
type Record struct {
	RecordID RecordID // RecordID is the identifier of the Record
	Record   []byte   // Record is the payload of the record
}

// Append adds a record to a stream
func (c *Client) Append(stream string, record []byte) error {
	return c.sendMessage(encodeAppendByName(stream, record))
}

// ContinuousQuery runs a streaming query against a stream or the key-value store.
func (c *Client) ContinuousQuery(ctx context.Context, dql string, handler func(*Record)) error {
	return c.runConsumer(ctx, newQueryConsumer(c, c.nextConsumerID(), dql, true, nil, handler))
}

// ContinuousQuery runs a streaming query against a stream or the key-value store.
// Use options to configure the behavior of the query, such as the first Record to query.
// c.f. QueryOptions for a more details.
func (c *Client) ContinuousQueryOptions(ctx context.Context, dql string, options *QueryOptions, handler func(*Record)) error {
	return c.runConsumer(ctx, newQueryConsumer(c, c.nextConsumerID(), dql, true, options, handler))
}

// ListStreams iterates all keys currently in the system.
// This operation returns a finite amount results.
// This operation can be canceled by cancelling the context.
func (c *Client) ListKeys(ctx context.Context, keyPattern string, handler func(string)) error {
	return c.runConsumer(ctx, newListConsumer(c, c.nextConsumerID(), false, keyPattern, handler))
}

// ListStreams iterates all streams in Driveline.
// This operation returns a finite amount results.
// This operation can be canceled by cancelling the context.
func (c *Client) ListStreams(ctx context.Context, streamPattern string, handler func(string)) error {
	return c.runConsumer(ctx, newListConsumer(c, c.nextConsumerID(), true, streamPattern, handler))
}

// Query runs a simple -- one shot -- query against a stream or a subset of the key-value store.
// Although this operation can generate large amounts of data, it will terminate.
func (c *Client) Query(ctx context.Context, dql string, handler func(*Record)) error {
	return c.runConsumer(ctx, newQueryConsumer(c, c.nextConsumerID(), dql, false, nil, handler))
}

// Query runs a simple -- one shot -- query against a stream or a subset of the key-value store.
// Use options to configure the behavior of the query, such as the first Record to query.
// c.f. QueryOptions for a more details.
func (c *Client) QueryOptions(ctx context.Context, dql string, options *QueryOptions, handler func(*Record)) error {
	return c.runConsumer(ctx, newQueryConsumer(c, c.nextConsumerID(), dql, false, options, handler))
}

// Remove deletes a key from the key-value store.
func (c *Client) Remove(key string) error {
	return c.sendMessage(encodeRemove(key))
}

// RemoveMatches deletes all keys matching the provided pattern from the key-value store.
func (c *Client) RemoveMatches(keyPattern string) error {
	return c.sendMessage(encodeRemoveMatches(keyPattern))
}

// Store writes data to the key-value store.
func (c *Client) Store(key string, record []byte) error {
	return c.sendMessage(encodeStore(key, record, nil))
}

// Store writes data to the key-value store.
// Use options to configure TTL and CAS.
func (c *Client) StoreOptions(key string, record []byte, options *StoreOptions) error {
	return c.sendMessage(encodeStore(key, record, options))
}

// Sync execute a sync cycle with the server.
func (c *Client) Sync(ctx context.Context) error {
	return c.runConsumer(ctx, newSyncConsumer(c, c.nextConsumerID()))
}

// Truncate removes all records of the specified stream.
func (c *Client) Truncate(stream string) error {
	return c.sendMessage(encodeTruncateByName(stream))
}

func (c *Client) append(streamID streamID, record []byte) error {
	if streamID.isNumeric() {
		return c.sendMessage(encodeAppendByID(streamID.numericID(), record))
	}
	return c.sendMessage(encodeAppendByName(streamID.textualID(), record))
}

func (c *Client) query(consumer *queryConsumer) error {
	return c.sendMessage(encodeQuery(consumer.isContinuous, consumer.ConsumerID, consumer.dql, consumer.options))
}

func (c *Client) truncate(streamID streamID) error {
	if streamID.isNumeric() {
		return c.sendMessage(encodeTruncateByID(streamID.numericID()))
	}
	return c.sendMessage(encodeTruncateByName(streamID.textualID()))
}

func (c *Client) list(consumer *listConsumer) error {
	return c.sendMessage(encodeList(consumer.isStream, consumer.ConsumerID, consumer.pattern))
}

// Load reads data from the key-value store.
func (c *Client) Load(ctx context.Context, key string) (*Record, error) {
	consumer := newLoadConsumer(c, c.nextConsumerID(), key)
	if err := c.runConsumer(ctx, consumer); err != nil {
		return nil, err
	}
	return consumer.record, nil
}

func (c *Client) load(consumer *loadConsumer) error {
	return c.sendMessage(encodeLoad(consumer.ConsumerID, consumer.key))
}

func (c *Client) sync(consumer *syncConsumer) error {
	return c.sendMessage(encodeSync(consumer.consumerID()))
}

func (c *Client) define(name string, id uint8) error {
	return c.sendMessage(encodeDefine(id, name))
}

func (c *Client) cancel(consumer consumer) error {
	return c.sendMessage(encodeCancel(consumer.consumerID()))
}
