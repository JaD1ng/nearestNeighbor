// @File: redisReader
// @Author: Nanjia Ding
// @Date: 2024/06/19
package tool

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"nearestNeighbor/common"
	"nearestNeighbor/config"
	"sync"
)

func LoadUserCoordinatesFromRedis() ([]*common.UserCoordinate, error) {
	conf, err := config.LoadConfig("config.json")
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	ctx := context.Background()

	// 创建一个channel用于接收数据
	dataCh := make(chan string, 1000)
	var wg sync.WaitGroup
	wg.Add(1)

	// 创建一个goroutine用于从Redis读取数据并发送到channel
	go func() {
		defer wg.Done()
		var cursor uint64
		for {
			var keys []string
			var err error
			keys, cursor, err = rdb.Scan(ctx, cursor, "user_coordinates:*", 1000).Result()
			if err != nil {
				close(dataCh)
				return
			}
			for _, key := range keys {
				val, err := rdb.Get(ctx, key).Result()
				if err != nil {
					continue
				}
				dataCh <- val
			}
			if cursor == 0 {
				break
			}
		}
		close(dataCh)
	}()

	// 创建一个goroutine用于从channel接收数据并处理
	users := make([]*common.UserCoordinate, 0)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range dataCh {
			var user common.UserCoordinate
			err := json.Unmarshal([]byte(data), &user)
			if err != nil {
				continue
			}
			users = append(users, &user)
		}
	}()

	wg.Wait()

	return users, nil
}
