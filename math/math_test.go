package math

import (
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		name   string
		x      float32
		wanted float32
	}{
		{
			name:   "normal test 1",
			x:      -1,
			wanted: 1,
		},
		{
			name:   "normal test 2",
			x:      1,
			wanted: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.x); got != tt.wanted {
				t.Errorf("Abs() = %v, want %v", got, tt.wanted)
			}
		})
	}
}
