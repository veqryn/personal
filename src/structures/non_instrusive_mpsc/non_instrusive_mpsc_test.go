package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// https://github.com/jakewins/4fq
func BenchmarkQueue(b *testing.B) {
	producerCount := 32
	runningProducers := int64(producerCount)

	queue := NewQueue()

	b.ReportAllocs()
	b.ResetTimer()
	//fmt.Println("---- starting ----", b.N)
	for i := 0; i < producerCount; i++ {
		go func() {
			for n := 0; n < b.N; n++ {
				//fmt.Println("---- pushing ----", n)
				queue.Push(n)
			}
			atomic.AddInt64(&runningProducers, -1)
			//fmt.Println("---- producer done ----")
		}()
	}

	for {
		if val := queue.Pop(); val == nil {
			if atomic.LoadInt64(&runningProducers) <= 0 {
				break
			} else {
				time.Sleep(1 * time.Nanosecond)
			}
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
	//fmt.Println("---- starting ----")
	for i := 0; i < producerCount; i++ {
		go func() {
			for n := 0; n < b.N; n++ {
				ch <- n
			}
			atomic.AddInt64(&runningProducers, -1)
			wg.Done()
			//fmt.Println("---- producer done ----")
		}()
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
