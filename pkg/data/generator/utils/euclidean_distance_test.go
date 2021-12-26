package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeEuclideanDistance(t *testing.T) {
	t.Parallel()

	type args struct {
		p1 *Point
		p2 *Point
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"euclidean distance for two different points",
			args{
				&Point{0, 0},
				&Point{3, 4},
			},
			5,
		},
		{
			"euclidean distance between the same point",
			args{
				&Point{1, 1},
				&Point{1, 1},
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeEuclideanDistance(tt.args.p1, tt.args.p2)
			assert.Equal(t, tt.want, got)
		})
	}
}
