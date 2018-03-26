package ffa_array_queue

import (
	"sync"
	"sync/atomic"
	"testing"
)

// https://github.com/jakewins/4fq
func BenchmarkQueue(b *testing.B) {
	producerCount := 32
	runningProducers := int64(producerCount)

	queue := NewQueue()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < producerCount; i++ {
		go func(partition int) {
			for n := 0; n < b.N; n++ {
				queue.Push(n)
			}
			atomic.AddInt64(&runningProducers, -1)
		}(i)
	}

	for {
		if ok := queue.Pop(); ok == nil && atomic.LoadInt64(&runningProducers) <= 0 {
			break
		}
	}
}

// For reference, same use case as above but using regular channels
func BenchmarkChannel(b *testing.B) {
	producerCount := 32
	runningProducers := int64(producerCount)

	wg := &sync.WaitGroup{}
	wg.Add(producerCount)

	ch := make(chan int, 1024)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < producerCount; i++ {
		go func(partition int) {
			for n := 0; n < b.N; n++ {
				ch <- n
			}
			atomic.AddInt64(&runningProducers, -1)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for {
		if _, ok := <-ch; !ok && atomic.LoadInt64(&runningProducers) <= 0 {
			break
		}
	}
}
