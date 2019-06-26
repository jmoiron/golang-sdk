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

// +build integration

package driveline

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestConnectionContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	_, err := NewClient(ctx, "ws://127.0.0.1:8080")
	if err != context.DeadlineExceeded {
		t.Fatalf("should have received an errror")
	}
}

func TestConnectionRetryCount(t *testing.T) {
	_, err := NewClient(context.Background(), "ws://127.0.0.1:32000", MaxReconnect(1))
	if err != ErrMaxReconnect {
		t.Fatalf("should have received an errror")
	}
}

func TestConnection(t *testing.T) {
	c, err := NewClient(context.Background(), "ws://127.0.0.1:8080")
	if err != nil {
		t.Fatalf("cannot create a client")
	}
	c.Close()
}

func TestWriteLoad(t *testing.T) {
	c, err := NewClient(context.Background(), "ws://127.0.0.1:8080")
	if err != nil {
		t.Fatalf("client cannot connect")
	}
	defer c.Close()

	expected := []byte{1, 2, 4}
	if err := c.Store("/test", expected); err != nil {
		t.Fatalf("cannot store data")
	}

	r, err := c.Load(context.Background(), "/test")
	if err != nil {
		t.Fatalf("cannot load data")
	}

	actual := r.Record
	if bytes.Compare(expected, actual) != 0 {
		t.Fatalf("loaded different data than was previously stored")
	}

}

func TestClient_ContinuousQuery(t *testing.T) {
	key := "/dl/infra/infra-manager/create-task/094134ab-97a6-4480-b8ab-ecc51e478b6d/ba9b63ea-15f2-4fb5-9f40-b6a84dadfe76/response"
	dql := fmt.Sprintf("SELECT * FROM '%s'", key)
	c, err := NewClient(context.Background(), "ws://127.0.0.1:8080")
	if err != nil {
		t.Fatalf("client cannot connect")
	}
	defer c.Close()
	c.Remove(key)
	expected := []byte{1, 2, 4}

	ctx, cancel := context.WithCancel(context.Background())

	var qErr error
	go func() {
		qErr = c.ContinuousQuery(ctx, dql, func(rec *Record) {
			actual := rec.Record
			if bytes.Compare(expected, actual) != 0 {
				t.Fatalf("records don't match!")
			}
		})
		if err != nil {
			t.Fatalf("query failed!")
		}
	}()

	if err := c.Store(key, expected); err != nil {
		t.Fatalf("cannot store key")
	}
	cancel()
	if qErr != nil {
		t.Fatalf("request failed:%s", qErr.Error())
	}

}

/*
func TestClient_QueryFuture(t *testing.T) {
	key := "/dl/infra/infra-manager/create-task/094134ab-97a6-4480-b8ab-ecc51e478b6d/ba9b63ea-15f2-4fb5-9f40-b6a84dadfe76/response"
	// dql := fmt.Sprintf("SELECT * FROM '%s'", key)
	c, err := NewClient(context.Background(), "ws://127.0.0.1:8080")
	if err != nil {
		t.Fatalf("client cannot connect")
	}
	defer c.Close()
	c.Remove(key)
	expected := []byte{1, 2, 4}
	f, err := c.QueryFuture(key)
	if err != nil {
		t.Fatalf("cannot create future query")
	}

	go func() {
		if err := c.Store(key, expected); err != nil {
			t.Fatalf("error while storing data %v", err)
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	rec, err := f.GetRecordContext(ctx)
	if err != nil {
		t.Fatalf("cannot get record %v", err)
	}

	actual := rec.Record

	if bytes.Compare(expected, actual) != 0 {
		t.Fatalf("records don't match!")
	}

}
*/
