package quickselect

import (
	"sort"
	"testing"
)

// --- quickSelect ---
// quickSelect returns the k-th largest element (0 = largest, 1 = second largest, …).

func Test_quickSelect_kthLargest(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{"single element", []int{42}, 0, 42},
		{"largest of two", []int{3, 1}, 0, 3},
		{"smallest of two", []int{3, 1}, 1, 1},
		{"largest of five", []int{5, 3, 1, 4, 2}, 0, 5},
		{"second largest", []int{5, 3, 1, 4, 2}, 1, 4},
		{"median", []int{1, 2, 3, 4, 5}, 2, 3},
		{"second smallest", []int{5, 3, 1, 4, 2}, 3, 2},
		{"smallest of five", []int{5, 3, 1, 4, 2}, 4, 1},
		{"negative values", []int{-3, -1, -4, -1, -5}, 0, -1},
		{"mixed neg/pos", []int{-1, 0, 1}, 1, 0},
		{"all equal", []int{5, 5, 5}, 1, 5},
		{"duplicates", []int{1, 3, 3, 5}, 1, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nums := make([]int, len(tt.nums))
			copy(nums, tt.nums)
			got := quickSelect(nums, tt.k)
			if got != tt.want {
				t.Errorf("quickSelect(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}

// --- TopN ---

func intGetter(a int) int { return a }

func TestTopN_emptyOnZeroK(t *testing.T) {
	result := TopN([]int{1, 2, 3}, 0, intGetter)
	if result == nil || len(result) != 0 {
		t.Errorf("expected empty slice for k=0, got %v", result)
	}
}

func TestTopN_emptyOnEmptyItems(t *testing.T) {
	result := TopN([]int{}, 3, intGetter)
	if result == nil || len(result) != 0 {
		t.Errorf("expected empty slice for empty items, got %v", result)
	}
}

func TestTopN_kGreaterThanLen_returnsAll(t *testing.T) {
	result := TopN([]int{1, 2, 3}, 10, intGetter)
	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d: %v", len(result), result)
	}
}

func TestTopN_kEqualsLen_returnsAll(t *testing.T) {
	result := TopN([]int{3, 1, 2}, 3, intGetter)
	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d: %v", len(result), result)
	}
}

func TestTopN_singleItem(t *testing.T) {
	result := TopN([]int{99}, 1, intGetter)
	if len(result) != 1 || result[0] != 99 {
		t.Errorf("expected [99], got %v", result)
	}
}

func TestTopN_top1_returnsMax(t *testing.T) {
	result := TopN([]int{4, 7, 1, 9, 3}, 1, intGetter)
	if len(result) != 1 || result[0] != 9 {
		t.Errorf("expected [9], got %v", result)
	}
}

func TestTopN_distinctValues(t *testing.T) {
	result := TopN([]int{1, 5, 3, 2, 4}, 3, intGetter)

	if len(result) != 3 {
		t.Errorf("expected exactly 3 items, got %d: %v", len(result), result)
	}

	got := make([]int, len(result))
	copy(got, result)
	sort.Ints(got)

	want := []int{3, 4, 5}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("expected top 3 to be %v, got %v", want, got)
			break
		}
	}
}

func TestTopN_customStruct(t *testing.T) {
	type pixel struct{ luma int }

	items := []pixel{{10}, {80}, {50}, {30}, {90}}
	getter := func(p pixel) int { return p.luma }

	result := TopN(items, 2, getter)

	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d: %v", len(result), result)
	}
	for _, r := range result {
		if r.luma < 80 {
			t.Errorf("unexpected low-luma pixel in top 2: %v", r)
		}
	}
}

// When duplicate values sit exactly on the threshold, TopN may return more
// than k items — all tied values are included. This is expected behavior.
func TestTopN_duplicatesAtBoundary_mayReturnMoreThanK(t *testing.T) {
	// top 2 of [1,3,3,5]: threshold = 3rd largest = 3; items >= 3 = [3,3,5]
	result := TopN([]int{1, 3, 3, 5}, 2, intGetter)

	if len(result) < 2 {
		t.Errorf("expected at least 2 items, got %d: %v", len(result), result)
	}
	for _, v := range result {
		if v < 3 {
			t.Errorf("result contains value below threshold: %d", v)
		}
	}
}

func TestTopN_allEqualValues(t *testing.T) {
	result := TopN([]int{5, 5, 5, 5}, 2, intGetter)
	if len(result) < 2 {
		t.Errorf("expected at least 2 items, got %d: %v", len(result), result)
	}
	for _, v := range result {
		if v != 5 {
			t.Errorf("expected all 5s, got %d", v)
		}
	}
}
