// Package lrucache provides a bounded, generation-guarded LRU cache with
// single-flight fills, built for in-memory mirrors of database rows that are
// kept consistent by explicit eviction rather than TTLs.
//
// Design points, all load-bearing for correctness:
//
//   - Fills are single-flight: concurrent callers for the same key wait on
//     one leader and share its result and error. Errors are propagated to
//     every waiter but never cached, so a failed load is retried by the next
//     caller. A panicking loader releases its waiters with an error and
//     re-raises on the leader's own stack.
//   - A generation counter increments on every explicit eviction or flush.
//     A fill records the generation before running its loader and its
//     install is discarded if the generation moved, so a fill racing an
//     eviction can never re-install the value the eviction just removed.
//     Capacity (LRU) eviction does not increment it: dropping a still-valid
//     entry for space is not an invalidation.
//   - An optional secondary index maps a caller-chosen string (typically a
//     row's primary key) to the cache key holding it, so write paths that
//     only know the row identity can evict without recomputing the key.
//   - An optional validator runs on every hit; an entry it rejects is
//     removed and reported as a miss, so callers fall through to their full
//     load path (used, for example, to treat expired credentials as misses
//     while keeping all expiry policy in the caller).
//
// The cache owns its own locking and never holds a lock across a loader
// call, so loaders are free to perform database reads or network I/O. It
// must not be guarded by any caller-side lock that is itself held across
// I/O; give the cache exclusive ownership of its consistency instead.
//
// All methods are safe on a nil receiver: reads miss, evictions no-op, and
// Fill degrades to calling the loader directly. This lets consumers treat
// an absent cache as "caching disabled" without branching.
package lrucache

import (
	"container/list"
	"fmt"
	"sync"
)

// Loader produces the value for a cache key on a miss. It returns the value,
// the secondary-index key to register for it (empty for none), and an error.
// Errors are shared with concurrent waiters but never cached.
type Loader[V any] func() (value V, indexKey string, err error)

// inflightCall represents one in-progress fill for a cache key. Concurrent
// callers for the same key wait on done and share result/err.
type inflightCall[V any] struct {
	done   chan struct{}
	result V
	err    error
}

type entry[V any] struct {
	key      string
	indexKey string
	value    V
}

// Cache is a bounded LRU cache with single-flight fills. Construct with New;
// the zero value is not usable (but a nil *Cache is, see the package doc).
type Cache[V any] struct {
	mu       sync.Mutex
	capacity int
	items    map[string]*list.Element
	order    *list.List // front = most recently used
	// byIndex maps a secondary-index key to the cache key holding it.
	byIndex map[string]string
	// generation increments on every explicit eviction or flush; see the
	// package doc for the fill-race guarantee it provides.
	generation uint64
	// validator, when non-nil, is consulted on every hit; entries it
	// rejects are removed and reported as misses.
	validator func(V) bool

	inflightMu sync.Mutex
	inflight   map[string]*inflightCall[V]
}

// Option configures a Cache at construction time.
type Option[V any] func(*Cache[V])

// WithValidator installs a hit-time validity check: entries for which valid
// returns false are removed and reported as misses.
func WithValidator[V any](valid func(V) bool) Option[V] {
	return func(c *Cache[V]) {
		c.validator = valid
	}
}

