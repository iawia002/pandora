package array

import (
	"testing"
)

func TestItemInArray(t *testing.T) {
	tests := []struct {
		name   string
		item   int
		items  []int
		wanted bool
	}{
		{
			name:   "int in array test 1",
			item:   1,
			items:  []int{1, 2},
			wanted: true,
		},
		{
			name:   "int in array test 2",
			item:   1,
			items:  []int{2, 3},
			wanted: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ItemInArray(tt.item, tt.items); got != tt.wanted {
				t.Errorf("ItemInArray() = %v, want %v", got, tt.wanted)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name   string
		items  []int
		wanted int
	}{
		{
			name:   "normal test",
			items:  []int{3, 1, 2},
			wanted: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.items...); got != tt.wanted {
				t.Errorf("Min() = %v, want %v", got, tt.wanted)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name   string
		items  []int
		wanted int
	}{
		{
			name:   "normal test",
			items:  []int{1, 2, 3},
			wanted: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.items...); got != tt.wanted {
				t.Errorf("Max() = %v, want %v", got, tt.wanted)
			}
		})
	}
}
