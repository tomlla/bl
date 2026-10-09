package main

import "testing"

func TestStepLevel(t *testing.T) {
	tests := []struct {
		name            string
		level, maxLevel uint32
		up              bool
		want            uint32
	}{
		{"up amdgpu", 60000, 65535, true, 63276},
		{"up clamps to max", 65000, 65535, true, 65535},
		{"down amdgpu", 60000, 65535, false, 56724},
		{"down clamps to zero", 1000, 65535, false, 0},
		{"tiny max still moves", 5, 10, true, 6},
	}
	for _, tt := range tests {
		if got := stepLevel(tt.level, tt.maxLevel, tt.up); got != tt.want {
			t.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
		}
	}
}
