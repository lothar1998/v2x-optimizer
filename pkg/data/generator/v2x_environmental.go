package generator

import (
	"math"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

const (
	squareSideLength = 1
	bigP             = 2
	alpha            = 3.5
	rate             = 1 * 1024
)

func GenerateV2XEnvironmental(itemCount, _, bucketCount, maxBucketSize int) *data.Data {
	itemSizes := generateItemSizesV2XEnvironmental(itemCount, bucketCount)
	bucketSizes := generateBucketsWithSizes(bucketCount, maxBucketSize)

	return &data.Data{R: itemSizes, MRB: bucketSizes}
}

func GenerateV2XEnvironmentalConstantBucketSize(itemCount, _, bucketCount, bucketSize int) *data.Data {
	itemSizes := generateItemSizesV2XEnvironmental(itemCount, bucketCount)
	bucketSizes := generateBucketsOfConstantSize(bucketCount, bucketSize)

	return &data.Data{R: itemSizes, MRB: bucketSizes}
}

func generateItemSizesV2XEnvironmental(itemCount, bucketCount int) [][]int {
	vehiclePoints := generateVehiclePoints(itemCount)
	stationPoints := generateStationPoints(bucketCount)

	itemSizes := make([][]int, itemCount)
	for i := range itemSizes {
		itemSizes[i] = make([]int, bucketCount)
		for j := range itemSizes[i] {
			itemSizes[i][j] = computeR(vehiclePoints[i], stationPoints[j])
		}
	}

	return itemSizes
}

func generateStationPoints(n int) []point {
	stationsCountInOneDimension := int(math.Ceil(math.Sqrt(float64(n))))
	radius := squareSideLength / float64(2*stationsCountInOneDimension)
	totalStationCount := int(math.Pow(float64(stationsCountInOneDimension), 2))

	stationIndices := random.Perm(totalStationCount)[0:n]

	stationPoints := make([]point, len(stationIndices))
	for i, index := range stationIndices {
		stationPoints[i] = toStationPoint(radius, stationsCountInOneDimension, index)
	}

	return stationPoints
}

func generateVehiclePoints(v int) []point {
	vehiclePoints := make([]point, v)
	for i := range vehiclePoints {
		x := random.Float64() * squareSideLength
		y := random.Float64() * squareSideLength
		vehiclePoints[i] = point{X: x, Y: y}
	}

	return vehiclePoints
}

func toStationPoint(radius float64, stationCount, i int) point {
	x := radius * float64((1+2*i)%(2*stationCount))
	y := radius + 2*(math.Floor(float64(i)/float64(stationCount)))*radius
	return point{X: x, Y: y}
}

func computeR(p1, p2 point) int {
	d := p1.Distance(p2)
	s := bigP * math.Pow(d, -alpha)

	return int(math.Ceil(rate / toRBDataRate(s)))
}

func toRBDataRate(s float64) float64 {
	switch {
	case s <= 0.4:
		return 101.07
	case s <= 2.5:
		return 147.34
	case s <= 4.5:
		return 197.53
	case s <= 6.5:
		return 248.07
	case s <= 8.5:
		return 321.57
	case s <= 10.3:
		return 404.26
	case s <= 12.3:
		return 458.72
	case s <= 14.2:
		return 558.15
	case s <= 15.9:
		return 655.59
	case s <= 17.8:
		return 759.93
	case s <= 19.8:
		return 859.35
	default:
		return 933.19
	}
}

type point struct {
	X float64
	Y float64
}

func (p point) Distance(d point) float64 {
	return math.Sqrt(math.Pow(p.X-d.X, 2) + math.Pow(p.Y-d.Y, 2))
}
