package gosh

import "hash/fnv"

// node represents an entry in the hash map.
// it is stored at an index in the bucket and
// contains a pointer to the next element in
// the chain.
type node struct {
	hash  uint32
	key   string
	value interface{}
	next  *node
}

// Map represents the hash map object.
type Map struct {
	// len keeps track of the number of elements in
	// the map.
	len uint32

	// size represents the capacity of the map
	// defined at the time of initialization.
	size uint32

	// buckets stores the slice of buckets for
	// the map.
	buckets []*node
}

// Initialize creates a map of the defined size and
// returns it. It accepts the size as the input parameter.
func Initialize(s uint32) *Map {
	return &Map{
		len:     0,
		size:    s,
		buckets: make([]*node, s),
	}
}

// Set inserts an object into the map. It returns true if the
// operation was successful.
// If the map has reached its size, then the set operation will
// return false. One can call the Load() method to check for
// the load factor of the map to verify.
func (m *Map) Set(key string, value interface{}) bool {
	keyHash, bucketIndex := getHashAndIndex(key, m.size)

	node := getValueFromBucketChain(m.buckets[bucketIndex], keyHash)

	// key is not present in the map and the
	// map has reached its limit.
	if node == nil && m.len == m.size {
		return false
	}

	// update the value for the given key if the node already
	// existed in the map.
	if node != nil {
		node.value = value
		return true
	}

	// insert a new node in the map.
	m.addNode(bucketIndex, createNode(keyHash, key, value))
	m.len++
	return true
}

// Get returns the object associated with the given key.
// It will return nil if the key doesn't exist in the map.
func (m *Map) Get(key string) interface{} {
	keyHash, bucketIndex := getHashAndIndex(key, m.size)
	n := getValueFromBucketChain(m.buckets[bucketIndex], keyHash)
	if n == nil {
		return nil
	}

	return n.value
}

// Delete deletes the entry from the map for the given key.
// It returns the value of the deleted object if it's present
// in the map. If not, it will return nil.
func (m *Map) Delete(key string) interface{} {
	keyHash, bucketIndex := getHashAndIndex(key, m.size)

	n := m.buckets[bucketIndex]

	// no chain of nodes exists for the given bucket index.
	if n == nil {
		return nil
	}

	// the key is the first node in the chain.
	if n.hash == keyHash {
		m.buckets[bucketIndex] = n.next
		m.len--

		// free the node
		n.next = nil
		return n.value
	}

	nex := n.next

	for nex != nil {
		// detach node from the chain if found.
		if nex.hash == keyHash {
			n.next = nex.next
			m.len--
			nex.next = nil
			return nex.value
		}
		n = nex
		nex = nex.next
	}

	return nil
}

// Load returns the load factor of the map. It is defined by
// (items in hashmap)/(size of hashmap).
// This value will never be greater than 1(represents a full map).
func (m *Map) Load() float64 {
	return float64(m.len) / float64(m.size)
}

// getHash computes the hash value using the FNV-1a hash function.
func getHashAndIndex(key string, size uint32) (uint32, uint32) {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32(), h.Sum32() % size
}

// getValueFromBucketChain searches for the node in the chain
// with the specified key hash.
func getValueFromBucketChain(nodeChain *node, keyHash uint32) *node {
	var n *node

	for nodeChain != nil {
		if nodeChain.hash == keyHash {
			n = nodeChain
			break
		}
		nodeChain = nodeChain.next
	}

	return n
}

// createNode defines the constructor for the node struct.
func createNode(h uint32, k string, val interface{}) *node {
	return &node{
		hash:  h,
		key:   k,
		value: val,
	}
}

// addNode inserts the new node at the defined bucket index.
func (m *Map) addNode(ind uint32, n *node) {
	if m.buckets[ind] == nil {
		m.buckets[ind] = n
		return
	}

	currNode := m.buckets[ind]
	for {
		if currNode.next == nil {
			currNode.next = n
			return
		}
	}
}
