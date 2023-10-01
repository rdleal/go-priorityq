// Package kpq provides a generic keyed priority queue implementation.
//
// A keyed priority queue is a data structure that allows you to associate keys with priority values
// and efficiently retrieve, update, and remove elements based on their priorities.
// This package offers concurrent-safe operations that leverages a binary heap to maintain the priority queue.
// Operations like Push, Pop, Update and Remove have O(log n) time complexity, where n is the size of the priority queue.
// The use of a map ensures fast lookups by key. Operations like Peek, Contains and ValueOf have O(1) time complexity.
package kpq

import (
	"fmt"
	"sync"
)

type keyError[K comparable] struct {
	key K
	msg string // description of the error
}

// Error returns the textual description of an error.
func (e keyError[K]) Error() string {
	return e.msg
}

// Key returns the key that caused the error.
func (e keyError[K]) Key() K {
	return e.key
}

// KeyAlreadyExistsError represents an error from calling a Push method
// with a key that already exists in the priority queue.
type KeyAlreadyExistsError[K comparable] struct {
	keyError[K]
}

func newKeyAlreadyExistsError[K comparable](k K) error {
	return KeyAlreadyExistsError[K]{
		keyError[K]{
			key: k,
			msg: fmt.Sprintf("keyed priority queue: key \"%v\" already exists", k),
		},
	}
}

// KeyNotFoundError represents an error from calling Update method
// with a key that doesn't exist in the priority queue.
type KeyNotFoundError[K comparable] struct {
	keyError[K]
}

func newKeyNotFoundError[K comparable](k K) error {
	return KeyNotFoundError[K]{
		keyError[K]{
			key: k,
			msg: fmt.Sprintf("keyed priority queue: key \"%v\" does not exist", k),
		},
	}
}

// CmpFunc is a generic function type used for ordering the priority queue.
type CmpFunc[V any] func(x, y V) bool

// KeyedPriorityQueue represents a generic keyed priority queue,
// where K is the key type and V is the priority value type.
//
// KeyedPriorityQueue must not be copied after first use.
type KeyedPriorityQueue[K comparable, V any] struct {
	mu sync.RWMutex

	pm   []K       // position map
	im   map[K]int // inverse map of pm; note that for a given key k, pm[im[k]] == k
	vals map[K]V   // generic priority values of key k
	cmp  CmpFunc[V]
}

// NewKeyedPriorityQueue returns a new keyed priority queue
// that uses the given cmp function for ordering the priority queue.
//
// NewKeyedPriorityQueue will panic if cmp is nil.
func NewKeyedPriorityQueue[K comparable, V any](cmp CmpFunc[V]) *KeyedPriorityQueue[K, V] {
	if cmp == nil {
		panic("keyed priority queue: comparison function cannot be nil")
	}
	return &KeyedPriorityQueue[K, V]{
		pm:   make([]K, 0),
		im:   make(map[K]int),
		vals: make(map[K]V),
		cmp:  cmp,
	}
}

// Push inserts the given priority value v onto the priority queue associated with the given key k.
// If the key already exists in the priority queue, it returns a KeyAlreadyExistsError error.
func (pq *KeyedPriorityQueue[K, V]) Push(k K, v V) error {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if _, ok := pq.im[k]; ok {
		return newKeyAlreadyExistsError(k)
	}

	n := len(pq.pm)
	pq.pm = append(pq.pm, k)
	pq.im[k] = n
	pq.vals[k] = v
	pq.swim(n)
	return nil
}

// Pop removes and returns the highest priority key and value from the priority queue.
// It returns false as its last return value if the priority queue is empty; otherwise, true.
func (pq *KeyedPriorityQueue[K, V]) Pop() (K, V, bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.pm) == 0 {
		var k K
		var v V
		return k, v, false
	}
	n := len(pq.pm) - 1
	k := pq.pm[0]
	v := pq.vals[k]
	pq.swap(0, n)
	pq.sink(0, n)
	pq.pm = pq.pm[:n]
	delete(pq.im, k)
	delete(pq.vals, k)
	return k, v, true
}

