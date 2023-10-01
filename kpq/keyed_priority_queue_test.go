package kpq

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewKeyedPriorityQueue_NilCmp(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("want NewKeyedPriorityQueue to panic when receiving a nil comparison cunction")
		}
	}()

	NewKeyedPriorityQueue[int, int](nil)
}

func TestKeyedPriorityQueue_Push(t *testing.T) {
	pq := NewKeyedPriorityQueue[string](func(x, y int) bool {
		return x < y
	})

	items := []struct {
		key string
		val int
	}{
		{key: "fourth", val: 10},
		{key: "second", val: 8},
		{key: "third", val: 9},
		{key: "first", val: 6},
		{key: "last", val: 20},
	}

	for _, item := range items {
		err := pq.Push(item.key, item.val)
		if err != nil {
			t.Fatalf("Push(%v, %v): got unexpected error %v", item.key, item.val, err)
		}
	}

	gotKey, gotVal, ok := pq.Peek()
	if !ok {
		t.Fatal("got no min value in the priority queue")
	}

	if want := 6; gotVal != want {
		t.Errorf("pq.Peek(): got value %v; want %v", gotVal, want)
	}

	if want := "first"; gotKey != want {
		t.Errorf("pq.Peek(): got key %v; want %v", gotKey, want)
	}
}

func TestKeyedPriorityQueue_Push_Error(t *testing.T) {
	t.Run("KeyAlreadyExists", func(t *testing.T) {
		pq := NewKeyedPriorityQueue[string](func(x, y int) bool {
			return x < y
		})

		k := "key"
		if err := pq.Push(k, 10); err != nil {
			t.Fatalf("pq.Push(%q, 10): got unexpected error %v", k, err)
		}

		err := pq.Push(k, 20)

		var wantErr KeyAlreadyExistsError[string]
		if !errors.As(err, &wantErr) {
			t.Errorf("pq.Push(%q, 20): got error type %T; want it to be %T", k, err, wantErr)
		}
	})
}

func TestKeyedPriorityQueue_Update(t *testing.T) {
	pq := NewKeyedPriorityQueue[string](func(x, y int) bool {
		return x < y
	})

	items := []struct {
		key string
		val int
	}{
		{key: "fourth", val: 10},
		{key: "second", val: 8},
		{key: "third", val: 9},
		{key: "first", val: 6},
		{key: "last", val: 20},
	}

	for _, item := range items {
		err := pq.Push(item.key, item.val)
		if err != nil {
			t.Fatalf("Push(%v, %v): got unexpected error %v", item.key, item.val, err)
		}
	}

	testCases := []struct {
		key           string
		newValue      int
		wantPeekKey   string
		wantPeekValue int
	}{
		{
			key:           "last",
			newValue:      5,
			wantPeekKey:   "last",
			wantPeekValue: 5,
		},
		{
			key:           "third",
			newValue:      25,
			wantPeekKey:   "last",
			wantPeekValue: 5,
		},
		{
			key:           "first",
			newValue:      7,
			wantPeekKey:   "last",
			wantPeekValue: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%d", tc.key, tc.newValue), func(t *testing.T) {
			err := pq.Update(tc.key, tc.newValue)
			if err != nil {
				t.Fatalf("pq.Update(%q, %d): got unexpected error: %v", tc.key, tc.newValue, err)
			}

			gotPeekKey, ok := pq.PeekKey()
			if !ok {
				t.Fatal("got no min key in the priority queue")
			}

			if gotPeekKey != tc.wantPeekKey {
				t.Errorf("pq.PeekKey(): got %q; want %q", gotPeekKey, tc.wantPeekKey)
			}

			gotMinVal, ok := pq.PeekValue()
			if !ok {
				t.Fatal("got no min value in the priority queue")
			}

			if gotMinVal != tc.wantPeekValue {
				t.Errorf("pq.PeekValue(): got %d; want %d", gotMinVal, tc.wantPeekValue)
			}
		})
	}
}

func TestKeyedPriorityQueue_Update_Error(t *testing.T) {
	t.Run("KeyNotFound", func(t *testing.T) {
		pq := NewKeyedPriorityQueue[string](func(x, y int) bool {
			return x < y
		})

		k := "key-not-found"
		err := pq.Update(k, 20)

		var wantErr KeyNotFoundError[string]
		if !errors.As(err, &wantErr) {
			t.Errorf("pq.Update(%q, 20): got error type %T; want it to be %T", k, err, wantErr)
		}
	})
}

