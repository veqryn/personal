package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

// http://www.1024cores.net/home/lock-free-algorithms/queues/non-intrusive-mpsc-node-based-queue
// https://github.com/akka/akka/blob/master/akka-actor/src/main/java/akka/dispatch/AbstractNodeQueue.java
// https://github.com/samanbarghi/MPSCQ/blob/master/src/MPSCQueue.h
// https://concurrencyfreaks.blogspot.com/2014/04/multi-producer-single-consumer-queue.html
type NonIntrusiveMpscQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer // volatile
}

func NewQueue() *NonIntrusiveMpscQueue {
	sentinelNode := (unsafe.Pointer)(&node{})
	// TODO: do I need to atomic.Store tail here???
	return &NonIntrusiveMpscQueue{
		head: sentinelNode,
		tail: sentinelNode,
	}
}

type node struct {
	next  unsafe.Pointer // volatile
	value interface{}
}

func (queue *NonIntrusiveMpscQueue) Push(val interface{}) {
	// TODO: is there a way to avoid having a sentinel node or put the new value into the tail instead of the next new node?
	newNode := (unsafe.Pointer)(&node{value: val})
	//fmt.Println("new: ", (*node)(newNode))
	prev := (*node)(atomic.SwapPointer(&queue.tail, newNode)) // serialization-point wrt producers, acquire-release
	atomic.StorePointer(&prev.next, newNode)                  // serialization-point wrt consumer, release
}

func (queue *NonIntrusiveMpscQueue) Pop() interface{} {
	head := (*node)(queue.head)
	next := atomic.LoadPointer(&head.next) // serialization-point wrt producers, acquire
	//fmt.Println("head: ", (unsafe.Pointer)(head), head)
	//fmt.Println("next: ", next, (*node)(next))
	if next != nil {
		//fmt.Println("Popping...")
		queue.head = next
		// Head node always has a nil value
		// TODO: this means we always keep a reference to the last pop'ed value, until the next pop
		return (*node)(next).value

		// TODO: why does this cost so much performance?
		//nextNode := (*node)(next)
		//val := nextNode.value // Head node always has a nil value
		//nextNode.value = nil  // Do not keep a reference to the value (my own modification)
		//return val
	}
	return nil
}

func print(queue *NonIntrusiveMpscQueue) {
	fmt.Println("\nQueue:")
	head := (*node)(atomic.LoadPointer(&queue.head))
	for current := head; current != nil; current = (*node)(atomic.LoadPointer(&current.next)) {
		fmt.Println("current: ", (unsafe.Pointer)(current), current)
		//fmt.Println(current.value)
	}
	last := (*node)(atomic.LoadPointer(&queue.tail))
	fmt.Println("tail: ", (unsafe.Pointer)(last), last)
	fmt.Println("Done\n")
}

func main2() {

	queue := NewQueue()
	print(queue)
	fmt.Println("pop: ", queue.Pop())
	print(queue)
	queue.Push("foo")
	queue.Push("bar")
	print(queue)
	fmt.Println("pop: ", queue.Pop())
	print(queue)
	fmt.Println("pop: ", queue.Pop())
	print(queue)
	fmt.Println("pop: ", queue.Pop())
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

	fmt.Println("pop: ", queue.Pop())
	fmt.Println("pop: ", queue.Pop())
	fmt.Println("pop: ", queue.Pop())
	fmt.Println("pop: ", queue.Pop())
	fmt.Println("pop: ", queue.Pop())
	fmt.Println("pop: ", queue.Pop())
	fmt.Println("pop: ", queue.Pop())
	fmt.Println("pop: ", queue.Pop())
	fmt.Println("pop: ", queue.Pop())
	fmt.Println("pop: ", queue.Pop())
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
