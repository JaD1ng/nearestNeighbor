// @File: lsh
// @Author: Nanjia Ding
// @Date: 2024/06/18
package lsh

import (
	"container/heap"
	"nearestNeighbor/common"
	"nearestNeighbor/tool"
	"runtime"
	"sync"
)

func FindClosestUsersLSH(users []*common.UserCoordinate, targetLatitude, targetLongitude float64, n int) []int {
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
		if i == numGroups-1 {
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
		} else if distance < (*distanceHeap)[0].Distance {
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

type lsh struct {
	data      []*common.UserCoordinate
	hashTable map[int][]*common.UserCoordinate
}

// 简单的哈希函数，将坐标映射到一个整数
func hash(coord *common.UserCoordinate) int {
	const p = 16777619
	hash := 2166136261

	// 简单的哈希函数，使用经纬度的整数部分
	hash = (hash ^ int(coord.Latitude)) * p
	hash = (hash ^ int(coord.Longitude)) * p

	// 确保hash值为非负
	hash += hash << 13
	hash ^= hash >> 7
	hash += hash << 3
	hash ^= hash >> 17
	hash += hash << 5

	if hash < 0 {
		hash = -hash
	}
	return hash
}

// 构建哈希表
func (lsh *lsh) newHashTable() {
	lsh.hashTable = make(map[int][]*common.UserCoordinate)
	for _, point := range lsh.data {
		h := hash(point)
		lsh.hashTable[h] = append(lsh.hashTable[h], point)
	}
}

func (lsh *lsh) nearest(coord *common.UserCoordinate, n int) []*common.UserCoordinate {
	h := hash(coord)
	candidates := lsh.hashTable[h]

	minHeap := &common.DistanceHeap{}
	heap.Init(minHeap)

	for _, candidate := range candidates {
		distance := tool.HaversineDistance(coord.Latitude, coord.Longitude, candidate.Latitude, candidate.Longitude)
		heap.Push(minHeap, common.DistanceHeapNode{UserID: candidate.Id, Distance: distance})
		if minHeap.Len() > n {
			heap.Pop(minHeap)
		}
	}

	result := make([]*common.UserCoordinate, minHeap.Len())
	for i := len(result) - 1; i >= 0; i-- {
		node := heap.Pop(minHeap).(common.DistanceHeapNode)
		result[i] = &common.UserCoordinate{Id: node.UserID, Latitude: coord.Latitude, Longitude: coord.Longitude}
	}

	return result
}
