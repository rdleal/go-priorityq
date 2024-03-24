package kpq_test

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"text/tabwriter"

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

func ExampleKeyedPriorityQueue_Set() {
	// This code implements Dijkstra's algorithm to find the shortest path in a
	// weighted graph from a source vertex to all other vertices.
	//
	// This example shows how to change the priority in `KeyedPriorityQueue`
	// when needed.
	graph := struct {
		len   int
		edges []int
	}{
		len: 8,
		// edges represents the adjacency matrix of a directed weighted Graph.
		edges: []int{
			0, 5, 0, 0, 9, 0, 0, 8,
			0, 0, 12, 15, 0, 0, 0, 4,
			0, 0, 0, 3, 0, 0, 11, 0,
			0, 0, 0, 0, 0, 0, 9, 0,
			0, 0, 0, 0, 0, 4, 20, 5,
			0, 0, 1, 0, 0, 0, 13, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 7, 0, 0, 6, 0, 0,
		},
	}

	edge := func(u, v int) (weight int) {
		return graph.edges[graph.len*u+v]
	}

	adj := func(v int) []int {
		vertices := make([]int, 0)

		for i := 0; i < graph.len; i++ {
			if weight := edge(v, i); weight > 0 {
				vertices = append(vertices, i)
			}
		}

		return vertices
	}

	src := 0

	distTo := make([]int, graph.len)
	for i := 0; i < graph.len; i++ {
		distTo[i] = math.MaxInt
	}
	distTo[src] = 0

	// cmpFunc maintains the variant of a min priority queue,
	// needed for relaxing all the edges from the source.
	cmpFunc := func(a, b int) bool {
		return a < b
	}
	pq := kpq.NewKeyedPriorityQueue[int](cmpFunc)
	pq.Push(src, 0) // starts with source vertex.

	for !pq.IsEmpty() {
		u, dist, _ := pq.Pop()
		// Iterate over vertices adjacent to vertex u, and relax each edge
		// between them.
		// Given a vertex u and v and a weighted edge e from u to v,
		// the relaxation algorithm updates the value in the priority queue
		// if the edge e provides a shorter path from u to v than previously known.
		for _, v := range adj(u) {
			weight := edge(u, v)
			if distTo[v] > dist+weight {
				distTo[v] = dist + weight
				pq.Set(v, distTo[v])
			}
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "Vertex\tDistance From Source")
	for i := 0; i < graph.len; i++ {
		fmt.Fprintf(w, "%3d\t%10d\n", i, distTo[i])
	}
	w.Flush()
	// Output:
	// Vertex    Distance From Source
	//   0                0
	//   1                5
	//   2               14
	//   3               17
	//   4                9
	//   5               13
	//   6               25
	//   7                8
}
