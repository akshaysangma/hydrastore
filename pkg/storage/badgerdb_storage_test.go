package storage

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dgraph-io/badger/v4"
)

// --- Test Helpers ---

func createTempBadgerDBStorage(t *testing.T) (*BadgerDBStorage, string) {
	dir := t.TempDir()
	store, err := NewBadgerDBStorage(dir)
	if err != nil {
		t.Fatalf("Error creating BadgerDB storage: %v", err)
	}
	return store, dir
}

func TestBadgerDBStorage_SetAndGet(t *testing.T) {
	store, _ := createTempBadgerDBStorage(t)
	defer store.Close()

	key := []byte("test-key")
	value := []byte("test-value")

	err := store.Set(key, value)
	if err != nil {
		t.Fatalf("Error during Set: %v", err)
	}

	result, err := store.Get(key)
	if err != nil {
		t.Fatalf("Error during Get: %v", err)
	}

	if !bytes.Equal(result, value) {
		t.Fatalf("Retrieved value mismatch. Expected: %v, got: %v", value, result)
	}
}

func TestBadgerDBStorage_Delete(t *testing.T) {
	store, _ := createTempBadgerDBStorage(t)
	defer store.Close()

	key := []byte("test-key")
	value := []byte("test-value")

	// Set the value first
	err := store.Set(key, value)
	if err != nil {
		t.Fatalf("Error during Set: %v", err)
	}

	// Delete the value
	err = store.Delete(key)
	if err != nil {
		t.Fatalf("Error during Delete: %v", err)
	}

	// Attempt to retrieve the deleted value
	result, err := store.Get(key)

	// Expect either 'not found' error or nil  result (BadgerDB behavior may change the specific signal on deleted items)
	if !errors.Is(err, badger.ErrKeyNotFound) && result != nil {
		t.Fatalf("Expected key to be deleted")
	}
}

func TestBadgerDBStorage_ErrorHandling(t *testing.T) {
	// Simulate an inaccessible path
	tempDir := "/nonexistent"
	store, err := NewBadgerDBStorage(tempDir)
	if err == nil {
		t.Fatalf("Expected storage creation to fail due to access issues")
	}

	// Even on creation errors, close if store isn't nil (might become important later)
	if store != nil {
		store.Close()
	}
}

// Large Value test
func TestBadgerDBStorage_LargeValue(t *testing.T) {
	store, _ := createTempBadgerDBStorage(t)
	defer store.Close()

	largeValue := make([]byte, 1024*1024) // 1MB

	key := []byte("large-key")

	err := store.Set(key, largeValue)
	if err != nil {
		t.Fatalf("Storing large value encountered an error: %v", err)
	}

	result, err := store.Get(key)
	if err != nil {
		t.Fatalf("Error retrieving large value: %v", err)
	}

	// Check values match if retrieval succeeded
	if !bytes.Equal(result, largeValue) {
		t.Fatal("Retrieved large value doesn't match expected")
	}
}

// Large Key test
func TestBadgerDBStorage_LargeKey(t *testing.T) {
	store, _ := createTempBadgerDBStorage(t)
	defer store.Close()

	largeKey := make([]byte, 65000) //default key limit in BadgerDB
	value := []byte("large-value")

	err := store.Set(largeKey, value)
	if err != nil {
		t.Fatalf("Storing value with large key encountered an error: %v", err)
	}

	result, err := store.Get(largeKey)
	if err != nil {
		t.Fatalf("Error retrieving large key's value: %v", err)
	}

	// Check values match if retrieval succeeded
	if !bytes.Equal(result, value) {
		t.Fatal("Retrieved value from large key doesn't match expected")
	}
}

func TestBadgerDBStorage_LargeKeyLargeValue(t *testing.T) {

	store, _ := createTempBadgerDBStorage(t)
	defer store.Close()

	// default opts key limit of BadgerDB.
	largeKey := make([]byte, 65000)

	largeValue := make([]byte, 1024*1024)

	err := store.Set(largeKey, largeValue)
	if err != nil {
		t.Fatalf("Storing large value with large key encountered an error: %v", err)
	}

	result, err := store.Get(largeKey)
	if err != nil {
		t.Fatalf("Error retrieving large key's large value: %v", err)
	}

	// Check values match if retrieval succeeded
	if !bytes.Equal(result, largeValue) {
		t.Fatalf("Retrieved large value from large key doesn't match expected")
	}
}
