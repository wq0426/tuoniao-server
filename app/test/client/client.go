package main

import (
	"fmt"
)

func splitIntoGroups(arr []int, groupCount int) [][]int {
	if groupCount <= 0 {
		return nil
	}

	// 计算每组的大小
	groupSize := len(arr) / groupCount
	remainder := len(arr) % groupCount

	// 初始化结果数组
	groups := make([][]int, 0, groupCount)
	start := 0

	for i := 0; i < groupCount; i++ {
		end := start + groupSize
		if remainder > 0 {
			end++
			remainder--
		}
		groups = append(groups, arr[start:end])
		start = end
	}

	return groups
}

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	groupCount := 4
	groups := splitIntoGroups(arr, groupCount)
	for i, group := range groups {
		fmt.Printf("Group %d: %v\n", i+1, group)
	}
}
