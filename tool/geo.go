// @File: geo
// @Author: Nanjia Ding
// @Date: 2024/06/18
package tool

import "math"

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	R := 6371.0 // 地球半径，单位为千米
	phi1 := ToRadians(lat1)
	phi2 := ToRadians(lat2)
	deltaPhi := ToRadians(lat2 - lat1)
	deltaLambda := ToRadians(lon2 - lon1)

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func ToRadians(degree float64) float64 {
	return degree * math.Pi / 180
}
