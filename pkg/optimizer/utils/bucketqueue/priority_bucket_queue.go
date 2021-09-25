package bucketqueue

import "container/heap"

type Bucket struct {
	Index     int
	LeftSpace int
	HeapIndex int
}

type PriorityBucketQueue []*Bucket

func NewPriority(bucketSizes []int) PriorityBucketQueue {
	pq := make(PriorityBucketQueue, len(bucketSizes))

	for i, bucketSize := range bucketSizes {
		pq[i] = &Bucket{
			Index:     i,
			LeftSpace: bucketSize,
			HeapIndex: i,
		}
	}

	heap.Init(&pq)
	return pq
}

func (pq PriorityBucketQueue) Len() int {
	return len(pq)
}

func (pq PriorityBucketQueue) Less(i, j int) bool {
	return pq[i].LeftSpace > pq[j].LeftSpace
}

func (pq PriorityBucketQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].HeapIndex = i
	pq[j].HeapIndex = j
}

func (pq *PriorityBucketQueue) Push(x interface{}) {
	n := len(*pq)
	b := x.(*Bucket)
	b.HeapIndex = n
	*pq = append(*pq, b)
}

func (pq *PriorityBucketQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	b := old[n-1]
	old[n-1] = nil
	b.HeapIndex = -1
	*pq = old[:n-1]
	return b
}

func (pq *PriorityBucketQueue) PopEmptiestBucket() *Bucket {
	return heap.Pop(pq).(*Bucket)
}

func (pq *PriorityBucketQueue) PushBucket(bucket *Bucket) {
	heap.Push(pq, bucket)
}

func (pq *PriorityBucketQueue) Decrease(b *Bucket, decreaseAmount int) {
	b.LeftSpace -= decreaseAmount
	heap.Fix(pq, b.HeapIndex)
}

func (pq PriorityBucketQueue) At(index int) *Bucket {
	return pq[index]
}
