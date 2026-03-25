package quickselect

import (
	"math/rand"
	"time"
)

func quickSelect(nums []int, k int) int {
	if len(nums) == 1 {
		return nums[0]
	}

	rand.Seed(time.Now().UnixNano())
	pivot := nums[rand.Intn(len(nums))]

	var left, right, pivots []int
	for _, n := range nums {
		switch {
		case n > pivot:
			left = append(left, n)
		case n < pivot:
			right = append(right, n)
		default:
			pivots = append(pivots, n)
		}
	}

	if k < len(left) {
		return quickSelect(left, k)
	} else if k < len(left)+len(pivots) {
		return pivot
	} else {
		return quickSelect(right, k-len(left)-len(pivots))
	}
}

func TopN[T any](items []T, k int, getter func(a T) int) []T {
	if k == 0 {
		return nil
	}
	if len(items) == 0 {
		return nil
	}

	if k > len(items) {
		k = len(items)
	}

	nums := make([]int, 0, len(items))
	for _, numss := range items {
		nums = append(nums, getter(numss))
	}

	threshold := quickSelect(nums, k-1)

	top := make([]T, 0, k)
	for _, item := range items {
		if getter(item) >= threshold {
			top = append(top, item)
		}
	}
	return top
}
