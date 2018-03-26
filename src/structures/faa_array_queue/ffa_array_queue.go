package ffa_array_queue

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

const buffer = 256

var takenPtr = &struct{}{}

// https://github.com/pramalhe/ConcurrencyFreaks/blob/master/Java/com/concurrencyfreaks/queues/array/FAAArrayQueue.java
type FaaArrayQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewQueue() *FaaArrayQueue {
	sentinelNode := (unsafe.Pointer)(&node{
		enqueueIdx: 1,
		items:      [buffer]unsafe.Pointer{},
	})
	return &FaaArrayQueue{
		head: sentinelNode,
		tail: sentinelNode,
	}
}

type node struct {
	dequeueIdx uint64
	enqueueIdx uint64
	items      [buffer]unsafe.Pointer
	next       unsafe.Pointer
}

func newNode(val interface{}) *node {
	return &node{
		enqueueIdx: 1,
		items:      [buffer]unsafe.Pointer{(unsafe.Pointer)(&val)},
	}
}

func print(queue *FaaArrayQueue) {
	fmt.Println("\nQueue:")
	head := (*node)(atomic.LoadPointer(&queue.head))
	for current := head; current != nil; current = (*node)(atomic.LoadPointer(&current.next)) {
		fmt.Println(current)
		for idx := 0; idx < len(current.items); idx++ {
			valPtr := atomic.LoadPointer(&current.items[idx])
			val := (*interface{})(valPtr)
			if val != nil && valPtr != (unsafe.Pointer)(takenPtr) {
				fmt.Println(*val)
			}
		}
		fmt.Println()
	}
	fmt.Println((*node)(atomic.LoadPointer(&queue.tail)))
	fmt.Println("Done\n")
}

func (queue *FaaArrayQueue) Push(val interface{}) {
	for {
		ltail := (*node)(atomic.LoadPointer(&queue.tail))
		idx := atomic.AddUint64(&ltail.enqueueIdx, 1) - 1
		if idx > buffer-1 {
			// Node is full
			if ltail != (*node)(atomic.LoadPointer(&queue.tail)) {
				continue
			}
			lnext := (*node)(atomic.LoadPointer(&ltail.next))
			if lnext == nil {
				newNode := newNode(val)
				if atomic.CompareAndSwapPointer(&ltail.next, nil, (unsafe.Pointer)(newNode)) {
					atomic.CompareAndSwapPointer(&queue.tail, (unsafe.Pointer)(ltail), (unsafe.Pointer)(newNode))
					return
				}
			} else {
				atomic.CompareAndSwapPointer(&queue.tail, (unsafe.Pointer)(ltail), (unsafe.Pointer)(lnext))
			}
			continue
		}
		if atomic.CompareAndSwapPointer(&ltail.items[idx], nil, (unsafe.Pointer)(&val)) {
			return
		}
	}
}

func (queue *FaaArrayQueue) Pop() interface{} {
	for {
		lhead := (*node)(atomic.LoadPointer(&queue.head))
		if atomic.LoadUint64(&lhead.dequeueIdx) >= atomic.LoadUint64(&lhead.enqueueIdx) && atomic.LoadPointer(&lhead.next) == nil {
			return nil
		}
		idx := atomic.AddUint64(&lhead.dequeueIdx, 1) - 1
		if idx > buffer-1 {
			// Node already drained
			if ((*node)(atomic.LoadPointer(&lhead.next))) == nil {
				return nil // No more nodes
			}
			atomic.CompareAndSwapPointer(&queue.head, (unsafe.Pointer)(lhead), (unsafe.Pointer)(lhead.next))
			continue
		}
		if val := (*interface{})(atomic.SwapPointer(&lhead.items[idx], (unsafe.Pointer)(takenPtr))); val != nil {
			return *val
		}
	}
}

func main2() {

	queue := NewQueue()
	print(queue)
	queue.Push("hello")
	queue.Push("w")
	queue.Push("o")
	queue.Push("r")
	queue.Push("l")
	queue.Push("d")
	queue.Push("omg")
	queue.Push("foo")
	queue.Push("bar")
	queue.Push("buz")

	fmt.Println(queue.Pop())
	fmt.Print(queue.Pop())
	fmt.Print(queue.Pop())
	fmt.Print(queue.Pop())
	fmt.Print(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
}

func main() {
	producerCount := runtime.NumCPU() - 1
	queue := NewQueue()

	for i := 0; i < producerCount; i++ {
		go func(partition int) {
			time.Sleep(time.Duration(20-partition) * time.Millisecond)
			for j := 0; j < 10000000; j++ {
				queue.Push(partition)
			}
			fmt.Println(partition, "done")
		}(i)
	}

	for {
		if val := queue.Pop(); val != nil {
			//fmt.Println(val)
		} else {
			fmt.Println("sleeping")
			time.Sleep(10 * time.Millisecond)
		}
	}
}
