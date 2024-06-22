// @File: main
// @Author: Nanjia Ding
// @Date: 2024/06/18
package main

import (
	"fmt"
	"nearestNeighbor/algorithms/kdtree"
	"nearestNeighbor/algorithms/lsh"
	"nearestNeighbor/algorithms/parallel"
	"nearestNeighbor/algorithms/rtree"
	"nearestNeighbor/tool"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	filePath := "./tests/users.txt"
	queryNum := 5
	targetLatitude := 37.7749
	targetLongitude := -122.4194

	totalTime := time.Now()
	users, err := tool.LoadUserCoordinatesFromTxt(filePath)
	if err != nil {
		panic(err)
	}

	kdtreeTime := time.Now()
	arr1 := kdtree.FindClosestUsersKDTree(users, targetLatitude, targetLongitude, queryNum)
	fmt.Println("KD树耗时：", time.Since(kdtreeTime))
	fmt.Println(arr1)

	rtreeTime := time.Now()
	arr2 := rtree.FindClosestUsersRTreeV2(users, targetLatitude, targetLongitude, queryNum)
	fmt.Println("R树耗时：", time.Since(rtreeTime))
	fmt.Println(arr2)

	lshTime := time.Now()
	arr3 := lsh.FindClosestUsersLSH(users, targetLatitude, targetLongitude, queryNum)
	fmt.Println("LSH耗时：", time.Since(lshTime))
	fmt.Println(arr3)

	parallelTime := time.Now()
	arr4 := parallel.FindClosestUsersParallel(users, targetLatitude, targetLongitude, queryNum)
	fmt.Println("分组耗时：", time.Since(parallelTime))
	fmt.Println(arr4)

	// violentTime := time.Now()
	// arr5 := violent.FindClosestPointsViolent(users, targetLatitude, targetLongitude, queryNum)
	// fmt.Println("暴力耗时：", time.Since(violentTime))
	// fmt.Println(arr5)

	fmt.Println("总耗时：", time.Since(totalTime))
}
