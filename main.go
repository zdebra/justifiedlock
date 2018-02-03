package main

import (
	"sync/atomic"
	"time"
	"fmt"
	"sync"
)

// JustifiedLock is implementation of lock with beautiful bisqwit's video
// https://www.youtube.com/watch?v=OrQ9swvm_VA
type JustifiedLock struct {
	locked atomic.Value
	reason string
}

// Lock locks with why
func (lock *JustifiedLock) Lock(why string) {
	reportTick := time.NewTicker(time.Second)
	unlockTick := time.NewTicker(250*time.Microsecond)
	for {
		select {
			case <- reportTick.C:
				fmt.Printf("paused processing of %s because %s\n", why, lock.reason)
			case <- unlockTick.C:
				val := lock.locked.Load()
				if locked, ok := val.(bool); !ok || !locked {
					lock.locked.Store(true)
					lock.reason = why
					return
				}
		}
	}
}

// Unlocks performs unlock
func (lock *JustifiedLock) Unlock() {
	lock.locked.Store(false)
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
