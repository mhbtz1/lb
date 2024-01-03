package concqueue

import (
	"fmt"
	"sync"
	"net/http"
)

type ConcurrentQueue struct {
	queue chan *http.Request
	Mutex sync.Mutex
}

func (cq *ConcurrentQueue) CheckSize() int {
	return len((*cq).queue)
}

func MakeQueue(size int) *ConcurrentQueue {
	//make a buffered channel
	return &ConcurrentQueue{ queue: make(chan *http.Request, size) }
}

func (cq *ConcurrentQueue) Enqueue (item *http.Request) {
	cq.Mutex.Lock()
	defer cq.Mutex.Unlock()
	select {
	case cq.queue <- item:
		default:
			fmt.Println("Queue full, dropping item:", item)
	}
}

func (cq *ConcurrentQueue) Dequeue() *http.Request {
	cq.Mutex.Lock()
	defer cq.Mutex.Unlock()
	select {
		case item := <- cq.queue:
			return item
		default:
			return nil
	}
}
