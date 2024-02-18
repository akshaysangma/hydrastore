package storage

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dgraph-io/badger/v4"
)

const (
	valueSize          = 1024 * 1024
	benchmarkKeyPrefix = "benchmark-key-"
)

// createTempBadgerDBStorage is a Helper
func createTempBadgerDBStorage(dir string) (*BadgerDBStorage, error) {
	store, err := NewBadgerDBStorage(dir)
	return store, err
}

func TestBadgerDBStorage_SetAndGet(t *testing.T) {
	store, _ := createTempBadgerDBStorage(t.TempDir())
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
	store, _ := createTempBadgerDBStorage(t.TempDir())
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
	store, _ := createTempBadgerDBStorage(t.TempDir())
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
	store, _ := createTempBadgerDBStorage(t.TempDir())
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

	store, _ := createTempBadgerDBStorage(t.TempDir())
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

func BenchmarkBadgerDBStorage_Set(b *testing.B) {
	store, _ := createTempBadgerDBStorage(b.TempDir())
	defer store.Close()

	key := []byte(benchmarkKeyPrefix + "set")
	value := make([]byte, valueSize)

	b.ResetTimer() // Reset the timer right before the benchmark loop

	for i := 0; i < b.N; i++ {
		err := store.Set(key, value)
		if err != nil {
			b.Errorf("Error during Set benchmark: %v", err)
		}
	}
}

func BenchmarkBadgerDBStorage_Get(b *testing.B) {
	store, _ := createTempBadgerDBStorage(b.TempDir())
	defer store.Close()

	key := []byte(benchmarkKeyPrefix + "get")
	value := make([]byte, valueSize)

	// Pre-store the value before starting measurement
	err := store.Set(key, value)
	if err != nil {
		b.Fatalf("Error preparing storage for Get benchmark: %v", err)
	}

	b.ResetTimer() // Reset timer before Get operations

	for i := 0; i < b.N; i++ {
		_, err := store.Get(key)
		if err != nil {
			b.Errorf("Error during Get benchmark: %v", err)
		}
	}
}
