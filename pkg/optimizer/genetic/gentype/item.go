package gentype

import (
	"sync"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

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
	mu    sync.RWMutex
}

func NewItemPool(data *data.Data) *ItemPool {
	return &ItemPool{items: make(map[itemPoolKey]*Item), data: data}
}

func (ip *ItemPool) Get(itemID, bucketID int) *Item {
	key := itemPoolKey{itemID: itemID, bucketID: bucketID}
	if item := ip.getExistingItem(key); item != nil {
		return item
	}

	return ip.createAndGetNewItem(key, itemID, bucketID)
}

func (ip *ItemPool) getExistingItem(key itemPoolKey) *Item {
	ip.mu.RLock()
	defer ip.mu.RUnlock()

	return ip.items[key]
}

func (ip *ItemPool) createAndGetNewItem(key itemPoolKey, itemID, bucketID int) *Item {
	ip.mu.Lock()
	defer ip.mu.Unlock()

	item := NewItem(itemID, ip.data.R[itemID][bucketID])
	ip.items[key] = item
	return item
}

func (ip *ItemPool) MaxID() int {
	return len(ip.data.R) - 1
}
