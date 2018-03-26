package main

import (
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"
)

type CowSlice struct {
	ptr unsafe.Pointer
}

func New(newSlice []interface{}) *CowSlice {
	cow := CowSlice{}
	cow.Set(newSlice)
	return &cow
}

func (cow *CowSlice) Get() []interface{} {
	return *(*[]interface{})(atomic.LoadPointer(&cow.ptr))
}

func (cow *CowSlice) Set(newSlice []interface{}) {
	atomic.StorePointer(&cow.ptr, (unsafe.Pointer)(&newSlice))
}
func (cow *CowSlice) Append(i ...interface{}) {
	cow.Set(append(cow.Get(), i...))
}

func main() {

	init := make([]interface{}, 0, 2)
	cow := New(init)

	go func() {
		time.Sleep(30 * time.Millisecond)
		current := cow.Get()
		fmt.Println("gor2", current)

		time.Sleep(100 * time.Millisecond)
		for _, v := range cow.Get() {
			fmt.Println("gor2", v)
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Println("gor2", current)

		time.Sleep(4 * time.Second)
		fmt.Println("gor2", current)

		current = cow.Get()
		fmt.Println("gor2", current)
	}()

	time.Sleep(100 * time.Millisecond)
	init = append(init, "blah")
	fmt.Println("gor1", init)

	cow.Append("hello")
	cow.Append("world")
	fmt.Println("gor1", cow.Get())
	time.Sleep(30 * time.Millisecond)
	cow.Append("daisy")
	fmt.Println("gor1", cow.Get())
	time.Sleep(100 * time.Millisecond)
	cow.Append("foo", "bar", "baz")
	fmt.Println("gor1", init)
	fmt.Println("gor1", cow.Get())

	time.Sleep(100 * time.Millisecond)
	cow.Set([]interface{}{map[string]string{
		"foo": "bar",
		"bar": "baz",
		"baz": "qux",
		"qux": "naar",
	}})

	time.Sleep(5 * time.Second)
}