func TestKeyedPriorityQueue_Pop(t *testing.T) {
	t.Run("Keys", func(t *testing.T) {
		pq := NewKeyedPriorityQueue[string](func(x, y int) bool {
			return x < y
		})

		items := []struct {
			key string
			val int
		}{
			{key: "fourth", val: 10},
			{key: "second", val: 8},
			{key: "third", val: 9},
			{key: "first", val: 6},
			{key: "last", val: 20},
		}

		for _, item := range items {
			err := pq.Push(item.key, item.val)
			if err != nil {
				t.Fatalf("Push(%v, %v): got unexpected error %v", item.key, item.val, err)
			}
		}

		testCases := []struct {
			wantKey       string
			wantValue     int
			wantPeekKey   string
			wantPeekValue int
			wantLen       int
		}{
			{
				wantKey:       "first",
				wantValue:     6,
				wantPeekKey:   "second",
				wantPeekValue: 8,
				wantLen:       4,
			},
			{
				wantKey:       "second",
				wantValue:     8,
				wantPeekKey:   "third",
				wantPeekValue: 9,
				wantLen:       3,
			},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s_%d", tc.wantKey, tc.wantValue), func(t *testing.T) {
				gotKey, gotValue, ok := pq.Pop()
				if !ok {
					t.Fatal("pq.Pop(): got unexpected empty prioriy queue")
				}

				if gotKey != tc.wantKey {
					t.Errorf("pq.Pop(): got key %q; want %q", gotKey, tc.wantKey)
				}

				if gotValue != tc.wantValue {
					t.Errorf("pq.Pop(): got value %d; want %d", gotValue, tc.wantValue)
				}

				gotPeekKey, gotPeekValue, ok := pq.Peek()
				if !ok {
					t.Fatal("got no min key and value in the priority queue")
				}

				if gotPeekKey != tc.wantPeekKey {
					t.Errorf("pq.Peek(): got key %q; want %q", gotPeekKey, tc.wantPeekKey)
				}

				if gotPeekValue != tc.wantPeekValue {
					t.Errorf("pq.Peek(): got value %d; want %d", gotPeekValue, tc.wantPeekValue)
				}

				if got := pq.Len(); got != tc.wantLen {
					t.Errorf("pq.Len(): got %d; want %d", got, tc.wantLen)
				}
			})
		}
	})

	t.Run("EmptyPQ", func(t *testing.T) {
		pq := NewKeyedPriorityQueue[string](func(x, y int) bool {
			return x < y
		})

		_, _, ok := pq.Pop()
		if ok {
			t.Errorf("pq.Pop(): got unexpected non-empty priorit queue")
		}
	})
}

func TestKeyedPriorityQueue_Remove(t *testing.T) {
	pq := NewKeyedPriorityQueue[string](func(x, y int) bool {
		return x < y
	})

	items := []struct {
		key string
		val int
	}{
		{key: "fourth", val: 10},
		{key: "second", val: 8},
		{key: "third", val: 9},
		{key: "first", val: 6},
		{key: "last", val: 20},
	}

	for _, item := range items {
		err := pq.Push(item.key, item.val)
		if err != nil {
			t.Fatalf("Push(%v, %v): got unexpected error %v", item.key, item.val, err)
		}
	}

	t.Run("Keys", func(t *testing.T) {
		testCases := []struct {
			key           string
			wantPeekKey   string
			wantPeekValue int
			wantLen       int
		}{
			{
				key:           "first",
				wantPeekKey:   "second",
				wantPeekValue: 8,
				wantLen:       4,
			},
			{
				key:           "third",
				wantPeekKey:   "second",
				wantPeekValue: 8,
				wantLen:       3,
			},
			{
				key:           "second",
				wantPeekKey:   "fourth",
				wantPeekValue: 10,
				wantLen:       2,
			},
			{
				key:           "last",
				wantPeekKey:   "fourth",
				wantPeekValue: 10,
				wantLen:       1,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.key, func(t *testing.T) {
				pq.Remove(tc.key)

				gotPeekKey, gotPeekValue, ok := pq.Peek()
				if !ok {
					t.Fatal("got no min key and value in the priority queue")
				}

				if gotPeekKey != tc.wantPeekKey {
					t.Errorf("pq.Peek(): got key %q; want %q", gotPeekKey, tc.wantPeekKey)
				}

				if gotPeekValue != tc.wantPeekValue {
					t.Errorf("pq.PeekValue(): got value %d; want %d", gotPeekValue, tc.wantPeekValue)
				}

				if got := pq.Len(); got != tc.wantLen {
					t.Errorf("pq.Len(): got %d; want %d", got, tc.wantLen)
				}
			})
		}
	})

	t.Run("NonExistingKey", func(t *testing.T) {
		want := pq.Len()
		pq.Remove("non-existing-key")

		if got := pq.Len(); got != want {
			t.Errorf("pq.Len(): got %d; want %d", got, want)
		}
	})
}

