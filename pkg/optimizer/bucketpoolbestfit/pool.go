package bucketpoolbestfit

import "errors"

type bucketPool struct {
	Buckets  []int
	InitSize int
}

func (bp *bucketPool) Expand() (int, error) {
	if bp.InitSize == bp.MaxSize() {
		return 0, errors.New("no items to expand the pool")
	}

	defer func() {
		bp.InitSize++
	}()
	return bp.Buckets[bp.InitSize], nil
}

func (bp *bucketPool) GetBuckets() []int {
	return bp.Buckets[:bp.InitSize]
}

func (bp *bucketPool) Size() int {
	return bp.InitSize
}

func (bp *bucketPool) MaxSize() int {
	return len(bp.Buckets)
}
