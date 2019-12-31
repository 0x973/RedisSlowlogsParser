// Copyright 2019 ShouDong Zheng. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

package concurrentmap

import "sync"

// Map is a thread-safe map
type Map struct {
	items map[interface{}]interface{}
	mu    sync.RWMutex
}

// New returns a new thread-safe map
func New() *Map {
	return &Map{items: make(map[interface{}]interface{})}
}

// Range - Iterate through all the data
func (m *Map) Range(action func(key, value interface{}) bool) {
	m.mu.RLock()
	for key, value := range m.items {
		if !action(key, value) {
			break
		}
	}
	m.mu.RUnlock()
}

// Values - Get all value in items
func (m *Map) Values() []interface{} {
	result := make([]interface{}, 0)
	m.mu.RLock()
	for _, value := range m.items {
		result = append(result, value)
	}
	m.mu.RUnlock()
	return result
}

// Keys - Get all key in items
func (m *Map) Keys() []interface{} {
	result := make([]interface{}, 0)
	m.mu.RLock()
	for key := range m.items {
		result = append(result, key)
	}
	m.mu.RUnlock()
	return result
}

// Get returns the value to which the specified key is mapped.
func (m *Map) Get(key interface{}) (value interface{}, found bool) {
	m.mu.RLock()
	value, found = m.items[key]
	m.mu.RUnlock()
	return
}

// Contains returns true if the specified key exists.
func (m *Map) Contains(key interface{}) bool {
	_, found := m.Get(key)
	return found
}

// Put associates the specified value with the specified key.
func (m *Map) Put(key interface{}, value interface{}) interface{} {
	m.mu.Lock()
	m.items[key] = value
	m.mu.Unlock()
	return value
}

// ComputeIfAbsent check if the specified key is not already associated with a value, attempts to compute its value using the given mapping function and enters it into this map.
func (m *Map) ComputeIfAbsent(key interface{}, compFunction func(key interface{}) interface{}) (value interface{}, computed bool) {
	value, found := m.Get(key)
	if !found {
		value = m.Put(key, compFunction(key))
	}

	return value, !found
}

// Remove the entry associated with the specified key.
func (m *Map) Remove(key interface{}) (found bool) {
	if _, found = m.Get(key); found {
		m.mu.Lock()
		delete(m.items, key)
		defer m.mu.Unlock()
	}
	return
}

// Size returns the number of items in this map
func (m *Map) Size() (size int) {
	m.mu.RLock()
	size = len(m.items)
	m.mu.RUnlock()
	return
}

// Clear - delete all items value
func (m *Map) Clear() {
	m.mu.Lock()
	m.items = nil
	m.items = make(map[interface{}]interface{})
	m.mu.Unlock()
}
