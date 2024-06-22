// @File: kdtree
// @Author: Nanjia Ding
// @Date: 2024/06/18
package kdtree

import (
	"container/heap"
	"math"
	"nearestNeighbor/common"
	"nearestNeighbor/tool"
	"runtime"
	"sort"
	"sync"
)

func FindClosestUsersKDTree(users []*common.UserCoordinate, targetLatitude, targetLongitude float64, n int) []int {
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
	root := newKDTree(users, 0) // 使用KD树
	target := &common.UserCoordinate{Latitude: targetLatitude, Longitude: targetLongitude}
	h := &common.DistanceHeap{}
	heap.Init(h)

	findNearest(root, target, h, n, 0) // 使用KD树查找最近的用户

	result := make([]*common.DistanceHeapNode, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		node := heap.Pop(h).(common.DistanceHeapNode)
		result[i] = &node
	}
	return result
}

type kdNode struct {
	point       *common.UserCoordinate
	left, right *kdNode
}

func newKDTree(points []*common.UserCoordinate, depth int) *kdNode {
	n := len(points)
	if n == 0 {
		return nil
	}

	sort.Slice(points, func(i, j int) bool {
		if depth%2 == 0 {
			return points[i].Latitude < points[j].Latitude
		}
		return points[i].Longitude < points[j].Longitude
	})

	median := n / 2

	return &kdNode{
		point: points[median],
		left:  newKDTree(points[:median], depth+1),
		right: newKDTree(points[median+1:], depth+1),
	}
}

func findNearest(node *kdNode, target *common.UserCoordinate, h *common.DistanceHeap, n int, depth int) {
	if node == nil {
		return
	}

	dist := tool.HaversineDistance(target.Latitude, target.Longitude, node.point.Latitude, node.point.Longitude)
	if h.Len() < n {
		heap.Push(h, common.DistanceHeapNode{UserID: node.point.Id, Distance: dist})
	} else if dist < (*h)[0].Distance {
		heap.Pop(h)
		heap.Push(h, common.DistanceHeapNode{UserID: node.point.Id, Distance: dist})
	}

	dim := depth % 2
	var first, second *kdNode
	if (dim == 0 && target.Latitude < node.point.Latitude) || (dim == 1 && target.Longitude < node.point.Longitude) {
		first, second = node.left, node.right
	} else {
		first, second = node.right, node.left
	}

	findNearest(first, target, h, n, depth+1)
	if h.Len() < n || (dim == 0 && math.Abs(target.Latitude-node.point.Latitude) < (*h)[0].Distance) || (dim == 1 && math.Abs(target.Longitude-node.point.Longitude) < (*h)[0].Distance) {
		findNearest(second, target, h, n, depth+1)
	}
}
