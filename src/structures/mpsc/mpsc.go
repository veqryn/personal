package mpsc

import (
	"fmt"
	"math"
	"runtime"
	"sync/atomic"
	"time"
)

// https://concurrencyfreaks.blogspot.com/2017/03/multilist-mpsc-wait-free-queue.html
// https://github.com/pramalhe/ConcurrencyFreaks/blob/master/papers/multilist-2017.pdf
type MPSC struct {
	partitions int
	counter    uint64
	heads      []*Node
	tails      []*Node
}

type Node struct {
	val   interface{}
	next  *Node
	count uint64
}

func New(partitions int) *MPSC {
	mpsc := MPSC{
		partitions: partitions,
		counter:    0,
		heads:      make([]*Node, partitions),
		tails:      make([]*Node, partitions),
	}
	for i := 0; i < partitions; i++ {
		mpsc.heads[i] = &Node{count: math.MaxUint64}
		mpsc.tails[i] = mpsc.heads[i]
	}
	return &mpsc
}

func (mpsc *MPSC) Push(partition int, val interface{}) {
	current := atomic.LoadUint64(&mpsc.counter)
	ltail := mpsc.tails[partition]
	mpsc.tails[partition] = &Node{count: math.MaxUint64}
	ltail.val = val
	ltail.next = mpsc.tails[partition]
	atomic.StoreUint64(&ltail.count, current) // Barrier
	atomic.AddUint64(&mpsc.counter, 1)
}

func (mpsc *MPSC) Pop() (interface{}, bool) {
	// TODO: no need for 2 iterations if the queue is really full
	//current := atomic.LoadUint64(&mpsc.counter)
	prevIdx := -2
	for {
		var minCount uint64 = math.MaxUint64
		minIdx := -1
		for i := 0; i < mpsc.partitions; i++ {
			lcount := atomic.LoadUint64(&mpsc.heads[i].count)
			if lcount < minCount {
				minCount = lcount
				minIdx = i
			}
		}
		//if minIdx >= 0 && minCount+1000 < current {
		//	lhead := mpsc.heads[minIdx]
		//	mpsc.heads[minIdx] = lhead.next
		//	return lhead.val, true
		//}
		//if minIdx >= 0 {
		//	lhead := mpsc.heads[minIdx]
		//	mpsc.heads[minIdx] = lhead.next
		//	return lhead.val, true
		//}
		if minIdx == -1 && prevIdx == minIdx {
			return nil, false
		}
		// TODO: fix bias towards lower indexed partitions
		if prevIdx == minIdx {
			lhead := mpsc.heads[minIdx]
			mpsc.heads[minIdx] = lhead.next
			return lhead.val, true
		}
		prevIdx = minIdx
	}
}

func main() {
	producerCount := runtime.NumCPU() - 1
	mpsc := New(producerCount)

	for i := 0; i < producerCount; i++ {
		go func(partition int) {
			time.Sleep(time.Duration(20-partition) * time.Millisecond)
			for j := 0; j < 10000000; j++ {
				mpsc.Push(partition, partition)
			}
			fmt.Println(partition, "done")
		}(i)
	}

	for {
		if _, ok := mpsc.Pop(); ok {
		} else {
			fmt.Println("sleeping")
			time.Sleep(10 * time.Millisecond)
		}
	}
}
