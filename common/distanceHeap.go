// @File: distanceHeap
// @Author: Nanjia Ding
// @Date: 2024/06/18
package common

type DistanceHeap []DistanceHeapNode

type DistanceHeapNode struct {
	UserID   int
	Distance float64
}

func (h *DistanceHeap) Len() int           { return len(*h) }
func (h *DistanceHeap) Less(i, j int) bool { return (*h)[i].Distance < (*h)[j].Distance }
func (h *DistanceHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *DistanceHeap) Push(x any) {
	*h = append(*h, x.(DistanceHeapNode))
}

func (h *DistanceHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
