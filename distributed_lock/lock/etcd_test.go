package lock

import (
	"context"
	"testing"
)

func TestNewEtcdClient(t *testing.T) {
	// Create a new etcd client
	client := NewEtcdClient()

	// Check that the client is not nil
	if client == nil {
		t.Errorf("NewEtcdClient() returned nil")
	}
	defer client.Close()
	resp, err := client.Client.Put(context.Background(), "test1", "1")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
	resp1, err := client.Client.Get(context.Background(), "test1")
	if err != nil {
		t.Error(err)
	}
	t.Log(resp1)
}
