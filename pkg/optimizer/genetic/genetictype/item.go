package genetictype

import "github.com/lothar1998/v2x-optimizer/pkg/data"

type Item struct {
	id   int
	size int
}

func NewItem(id, size int) *Item {
	return &Item{id: id, size: size}
}

func (i *Item) ID() int {
	return i.id
}

func (i *Item) Size() int {
	return i.size
}

type itemPoolKey struct {
	itemID   int
	bucketID int
}

type ItemPool struct {
	items map[itemPoolKey]*Item
	data  *data.Data
}

func NewItemPool(data *data.Data) *ItemPool {
	return &ItemPool{items: make(map[itemPoolKey]*Item), data: data}
}

func (ip *ItemPool) Get(itemID, bucketID int) *Item {
	key := itemPoolKey{itemID: itemID, bucketID: bucketID}
	if item, ok := ip.items[key]; ok {
		return item
	}

	item := NewItem(itemID, ip.data.R[itemID][bucketID])
	ip.items[key] = item
	return item
}
