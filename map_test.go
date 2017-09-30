package gosh

import (
	"testing"
)

func TestMapInitialization(t *testing.T) {
	size := uint32(5)
	m := Initialize(size)
	assert(t, m != nil, "expect map to be initialized and not nil")
	equals(t, m.size, size)
}

func TestMapSet(t *testing.T) {
	size := uint32(5)
	m := Initialize(size)
	status := m.Set("test", "value")
	equals(t, status, true)
	assert(t, m.len == 1, "expect the len to be 1")
}

func TestMapSetSameKey(t *testing.T) {
	size := uint32(5)
	m := Initialize(size)
	status := m.Set("key1", "value")
	equals(t, status, true)
	val := m.Get("key1")
	equals(t, val.(string), "value")

	status = m.Set("key1", "value2")
	equals(t, status, true)
	assert(t, m.len == 1, "expect the len to stay unchanged")

	val = m.Get("key1")
	assert(t, val.(string) == "value2", "expect to receive the updated value")
}

func TestMapSetOverflow(t *testing.T) {
	size := uint32(2)
	m := Initialize(size)

	m.Set("key1", "value1")
	m.Set("key2", "value2")
	equals(t, m.len, uint32(2))

	status := m.Set("key3", "value3")
	assert(t, status == false, "expect failure to be returned when there's an overflow")
}

func TestMapGetValidKey(t *testing.T) {
	size := uint32(5)
	m := Initialize(size)

	value := "value1"
	m.Set("key1", value)
	val := m.Get("key1")

	equals(t, val.(string), value)
}

func TestMapGetInvalidKey(t *testing.T) {
	size := uint32(5)
	m := Initialize(size)

	value := "value1"
	m.Set("key1", value)
	val := m.Get("key2")

	equals(t, nil, val)
}

func TestMapDeletionValidKey(t *testing.T) {
	size := uint32(5)
	m := Initialize(size)

	m.Set("key1", "value1")
	m.Set("key2", "value2")

	val := m.Get("key2")

	equals(t, "value2", val.(string))

	val2 := m.Delete("key2")
	assert(t, val2 == val, "expect the deleted value to be returned")

	val = m.Get("key2")
	assert(t, val == nil, "expect key to be deleted")

	// expect len to be decreased by 1
	equals(t, m.len, uint32(1))
}

func TestMapDeletionInvalidKey(t *testing.T) {
	size := uint32(5)
	m := Initialize(size)

	m.Set("key1", "value1")

	val := m.Delete("key2")
	assert(t, val == nil, "expect nil to be returned for invalid key")

	// expect the len to stay unchanged
	equals(t, m.len, uint32(1))
}

func TestLoadFactor(t *testing.T) {
	m := Initialize(2)

	m.Set("key1", "value1")

	assert(t, m.Load() == 0.5, "expect the load to be half the size")

	m.Set("key2", "value2")

	assert(t, m.Load() == 1, "expect the load to be 1")
}
