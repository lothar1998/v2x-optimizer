package genetictype

import (
	"errors"
	"fmt"
)

var (
	ErrBucketCapacityLimit = errors.New("cannot add item to bucket due to capacity restriction")
	ErrItemAlreadyInBucket = errors.New("item is already in bucket")
)

type Bucket struct {
	id          int
	capacity    int
	freeSpace   int
	itemMapping map[int]*Item
}

func NewBucket(id int, capacity int) *Bucket {
	return &Bucket{
		id:          id,
		capacity:    capacity,
		freeSpace:   capacity,
		itemMapping: make(map[int]*Item),
	}
}

func (b *Bucket) AddItem(item *Item) error {
	if _, ok := b.itemMapping[item.ID()]; ok {
		return fmt.Errorf("%w: bucketId=%d, itemId=%d", ErrItemAlreadyInBucket, b.id, item.ID())
	} else if b.freeSpace < item.Size() {
		return fmt.Errorf(
			"%w: bucketId=%d, itemId=%d, additional space needed=%d",
			ErrBucketCapacityLimit,
			b.id,
			item.ID(),
			item.Size()-b.freeSpace,
		)
	}

	b.freeSpace -= item.Size()
	b.itemMapping[item.ID()] = item
	return nil
}

func (b *Bucket) DeepCopy() *Bucket {
	return &Bucket{
		id:          b.id,
		capacity:    b.capacity,
		freeSpace:   b.freeSpace,
		itemMapping: b.copyMapping(),
	}
}

func (b *Bucket) copyMapping() map[int]*Item {
	newMapping := make(map[int]*Item)
	for id, item := range b.itemMapping {
		newMapping[id] = item
	}
	return newMapping
}

func (b *Bucket) IsEmpty() bool {
	return len(b.itemMapping) == 0
}

func (b *Bucket) ID() int {
	return b.id
}

func (b *Bucket) Capacity() int {
	return b.capacity
}

func (b *Bucket) FreeSpace() int {
	return b.freeSpace
}

func (b *Bucket) Map() map[int]*Item {
	return b.itemMapping
}