// Update changes the priority value associated with the given key k to the given value v.
// If there's no key k in the priority queue, it returns a KeyNotFoundError error.
func (pq *KeyedPriorityQueue[K, V]) Update(k K, v V) error {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	i, ok := pq.im[k]
	if !ok {
		return newKeyNotFoundError(k)
	}
	pq.vals[k] = v
	pq.swim(i)
	pq.sink(i, len(pq.vals))
	return nil
}

// Peek returns the highest priority key and value from the priority queue.
// It returns false as its last return value if the priority queue is empty; otherwise, true.
func (pq *KeyedPriorityQueue[K, V]) Peek() (K, V, bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if len(pq.pm) == 0 {
		var k K
		var v V
		return k, v, false
	}
	return pq.pm[0], pq.vals[pq.pm[0]], true
}

// PeekKey returns the highest priority key from the priority queue.
// It returns false as its last return value if the priority queue is empty; otherwise, true.
func (pq *KeyedPriorityQueue[K, V]) PeekKey() (K, bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if len(pq.pm) == 0 {
		var k K
		return k, false
	}
	return pq.pm[0], true
}

// PeekValue returns the highest priority value from the priority queue.
// It returns false as its last return value if the priority queue is empty; otherwise, true.
func (pq *KeyedPriorityQueue[K, V]) PeekValue() (V, bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if len(pq.pm) == 0 {
		var v V
		return v, false
	}
	return pq.vals[pq.pm[0]], true
}

// Contains returns true if the given key k exists in the priority queue; otherwise, false.
func (pq *KeyedPriorityQueue[K, V]) Contains(k K) bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	_, ok := pq.im[k]
	return ok
}

// ValueOf returns the priority value associated with the given key k.
// It returns false as its last return value if there's no such key k
// in the priority queue; otherwise, true.
func (pq *KeyedPriorityQueue[K, V]) ValueOf(k K) (V, bool) {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	v, ok := pq.vals[k]
	return v, ok
}

// Remove removes the priority value associated with the given key k from the priority queue.
// It's a no-op if there's no such key k in the priority queue.
func (pq *KeyedPriorityQueue[K, V]) Remove(k K) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	i, ok := pq.im[k]
	if !ok {
		return
	}
	n := len(pq.pm) - 1
	if i != n {
		pq.swap(i, n)
		pq.sink(i, n)
		pq.swim(i)
	}
	pq.pm = pq.pm[:n]
	delete(pq.im, k)
	delete(pq.vals, k)
}

// Len returns the size of the priority queue.
func (pq *KeyedPriorityQueue[K, V]) Len() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return len(pq.pm)
}

// IsEmpty returns true if the priority queue is empty; otherwise, false.
func (pq *KeyedPriorityQueue[K, V]) IsEmpty() bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return len(pq.pm) == 0
}

func (pq *KeyedPriorityQueue[K, V]) swap(i, j int) {
	pq.pm[i], pq.pm[j] = pq.pm[j], pq.pm[i]
	pq.im[pq.pm[i]], pq.im[pq.pm[j]] = i, j
}

func (pq *KeyedPriorityQueue[K, V]) swim(i int) {
	for i > 0 && pq.compare(i, parent(i)) {
		pq.swap(i, parent(i))
		i = parent(i)
	}
}

func (pq *KeyedPriorityQueue[K, V]) sink(i, n int) {
	for leftChild(i) < n {
		j := leftChild(i)
		if j < 0 { // j < 0 after int overflow
			break
		}
		if r := j + 1; r < n && pq.compare(r, j) {
			j = r // r == j + 1 == right child
		}
		if !pq.compare(j, i) {
			break
		}
		pq.swap(i, j)
		i = j
	}
}

func (pq *KeyedPriorityQueue[K, V]) compare(i, j int) bool {
	return pq.cmp(pq.vals[pq.pm[i]], pq.vals[pq.pm[j]])
}

func leftChild(i int) int {
	return (i * 2) + 1
}

func parent(i int) int {
	return (i - 1) / 2
}
