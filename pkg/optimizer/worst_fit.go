package optimizer

import (
	"container/heap"
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

type bucket struct {
	Index     int
	LeftSpace int
	heapIndex int
}

type priorityBucketQueue []*bucket

func (pq priorityBucketQueue) Len() int {
	return len(pq)
}

func (pq priorityBucketQueue) Less(i, j int) bool {
	return pq[i].LeftSpace > pq[j].LeftSpace
}

func (pq priorityBucketQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].heapIndex = i
	pq[j].heapIndex = j
}

func (pq *priorityBucketQueue) Push(x interface{}) {
	n := len(*pq)
	b := x.(*bucket)
	b.heapIndex = n
	*pq = append(*pq, b)
}

func (pq *priorityBucketQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	b := old[n-1]
	old[n-1] = nil
	b.heapIndex = -1
	*pq = old[:n-1]
	return b
}

func (pq *priorityBucketQueue) Decrease(b *bucket, decreaseAmount int) {
	b.LeftSpace -= decreaseAmount
	heap.Fix(pq, b.heapIndex)
}

func (pq priorityBucketQueue) At(index int) *bucket {
	return pq[index]
}

// WorstFit is an optimizer that implements the worst-fit algorithm expanded to solve the bin packing problem
// with heterogeneous bins and items with different sizes that depend on the bin.
// If the item cannot be assigned to the bucket from heap top, it is tried to assign it to any of the other buckets
// in heap array order. This means that if the least filled bucket is not suitable
// the algorithm doesn't look for the second least filled bucket but gets next from the current heap array.
type WorstFit struct{}

func (w WorstFit) Optimize(ctx context.Context, data *data.Data) (*Result, error) {
	v := len(data.R)
	n := len(data.MRB)
	sequence := make([]int, v)
	leftSpace := data.MRB

	pq := make(priorityBucketQueue, n)
	for i, bucketSize := range leftSpace {
		pq[i] = &bucket{
			Index:     i,
			LeftSpace: bucketSize,
			heapIndex: i,
		}
	}

	heap.Init(&pq)

	for i := 0; i < v; i++ {
		j := 0
		for ; j < pq.Len(); j++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			b := pq.At(j)

			if data.R[i][b.Index] < b.LeftSpace {
				sequence[i] = b.Index
				pq.Decrease(b, data.R[i][b.Index])
				break
			}
		}

		if j >= pq.Len() {
			return nil, ErrCannotAssignToBucket
		}
	}

	return toResult(sequence, n), nil
}
