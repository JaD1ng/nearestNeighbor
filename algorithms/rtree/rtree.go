// @File: rtree
// @Author: Nanjia Ding
// @Date: 2024/06/17
package rtree

import (
	"container/heap"
	"math"
	"nearestNeighbor/common"
	"nearestNeighbor/tool"
	"runtime"
	"sync"
)

func FindClosestUsersRTree(users []*common.UserCoordinate, targetLatitude, targetLongitude float64, n int) []int {
	numGroups := runtime.NumCPU()
	if numGroups > len(users) {
		numGroups = len(users)
	}
	groupSize := len(users) / numGroups
	var wg sync.WaitGroup
	groupRes := sync.Map{}

	for i := 0; i < numGroups; i++ {
		start := i * groupSize
		end := start + groupSize
		if end > len(users) {
			end = len(users)
		}
		wg.Add(1)
		go func(start, end, i int) {
			defer wg.Done()
			groupRes.Store(i, processGroup(users[start:end], targetLatitude, targetLongitude, n))
		}(start, end, i)
	}
	wg.Wait()

	res := make(common.DistanceHeap, 0, numGroups*n)
	groupRes.Range(func(_, value interface{}) bool {
		group := value.([]*common.DistanceHeapNode)
		for _, node := range group {
			res = append(res, *node)
		}
		return true
	})

	heap.Init(&res)
	closestUsers := make([]int, n)
	for i := 0; i < n; i++ {
		closestUsers[i] = heap.Pop(&res).(common.DistanceHeapNode).UserID
	}

	return closestUsers
}

func processGroup(users []*common.UserCoordinate, targetLatitude, targetLongitude float64, n int) []*common.DistanceHeapNode {
	rt := newRTree()
	for _, user := range users {
		rt.insert(user)
	}
	target := &common.UserCoordinate{Latitude: targetLatitude, Longitude: targetLongitude}
	nearestUsers := rt.nearestNeighbors(target, n)

	result := make([]*common.DistanceHeapNode, len(nearestUsers))
	for i, user := range nearestUsers {
		distance := tool.HaversineDistance(targetLatitude, targetLongitude, user.Latitude, user.Longitude)
		result[i] = &common.DistanceHeapNode{
			UserID:   user.Id,
			Distance: distance,
		}
	}
	return result
}

type rect struct {
	minX, minY, maxX, maxY float64
}

func (r *rect) contains(point *common.UserCoordinate) bool {
	return point.Latitude >= r.minX && point.Latitude <= r.maxX &&
		point.Longitude >= r.minY && point.Longitude <= r.maxY
}

type rTreeNode struct {
	bounds   *rect
	entries  []*common.UserCoordinate
	children []*rTreeNode
	isLeaf   bool
}

type rTree struct {
	root *rTreeNode
}

func newRTreeNode(bounds *rect, isLeaf bool) *rTreeNode {
	return &rTreeNode{
		bounds:   bounds,
		isLeaf:   isLeaf,
		entries:  []*common.UserCoordinate{},
		children: []*rTreeNode{},
	}
}

func newRTree() *rTree {
	return &rTree{
		root: newRTreeNode(&rect{minX: -90, minY: -180, maxX: 90, maxY: 180}, true),
	}
}

func (tree *rTree) insert(entry *common.UserCoordinate) {
	node := tree.root
	for !node.isLeaf {
		var chosenChild *rTreeNode
		closestChild := node.children[0] // 默认选择第一个子节点
		minDist := math.MaxFloat64

		for _, child := range node.children {
			if child.bounds.contains(entry) {
				chosenChild = child
				break
			}
			// 计算距离并选择最近的子节点
			dist := distanceToRect(child.bounds, entry)
			if dist < minDist {
				minDist = dist
				closestChild = child
			}
		}

		if chosenChild == nil {
			chosenChild = closestChild
		}
		node = chosenChild
	}
	node.entries = append(node.entries, entry)
}

func distanceToRect(r *rect, point *common.UserCoordinate) float64 {
	dx := math.Max(math.Max(r.minX-point.Latitude, 0), point.Latitude-r.maxX)
	dy := math.Max(math.Max(r.minY-point.Longitude, 0), point.Longitude-r.maxY)
	return math.Sqrt(dx*dx + dy*dy)
}

func (tree *rTree) nearestNeighbors(target *common.UserCoordinate, n int) []*common.UserCoordinate {
	h := &common.DistanceHeap{}
	heap.Init(h)
	tree.search(tree.root, target, h, n)
	result := make([]*common.UserCoordinate, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		node := heap.Pop(h).(common.DistanceHeapNode)
		result[i] = &common.UserCoordinate{Id: node.UserID}
	}
	return result
}

func (tree *rTree) search(node *rTreeNode, target *common.UserCoordinate, h *common.DistanceHeap, n int) {
	if node.isLeaf {
		for _, entry := range node.entries {
			dist := tool.HaversineDistance(target.Latitude, target.Longitude, entry.Latitude, entry.Longitude)
			heap.Push(h, common.DistanceHeapNode{UserID: entry.Id, Distance: dist})
			if h.Len() > n {
				heap.Pop(h)
			}
		}
	} else {
		for _, child := range node.children {
			if child.bounds.contains(target) {
				tree.search(child, target, h, n)
			}
		}
	}
}
