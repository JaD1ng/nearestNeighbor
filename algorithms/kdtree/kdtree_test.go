// @File: kdtree_test
// @Author: Nanjia Ding
// @Date: 2024/06/20
package kdtree

import (
	"nearestNeighbor/common"
	"testing"
)

func TestFindClosestUsersKDTree(t *testing.T) {
	users := []*common.UserCoordinate{
		{Id: 1, Latitude: 40.7128, Longitude: -74.0060},
		{Id: 2, Latitude: 34.0522, Longitude: -118.2437},
		{Id: 3, Latitude: 51.5074, Longitude: -0.1278},
		{Id: 4, Latitude: 37.7749, Longitude: -122.4194},
		{Id: 5, Latitude: 34.9522, Longitude: -118.2437},
		{Id: 6, Latitude: 51.5074, Longitude: -0.1278},
	}
	targetLatitude := 37.7749
	targetLongitude := -122.4194
	n := 2

	expected := []int{4, 5}
	result := FindClosestUsersKDTree(users, targetLatitude, targetLongitude, n)

	if len(result) != len(expected) {
		t.Fatalf("Expected result length: %d, but got: %d", len(expected), len(result))
	}

	for i, userID := range result {
		if userID != expected[i] {
			t.Errorf("Expected user ID at index %d to be: %d, but got: %d", i, expected[i], userID)
		}
	}
}
