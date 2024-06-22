// @File: parallel
// @Author: Nanjia Ding
// @Date: 2024/06/18
package parallel

import (
	"container/heap"
	"nearestNeighbor/common"
	"nearestNeighbor/tool"
	"runtime"
	"sync"
)

func FindClosestUsersParallel(users []*common.UserCoordinate, targetLatitude, targetLongitude float64, n int) []int {
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
	distanceHeap := &common.DistanceHeap{}
	heap.Init(distanceHeap)

	for _, user := range users {
		distance := tool.HaversineDistance(targetLatitude, targetLongitude, user.Latitude, user.Longitude)
		if distanceHeap.Len() < n {
			heap.Push(distanceHeap, common.DistanceHeapNode{
				UserID:   user.Id,
				Distance: distance,
			})
		} else if distance < (*distanceHeap)[0].Distance-1e-9 {
			heap.Pop(distanceHeap)
			heap.Push(distanceHeap, common.DistanceHeapNode{
				UserID:   user.Id,
				Distance: distance,
			})
		}
	}

	result := make([]*common.DistanceHeapNode, distanceHeap.Len())
	for i := distanceHeap.Len() - 1; i >= 0; i-- {
		node := heap.Pop(distanceHeap).(common.DistanceHeapNode)
		result[i] = &node
	}
	return result
}
