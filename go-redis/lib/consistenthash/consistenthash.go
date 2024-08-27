package consistenthash

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32

type NodeMap struct {
	hashFunc    HashFunc
	nodeHashes  []int
	nodeHashMap map[int]string
}

func NewNodeMap(fn HashFunc) *NodeMap {
	mp := &NodeMap{
		hashFunc:    fn,
		nodeHashMap: make(map[int]string),
	}
	if mp.hashFunc == nil {
		mp.hashFunc = crc32.ChecksumIEEE
	}
	return mp
}

func (mp *NodeMap) IsEmpty() bool {
	return len(mp.nodeHashMap) == 0
}

func (mp *NodeMap) AddNode(keys ...string) {
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(mp.hashFunc([]byte(key)))
		mp.nodeHashes = append(mp.nodeHashes, hash)
		mp.nodeHashMap[hash] = key
	}
	sort.Ints(mp.nodeHashes)
}

func (mp *NodeMap) PickNode(key string) string {
	if mp.IsEmpty() {
		return ""
	}
	hash := int(mp.hashFunc([]byte(key)))
	index := sort.Search(len(mp.nodeHashes), func(i int) bool {
		return mp.nodeHashes[i] >= hash
	})
	if index == len(mp.nodeHashes) {
		index = 0
	}
	return mp.nodeHashMap[mp.nodeHashes[index]]
}
