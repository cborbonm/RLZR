package rlzr

import (
	"sync"
)

//type pStateShared struct {
//	items        map[string]*packet_state
//	sync.RWMutex // Read Write mutex, guards access to internal map.
//}

//type pState []*pStateShared



/* NEW START
 * This file includes code from 
 * https://github.com/streamrail/concurrent-map/blob/master/concurrent_map.go 
 * which is licensed under the MIT license (Copyright (c) 2014 streamrail).
 */

var SHARD_COUNT = 32

type pStateShared struct {
	items          map[string]*packet_state
	sync.RWMutex   // Read Write mutex, guards access to internal map.
}

/* Thread-safe map that maps strings to packets. 
 * This map avoid mutex bottlenecks by dividing up
 * into SHARD_COUNT shards. */
/* TODO: change any references from: type pState []*pStateShared to: */
type rCMap []*rCMapShard

/* A "thread" safe string-to-packet_state map. */
type rCMapShard struct {
	items          map[string]*packet_state
	sync.RWMutex   // Read-Write mutex, guards access to internal map.
}

/* Creates a new concurrent map. */
func New() rCMap {
	m := make(rCMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &rCMapShard{items: make(map[string]*packet_state)}
	}
	return m
}

/* Returns shard for given key. */
func (m rCMap) GetShard(key string) *rCMapShard {
	return m[uint(fnv32(key)) % uint(SHARD_COUNT)]
}

/* Sets the given value under the specified key. */
func (m *rCMap) Set(key string, state *packet_state) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = state
	shard.Unlock()
}

/* Retrieves an element from map under given key. 
 * Returns value AND whether it ws found or not. */
func (m rCMap) Get(key string) (state *packet_state, ok bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

/* Returns the number of elements within the map. */
func (m rCMap) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

/* Looks up an item under specified key. */
func (m *rCMap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

/* Removes an element from the map. */
func (m *rCMap) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

/* Checks if map is empty. */
func (m *rCMap) IsEmpty() bool {
	return m.Count() == 0
}

/* Helper function for getting the appropriate shard within the map (GetShard). */
func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