// New constructs a Cache bounded to capacity entries. A non-positive
// capacity panics: the bound is part of the type's contract, and silently
// defaulting it would hide a caller bug.
func New[V any](capacity int, opts ...Option[V]) *Cache[V] {
	if capacity <= 0 {
		panic("lrucache: capacity must be positive")
	}
	c := &Cache[V]{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
		byIndex:  make(map[string]string),
		inflight: make(map[string]*inflightCall[V]),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Get returns the cached value for key. An entry rejected by the validator
// is removed and reported as a miss.
func (c *Cache[V]) Get(key string) (V, bool) {
	var zero V
	if c == nil {
		return zero, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.items[key]
	if !ok {
		return zero, false
	}
	ent := elem.Value.(*entry[V])
	if c.validator != nil && !c.validator(ent.value) {
		c.removeLocked(elem)
		return zero, false
	}
	c.order.MoveToFront(elem)
	return ent.value, true
}

// Fill runs load for key with single-flight deduplication: concurrent
// callers for the same key wait for one leader and share its result and
// error. A successful result is cached unless an eviction or flush landed
// while the load was in flight; errors are propagated but never cached.
func (c *Cache[V]) Fill(key string, load Loader[V]) (V, error) {
	if c == nil {
		v, _, err := load()
		return v, err
	}
	c.inflightMu.Lock()
	if call, ok := c.inflight[key]; ok {
		c.inflightMu.Unlock()
		<-call.done
		return call.result, call.err
	}
	call := &inflightCall[V]{done: make(chan struct{})}
	c.inflight[key] = call
	c.inflightMu.Unlock()

	// Deferred so a panicking loader still releases the waiters and clears
	// the inflight slot instead of wedging this key for the process
	// lifetime. Waiters must never mistake a panicked load for success, so
	// an error is set before the channel closes and the panic is re-raised
	// for the leader's own stack.
	defer func() {
		if r := recover(); r != nil {
			call.err = fmt.Errorf("lrucache: loader panicked: %v", r)
			close(call.done)
			c.inflightMu.Lock()
			delete(c.inflight, key)
			c.inflightMu.Unlock()
			panic(r)
		}
		close(call.done)
		c.inflightMu.Lock()
		delete(c.inflight, key)
		c.inflightMu.Unlock()
	}()

	gen := c.currentGeneration()
	value, indexKey, err := load()
	call.result, call.err = value, err

	if err == nil {
		c.setIfGeneration(key, indexKey, value, gen)
	}
	return value, err
}

func (c *Cache[V]) currentGeneration() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.generation
}

// setIfGeneration installs value under key unless the generation moved since
// gen was read, meaning an eviction or flush invalidated state while the
// loader was running.
func (c *Cache[V]) setIfGeneration(key, indexKey string, value V, gen uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.generation != gen {
		return
	}
	if elem, ok := c.items[key]; ok {
		ent := elem.Value.(*entry[V])
		if ent.indexKey != indexKey && ent.indexKey != "" {
			if mapped, ok := c.byIndex[ent.indexKey]; ok && mapped == key {
				delete(c.byIndex, ent.indexKey)
			}
		}
		ent.value = value
		ent.indexKey = indexKey
		if indexKey != "" {
			c.byIndex[indexKey] = key
		}
		c.order.MoveToFront(elem)
		return
	}
	if c.order.Len() >= c.capacity {
		if tail := c.order.Back(); tail != nil {
			c.removeLocked(tail)
		}
	}
	elem := c.order.PushFront(&entry[V]{key: key, indexKey: indexKey, value: value})
	c.items[key] = elem
	if indexKey != "" {
		c.byIndex[indexKey] = key
	}
}

// removeLocked drops an entry from the list, the key map, and the secondary
// index. Caller holds c.mu.
func (c *Cache[V]) removeLocked(elem *list.Element) {
	ent := elem.Value.(*entry[V])
	c.order.Remove(elem)
	delete(c.items, ent.key)
	if ent.indexKey != "" {
		if mapped, ok := c.byIndex[ent.indexKey]; ok && mapped == ent.key {
			delete(c.byIndex, ent.indexKey)
		}
	}
}

// Evict removes the entry for an exact key, if present.
func (c *Cache[V]) Evict(key string) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.generation++
	if elem, ok := c.items[key]; ok {
		c.removeLocked(elem)
	}
}

// EvictByIndex removes the entry registered under the given secondary-index
// key, if any.
func (c *Cache[V]) EvictByIndex(indexKey string) {
	if c == nil || indexKey == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.generation++
	if key, ok := c.byIndex[indexKey]; ok {
		if elem, ok := c.items[key]; ok {
			c.removeLocked(elem)
		}
	}
}

// EvictWhere removes every entry whose cache key matches pred. A linear
// sweep under the lock; intended for rare bulk invalidations, not hot
// paths.
func (c *Cache[V]) EvictWhere(pred func(key string) bool) {
	if c == nil || pred == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.generation++
	for elem := c.order.Front(); elem != nil; {
		next := elem.Next()
		if pred(elem.Value.(*entry[V]).key) {
			c.removeLocked(elem)
		}
		elem = next
	}
}

// Flush drops every cached entry.
func (c *Cache[V]) Flush() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.generation++
	c.items = make(map[string]*list.Element)
	c.byIndex = make(map[string]string)
	c.order.Init()
}

// Len reports the number of cached entries.
func (c *Cache[V]) Len() int {
	if c == nil {
		return 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}
