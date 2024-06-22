// @File: external_rtree
// @Author: Nanjia Ding
// @Date: 2024/06/21
package rtree

import (
	"github.com/dhconnelly/rtreego"
	"nearestNeighbor/common"
	"nearestNeighbor/tool"
	"runtime"
	"sort"
	"sync"
)

func FindClosestUsersRTreeV2(users []*common.UserCoordinate, targetLatitude, targetLongitude float64, n int) []int {
	numCPU := runtime.NumCPU()                      // 获取CPU数量
	chunkSize := (len(users) + numCPU - 1) / numCPU // 计算每组的大小

	var wg sync.WaitGroup
	resultsChan := make(chan []*common.UserCoordinate, numCPU)

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * chunkSize
			end := start + chunkSize
			if end > len(users) {
				end = len(users)
			}
			rt := rtreego.NewTree(2, 15, 30)
			for _, user := range users[start:end] {
				rt.Insert(user)
			}
			targetPoint := rtreego.Point{targetLatitude, targetLongitude}
			nearest := rt.NearestNeighbors(5, targetPoint)
			var localResults []*common.UserCoordinate
			for _, result := range nearest {
				user := result.(*common.UserCoordinate)
				localResults = append(localResults, user)
			}
			resultsChan <- localResults
		}(i)
	}

	wg.Wait()
	close(resultsChan)

	allResults := make([]*common.UserCoordinate, 0)
	for res := range resultsChan {
		allResults = append(allResults, res...)
	}

	// 对所有结果进行排序和去重
	sort.Slice(allResults, func(i, j int) bool {
		return tool.HaversineDistance(targetLatitude, targetLongitude, allResults[i].Latitude, allResults[i].Longitude) < tool.HaversineDistance(targetLatitude, targetLongitude, allResults[j].Latitude, allResults[j].Longitude)
	})
	uniqueIds := make(map[int]bool)
	closestUsers := make([]int, 0, n)
	for _, res := range allResults {
		if !uniqueIds[res.Id] {
			uniqueIds[res.Id] = true
			closestUsers = append(closestUsers, res.Id)
			if len(closestUsers) == n {
				break
			}
		}
	}

	return closestUsers
}
