package bucketpoolbestfit

import "errors"

type BucketPool struct {
	Buckets  []int
	InitSize int
}

func (bp *BucketPool) Expand() (int, error) {
	if bp.InitSize == bp.MaxSize() {
		return 0, errors.New("no items to expand the pool")
	}

	defer func() {
		bp.InitSize++
	}()
	return bp.Buckets[bp.InitSize], nil
}

func (bp *BucketPool) GetBuckets() []int {
	return bp.Buckets[:bp.InitSize]
}

func (bp *BucketPool) Size() int {
	return bp.InitSize
}

func (bp *BucketPool) MaxSize() int {
	return len(bp.Buckets)
}
