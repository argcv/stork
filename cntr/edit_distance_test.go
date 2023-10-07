package cntr

import (
	"fmt"
	"testing"

	"github.com/argcv/stork/assert"
)

func TestGetEditDistance(t *testing.T) {
	tests := []struct {
		name      string
		s1        string
		s2        string
		threshold int
		want      int
	}{
		{
			name:      "empty",
			s1:        "",
			s2:        "",
			threshold: 0,
			want:      0,
		},
		{
			name:      "empty-vs-nonempty",
			s1:        "",
			s2:        "x",
			threshold: 0,
			want:      -1,
		},
		{
			name:      "equal",
			s1:        "abc",
			s2:        "abc",
			threshold: 1,
			want:      0,
		},
		{
			name:      "non-equal",
			s1:        "bbc",
			s2:        "abc",
			threshold: 2,
			want:      1,
		},
		{
			name:      "non-equal",
			s1:        "abc",
			s2:        "abd",
			threshold: 3,
			want:      1,
		},
		{
			name:      "too-far",
			s1:        "hello",
			s2:        "world",
			threshold: 3,
			want:      -1,
		},
		{
			name:      "not-too-far",
			s1:        "hello",
			s2:        "world",
			threshold: 5,
			want:      4,
		},
		{
			name:      "unicode-distance",
			s1:        "北京西门",
			s2:        "北海西",
			threshold: 5,
			want:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetEditDistanceT(tt.s1, tt.s2, tt.threshold)
			assert.ExpectEQ(t, got, tt.want,
				fmt.Sprintf("GetEditDistanceT(%v, %v) = %v, want %v", tt.s1, tt.s2, got, tt.want))
		})
	}
}
