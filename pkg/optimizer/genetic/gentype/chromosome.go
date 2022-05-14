package gentype

type Chromosome struct {
	buckets      []*Bucket
	idsToBuckets map[int]*Bucket
}

func NewChromosome(size int) *Chromosome {
	return &Chromosome{buckets: make([]*Bucket, size), idsToBuckets: make(map[int]*Bucket)}
}

func (c *Chromosome) At(index int) *Bucket {
	return c.buckets[index]
}

func (c *Chromosome) Slice(left, right int) []*Bucket {
	return c.buckets[left:right]
}

func (c *Chromosome) SetAt(index int, bucket *Bucket) {
	c.buckets[index] = bucket
	c.idsToBuckets[bucket.ID()] = bucket
}

func (c *Chromosome) Append(bucket *Bucket) {
	c.buckets = append(c.buckets, bucket)
	c.idsToBuckets[bucket.ID()] = bucket
}

func (c *Chromosome) ContainsBucket(bucketID int) bool {
	_, ok := c.idsToBuckets[bucketID]
	return ok
}

func (c *Chromosome) Len() int {
	return len(c.buckets)
}
