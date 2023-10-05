package kpq_test

import (
	"errors"
	"fmt"
	"log"

	"github.com/rdleal/go-priorityq/kpq"
)

func ExampleKeyAlreadyExistsError() {
	pq := kpq.NewKeyedPriorityQueue[string](func(a, b int) bool {
		return a > b // max priority queue
	})

	key := "key1"
	pq.Push(key, 10)
	err := pq.Push(key, 20) // pushing the same key should return an error
	if err != nil {
		var keyErr kpq.KeyAlreadyExistsError[string]
		if errors.As(err, &keyErr) {
			fmt.Println(keyErr.Key() == key)
		}
	}
	// Output:
	// true
}

func Example() {
	// Create a new KeyedPriorityQueue with a custom comparison function
	cmp := func(a, b int) bool {
		return a < b
	}
	pq := kpq.NewKeyedPriorityQueue[string](cmp)

	// Insert elements onto the priority queue
	pq.Push("key1", 42)
	pq.Push("key2", 30)
	pq.Push("key3", 50)

	// Remove and retrieve the element with the highest priority
	key, value, ok := pq.Pop()
	if !ok {
		log.Fatal("priority queue is empty")
	}

	fmt.Printf("Key: %q, Value: %d\n", key, value)

	// Update the priority value of an element
	if err := pq.Update("key3", 20); err != nil {
		log.Fatalf("got unexpected error: %v\n", err)
	}

	// Remove an element from the priority queue
	pq.Remove("key1")

	// Check if a key exists in the priority queue
	exists := pq.Contains("key3")

	fmt.Println("Key 'key3' exists:", exists)
	// Output:
	// Key: "key2", Value: 30
	// Key 'key3' exists: true
}
