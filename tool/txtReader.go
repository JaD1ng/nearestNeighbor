// @File: txtReader
// @Author: Nanjia Ding
// @Date: 2024/06/19
package tool

import (
	"encoding/csv"
	"io"
	"nearestNeighbor/common"
	"os"
	"runtime"
	"strconv"
	"sync"
)

func LoadUserCoordinatesFromTxt(filePath string) ([]*common.UserCoordinate, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []*common.UserCoordinate
	reader := csv.NewReader(file)
	reader.Comma = ','        // 确保CSV分隔符正确设置
	reader.ReuseRecord = true // 优化内存使用

	var mu sync.Mutex
	var wg sync.WaitGroup
	numCPU := runtime.NumCPU()
	chunks := make(chan []string, numCPU)

	// 启动多个goroutine并行处理数据
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range chunks {
				id, _ := strconv.Atoi(record[0])
				latitude, _ := strconv.ParseFloat(record[1], 64)
				longitude, _ := strconv.ParseFloat(record[2], 64)
				user := &common.UserCoordinate{
					Id:        id,
					Latitude:  latitude,
					Longitude: longitude,
				}
				mu.Lock()
				users = append(users, user)
				mu.Unlock()
			}
		}()
	}

	// 主线程读取文件并分发到各个goroutine
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		chunks <- record
	}
	close(chunks)
	wg.Wait()

	return users, nil
}
