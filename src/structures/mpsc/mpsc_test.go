package mpsc

import (
	"sync"
	"sync/atomic"
	"testing"
)

// https://github.com/jakewins/4fq
func BenchmarkMpscQueue(b *testing.B) {
	producerCount := 32
	runningProducers := int64(producerCount)

	mpsc := New(producerCount)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < producerCount; i++ {
		go func(partition int) {
			for n := 0; n < b.N; n++ {
				mpsc.Push(partition, n)
			}
			atomic.AddInt64(&runningProducers, -1)
		}(i)
	}

	for {
		if _, ok := mpsc.Pop(); !ok && atomic.LoadInt64(&runningProducers) <= 0 {
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
