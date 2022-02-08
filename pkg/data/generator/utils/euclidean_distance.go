package utils

import "math"

type Point struct {
	X float64
	Y float64
}

func (p Point) Distance(d Point) float64 {
	return math.Sqrt(math.Pow(p.X-d.X, 2) + math.Pow(p.Y-d.Y, 2))
}

func ComputeEuclideanDistance(p1, p2 *Point) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}
