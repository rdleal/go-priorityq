# Keyed Priority Queue

[![Go Reference](https://pkg.go.dev/badge/github.com/rdleal/go-priorityq/kpq.svg)](https://pkg.go.dev/github.com/rdleal/go-priorityq/kpq)
[![Go Report Card](https://goreportcard.com/badge/github.com/rdleal/go-priorityq)](https://goreportcard.com/report/github.com/rdleal/go-priorityq)
[![codecov](https://codecov.io/gh/rdleal/go-priorityq/graph/badge.svg?token=DEVXQHRRQD)](https://codecov.io/gh/rdleal/go-priorityq)

A keyed priority queue is a data structure that allows you to associate keys with priority values
and efficiently retrieve, update, and remove elements based on their priorities.
This package offers concurrent-safe operations that leverages a binary heap to maintain the priority queue.
Operations like Push, Pop, Update and Remove have O(log n) time complexity, where n is the size of the priority queue.
The use of a map ensures fast lookups by key. Operations like Peek, Contains and ValueOf have O(1) time complexity.

# Installing

```
go get github.com/rdleal/go-priorityq
```
# Usage

Importing the package:
```go
import "github.com/rdleal/go-priorityq/kpq"
```

Creating a priority queue with `string` as the key type, and `int` as the priority value type:
```go
// Create a new KeyedPriorityQueue with a custom comparison function
cmp := func(a, b int) bool {
	return a < b
}
pq := kpq.NewKeyedPriorityQueue[string](cmp)
```

Pushing elements into the priority queue:
```go
pq.Push("key1", 42)
pq.Push("key2", 30)
pq.Push("key3", 50)
```

Popping the element with the highest priority: 
```go
key, value, ok := pq.Pop()
if !ok {
	log.Fatal("priority queue is empty")
}
```

Updating the priority of an element:
```go
if err := pq.Update("key3", 20); err != nil {
	log.Fatalf("got unexpected error: %v\n", err)
}
```

Removing an element from the priority queue:
```go
pq.Remove("key1")
```

Checking if a key exists in the priority queue:
```go
exists := pq.Contains("key3")
fmt.Println("Key 'key3' exists:", exists)
```

For more operations, check out the [GoDoc page](https://pkg.go.dev/github.com/rdleal/go-priorityq/kpq).

# Testing

Run unit tests:
```sh
go test -v -cover -race ./...
```

# License
[MIT](./LICENSE)
