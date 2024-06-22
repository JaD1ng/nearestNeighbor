// @File: points
// @Author: Nanjia Ding
// @Date: 2024/06/17
package main

import (
	"encoding/csv"
	"math/rand"
	"nearestNeighbor/common"
	"os"
	"strconv"
)

func main() {
	numUsers := 1000000
	file, err := os.Create("users.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < numUsers; i++ {
		user := common.UserCoordinate{
			Id:        i,
			Latitude:  rand.Float64()*180 - 90,
			Longitude: rand.Float64()*360 - 180,
		}
		writer.Write([]string{
			strconv.Itoa(user.Id),
			strconv.FormatFloat(user.Latitude, 'f', 6, 64),
			strconv.FormatFloat(user.Longitude, 'f', 6, 64),
		})
	}
}
