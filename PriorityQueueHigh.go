// this is roughly started from the code at
// https://golang.org/pkg/container/heap/

package main

import (
	"container/heap"
)

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueueHigh []*PqItem

//type PqItem struct {
//	containedItem *SequentialInterface
//	priority      int
//	index         int
//}

func (pq PriorityQueueHigh) Len() int {
	return len(pq)
}

func (pq PriorityQueueHigh) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueueHigh) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueueHigh) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueueHigh) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueueHigh) update(item *PqItem, value *SequentialInterface, priority int) {
	item.containedItem = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

//Added priority to arguments to ease switch to greedy alg
func (pq *PriorityQueueHigh) PushSequentialInterface(newItem *SequentialInterface,priority int) {
	heldItem := PqItem{
		containedItem: newItem,
		priority:      priority,
		index:         len(*pq),
	}
	heap.Push(pq, &heldItem)
}

func (pq *PriorityQueueHigh) PopSequentialInterface() *SequentialInterface {
	item := heap.Pop(pq).(*PqItem)
	return item.containedItem
}
