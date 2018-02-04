package main

import (
	"time"
	"fmt"
	"sync"
)

// JustifiedLock is implementation of lock with beautiful bisqwit's video
// https://www.youtube.com/watch?v=OrQ9swvm_VA
type JustifiedLock struct {
	lock   sync.Mutex
	reason string
}

// Lock locks with why
func (l *JustifiedLock) Lock(why string) {
	reportTick := time.NewTicker(time.Second)
	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <- reportTick.C:
				fmt.Printf("paused processing of %s because %s\n", why, l.reason)
			case <- done:
				return
			}
		}
	}()

	l.lock.Lock()
	l.reason = why
	done <- struct{}{}
}

// Unlocks performs unlock
func (l *JustifiedLock) Unlock() {
	l.lock.Unlock()
}


func main() {

	g := sync.WaitGroup{}
	jl := JustifiedLock{}


	fmt.Println("1")
	jl.Lock("i wanted")
	go func() {
		g.Add(1)
		time.Sleep(4*time.Second)
		jl.Unlock()
		g.Done()
	}()

	go func() {
		g.Add(1)
		time.Sleep(2*time.Second)
		jl.Lock("druhy hovno")
		fmt.Println("after druhy hovno")
		jl.Unlock()
		g.Done()
	}()


	fmt.Println("2")
	jl.Lock("hovno")
	jl.Unlock()
	fmt.Println("after hovno")

	g.Wait()


}
