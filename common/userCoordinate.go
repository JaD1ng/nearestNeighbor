// @File: userCoordinate
// @Author: Nanjia Ding
// @Date: 2024/06/18
package common

import (
	"github.com/dhconnelly/rtreego"
)

type UserCoordinate struct {
	Id        int
	Latitude  float64
	Longitude float64
}

func (u *UserCoordinate) Bounds() rtreego.Rect {
	point := rtreego.Point{u.Latitude, u.Longitude}
	// 0.01表示用户坐标的误差范围
	return point.ToRect(0.01)
}
