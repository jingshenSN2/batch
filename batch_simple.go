package main

import (
	"fmt"
	"time"
)

type BatchSimple struct {
	queue    []string
	maxDelay time.Duration

	inChan chan string
}

func NewBatchSimple(maxDelay time.Duration) *BatchSimple {
	b := &BatchSimple{
		queue:    make([]string, 0),
		maxDelay: maxDelay,

		inChan: make(chan string),
	}

	go b.run()

	return b
}

func (b *BatchSimple) Post(data string) {
	// Post some data to batching service
	b.inChan <- data
}

func (b *BatchSimple) flush() {
	// Flush data to console
	if len(b.queue) > 0 {
		fmt.Println("Batched result:", b.queue)
		b.queue = make([]string, 0)
	}
	b.queue = make([]string, 0)
}

func (b *BatchSimple) run() {
	ticker := time.NewTicker(b.maxDelay)
	for {
		select {
		case data := <-b.inChan:
			b.queue = append(b.queue, data)
		case <-ticker.C:
			b.flush()
		}
	}
}

func main() {
	b := NewBatchSimple(5 * time.Second)
	// Capture keyboard input
	for {
		var data string
		fmt.Scanln(&data)
		b.Post(data)
	}
}
