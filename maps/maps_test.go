package maps

import (
	"reflect"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	tests := []struct {
		name   string
		src    map[string]int
		wanted map[string]int
	}{
		{
			name: "normal test",
			src: map[string]int{
				"aa": 1,
				"bb": 2,
			},
			wanted: map[string]int{
				"aa": 1,
				"bb": 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeepCopy(tt.src); !reflect.DeepEqual(got, tt.wanted) {
				t.Errorf("DeepCopy() = %v, want %v", got, tt.wanted)
			}
		})
	}
}
