// this is roughly started from the code at
// https://golang.org/pkg/container/heap/

package main

import (
	"container/heap"
)

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*PqItem

type PqItem struct {
	containedItem *SequentialInterface
	priority      int
	index         int
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *PqItem, value *SequentialInterface, priority int) {
	item.containedItem = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

//Added priority to arguments to ease switch to greedy alg
func (pq *PriorityQueue) PushSequentialInterface(newItem *SequentialInterface,priority int) {
	heldItem := PqItem{
		containedItem: newItem,
		priority:      priority,
		index:         len(*pq),
	}
	heap.Push(pq, &heldItem)
}

func (pq *PriorityQueue) PopSequentialInterface() *SequentialInterface {
	item := heap.Pop(pq).(*PqItem)
	return item.containedItem
}
