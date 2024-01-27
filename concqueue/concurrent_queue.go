package concqueue

import (
	"fmt"
	"sync"
	"net/http"
)

type ConcurrentQueue struct {
	queue chan *http.Request
	Mutex sync.Mutex
	cond *sync.Cond
}

func (cq *ConcurrentQueue) CheckSize() int {
	return len((*cq).queue)
}

func MakeQueue(size int) *ConcurrentQueue {
	//make a buffered channel
	return &ConcurrentQueue{ queue: make(chan *http.Request, size) }
}

func (cq *ConcurrentQueue) Enqueue (item *http.Request) {
	cq.Mutex.Lock();
	select {
	case cq.queue <- item:
		fmt.Printf("Adding HTTP request to queue at address: %d", &item);
	default:
			fmt.Println("Queue full, dropping item:", item)
	}
	cq.Mutex.Unlock();
}

func (cq *ConcurrentQueue) Dequeue() *http.Request {
	cq.Mutex.Lock();
	select {
		case item := <- cq.queue:
			return item
		default:
			return nil
	}
	cq.Mutex.Unlock();
	return nil
}
