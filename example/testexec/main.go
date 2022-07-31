package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	m := sync.Mutex{}

	go func() {
		t := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-t.C:
				m.Lock()
				fmt.Println("tick1")
				m.Unlock()
			}
		}
	}()

	go func() {
		t := time.NewTicker(999 * time.Millisecond)

		for {
			select {
			case <-t.C:
				m.Lock()
				fmt.Println("tick2")
				m.Unlock()
			}
		}
	}()

	w := sync.WaitGroup{}
	w.Add(2)

	for {
		time.Sleep(1 * time.Second)
	}
}
