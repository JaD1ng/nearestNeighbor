package violent

import (
	"nearestNeighbor/common"
	"nearestNeighbor/tool"
	"sort"
)

func FindClosestPointsViolent(users []*common.UserCoordinate, targetLatitude, targetLongitude float64, n int) []int {
	sort.Slice(users, func(i, j int) bool {
		return tool.HaversineDistance(targetLatitude, targetLongitude, users[i].Latitude, users[i].Longitude) < tool.HaversineDistance(targetLatitude, targetLongitude, users[j].Latitude, users[j].Longitude)
	})

	closestUsers := make([]int, n)
	for i := 0; i < n; i++ {
		closestUsers[i] = users[i].Id
	}
	return closestUsers
}
