package bucketpoolbestfit

type BucketPool struct {
	Buckets  []int
	InitSize int
}

func (bp *BucketPool) Expand() int {
	defer func() {
		bp.InitSize++
	}()
	return bp.Buckets[bp.InitSize]
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
