package poller

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestHeap(t *testing.T) {
	h := &Heap{}

	num := 100
	for i := 0; i < num; i++ {
		h.Push(HeapPollInfo{
			Info: HTTPPollInfo{
				Interval: time.Second * time.Duration(8),
				ID:       fmt.Sprintf("%v", i),
			},
			Next: time.Now().Add(time.Second * time.Duration(i)), // time.Duration((i%2)*-1)
		})
	}

	for i := 0; i < num; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if val.Info.ID != fmt.Sprintf("%v", i) {
			t.Errorf("expected pop ID %v got %v next %v", i, val.Info.ID, val.Next)
		}
	}
}

func TestHeapRandom(t *testing.T) {
	h := &Heap{}

	num := 10
	for i := 0; i < num; i++ {
		h.Push(HeapPollInfo{
			Info: HTTPPollInfo{
				Interval: time.Second * time.Duration(8),
				ID:       fmt.Sprintf("%v", i),
			},
			Next: time.Now().Add(time.Duration(rand.Int63())),
		})
	}

	previousTime := time.Now()
	for i := 0; i < num; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}
}

func TestHeapRandomPopping(t *testing.T) {
	h := &Heap{}

	randInfo := func(id int) HeapPollInfo {
		return HeapPollInfo{
			Info: HTTPPollInfo{
				Interval: time.Second * time.Duration(8),
				ID:       fmt.Sprintf("%v", id),
			},
			Next: time.Now().Add(time.Duration(rand.Int63())),
		}
	}

	num := 10
	for i := 0; i < num; i++ {
		h.Push(randInfo(i))
	}

	previousTime := time.Now()
	for i := 0; i < num/2; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}

	for i := 0; i < num; i++ {
		h.Push(randInfo(i))
	}
	val, ok := h.Pop()
	if !ok {
		t.Errorf("expected pop, got empty heap")
	} else {
		previousTime = val.Next
	}

	for i := 0; i < num; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}

	for i := 0; i < num; i++ {
		h.Push(randInfo(i))
	}
	val, ok = h.Pop()
	if !ok {
		t.Errorf("expected pop, got empty heap")
	} else {
		previousTime = val.Next
	}

	for i := 0; i < num; i++ {
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}

	for i := 0; i < num/2-2; i++ { // -2 for the two we manually popped in order to get the max
		val, ok := h.Pop()
		if !ok {
			t.Errorf("expected pop ID %v got empty heap", i)
		} else if previousTime.After(val.Next) {
			t.Errorf("heap pop %v < previous %v expected >", val.Next, previousTime)
		}
		previousTime = val.Next
	}

	val, ok = h.Pop()
	if ok {
		t.Errorf("expected empty, got %+v", val)
	}
}