func TestKeyedPriorityQueue_Peek_EmptyQueue(t *testing.T) {
	pq := NewKeyedPriorityQueue[int](func(x, y int) bool { return x < y })

	_, _, ok := pq.Peek()
	if ok {
		t.Error("pq.Peek(): got unexpected non empty priority queue")
	}
}

func TestKeyedPriorityQueue_PeekValue_EmptyQueue(t *testing.T) {
	pq := NewKeyedPriorityQueue[int](func(x, y int) bool { return x < y })

	_, ok := pq.PeekValue()
	if ok {
		t.Error("pq.PeekValue(): got unexpected non empty priority queue")
	}
}

func TestKeyedPriorityQueue_PeekKey_EmptyQueue(t *testing.T) {
	pq := NewKeyedPriorityQueue[int](func(x, y int) bool { return x < y })

	_, ok := pq.PeekKey()
	if ok {
		t.Error("pq.PeekKey(): got unexpected non empty priority queue")
	}
}

func TestKeyedPriorityQeue_Contains(t *testing.T) {
	pq := NewKeyedPriorityQueue[string](func(x, y int) bool { return x < y })

	k := "user"
	if err := pq.Push(k, 10); err != nil {
		t.Fatalf("pq.Push(%q, 10): got unexpected error %v", k, err)
	}

	if ok := pq.Contains(k); !ok {
		t.Errorf("pq.Contains(%q): got no key in priority queue", k)
	}
}

func TestKeyedPriorityQueue_ValueOf(t *testing.T) {
	pq := NewKeyedPriorityQueue[string](func(x, y int) bool { return x < y })

	k, v := "user", 10
	if err := pq.Push(k, v); err != nil {
		t.Fatalf("pq.Push(%q, 10): got unexpected error %v", k, err)
	}

	t.Run("ExistingKey", func(t *testing.T) {
		got, ok := pq.ValueOf(k)
		if !ok {
			t.Fatalf("pq.ValueOf(%q): got no key %q in priority queue", k, k)
		}

		if want := v; got != want {
			t.Errorf("pq.ValueOf(%q): got unexpected value %v for key %q; want %v", k, got, k, want)
		}
	})

	t.Run("NonExistingKey", func(t *testing.T) {
		k := "non-existing-key"
		_, ok := pq.ValueOf(k)
		if ok {
			t.Errorf("pq.Contains(%q): got unexpected key %q in priority queue", k, k)
		}
	})
}

func TestKeyedPriorityQueue_IsEmpty(t *testing.T) {
	pq := NewKeyedPriorityQueue[string](func(x, y int) bool {
		return x < y
	})

	k := "key"
	if err := pq.Push(k, 10); err != nil {
		t.Fatalf("pq.Push(%q, 10): got unexpected error %v", k, err)
	}

	if empty := pq.IsEmpty(); empty {
		t.Fatal("pq.IsEmpty(): got unexpected empty priority queue")
	}

	pq.Remove(k)

	if empty := pq.IsEmpty(); !empty {
		t.Fatal("pq.IsEmpty(): got unexpected non-empty priority queue")
	}
}

func benchmarkKeyedPriorityQueue_PushPop(b *testing.B, n int) {
	pq := NewKeyedPriorityQueue[int](func(a, b int) bool {
		return a > b
	})
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			pq.Push(j, i)
		}

		for !pq.IsEmpty() {
			pq.Pop()
		}
	}
}

func BenchmarkKeyedPriorityQueue_PushPop_10(b *testing.B) {
	benchmarkKeyedPriorityQueue_PushPop(b, 10)
}
func BenchmarkKeyedPriorityQueue_PushPop_100(b *testing.B) {
	benchmarkKeyedPriorityQueue_PushPop(b, 100)
}
func BenchmarkKeyedPriorityQueue_PushPop_1000(b *testing.B) {
	benchmarkKeyedPriorityQueue_PushPop(b, 1000)
}
func BenchmarkKeyedPriorityQueue_PushPop_10000(b *testing.B) {
	benchmarkKeyedPriorityQueue_PushPop(b, 10000)
}
func BenchmarkKeyedPriorityQueue_PushPop_100000(b *testing.B) {
	benchmarkKeyedPriorityQueue_PushPop(b, 100000)
}
func BenchmarkKeyedPriorityQueue_PushPop_1000000(b *testing.B) {
	benchmarkKeyedPriorityQueue_PushPop(b, 1000000)
}
