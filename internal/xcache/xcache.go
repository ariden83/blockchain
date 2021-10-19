// Package xcache provides an in-memory LRU cache with extra features inspired from nginx caching.
// It is based on https://github.com/karlseguin/ccache for the basic caching features
// (LRU, concurrency optimizations).
//
// It adds cache locking (prevents 2 concurrent fetches for the same item),
// infinite serving of stale values in case of fetch errors,
// asynchronous refresh of stale values,
// concurrency-limited refresh fetchers,
// and negative cache.
//
// Its usage makes use of a single function Fetch() (no Get()/Set()), which is provided
// with a closure capturing the parameters necessary to fetch for the given key.
package xcache

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/karlseguin/ccache"
)

// Cache is the type supporting the local caching of objects.
type Cache struct {
	posCache     *ccache.Cache // positive cache: store valid entries
	posSize      int32         // how many max cached entries
	posPruneSize int32         // how many entries to evict on cache full
	posTTL       time.Duration // how long until a positive entry is considered stale

	negCache     *ccache.Cache // negative cache: store currently invalid entries
	negSize      int32         // how many max neg cached entries
	negPruneSize int32         // how many entries to evict on cache full
	negTTL       time.Duration // how long until a neg entry is considered stale

	fetching  map[string]struct{} // is an item being fetched?
	fetchLock sync.Mutex          // guard access to "fetching" map
	fetchCond *sync.Cond          // for waking up goroutines waiting for a fetcher

	fetchQueue chan fetchReq       // queue for async fetches
	queued     map[string]struct{} // is an item already queued?
	queuedLock sync.Mutex          // guard access to "queued" map

	staleFetchers  int // number of fetcher goroutines
	staleQueueSize int // size of the chan storing async fetch request

	staleValidator func(interface{}, time.Duration) bool // do serve stale item given its stale age?
	canUseStale    bool                                  // allow serving stale content

	maxFetchers  int           // max number of concurrent fetches
	fetchLimiter chan struct{} // used as a semaphore for concurrency limit

	// instrumentation
	hits         uint64 // cache hit counter
	requests     uint64 // requests (hit+miss) counter
	newFetches   uint64 // fetch counter for item not yet in cache
	staleFetches uint64 // fetch counter for refreshing expired items
}

// Fetcher is the type of the closure passed to Fetch() for fetching the desired object if missing or stale.
//
// Depending on the boolean validity and error returned by the closure,
// the entry will end in the positive or the negative cache (and possibly removed from the other).
//
// If non-nil error, entry will be stored in negative cache (but positive entry will be untouched -
// use this for transient backend failures);
// else if validity is false, will be stored in negative cache and positive entry will be deleted
// (use this case when the item is not "positive" anymore - has invalid state or has been removed);
// else (nil error nil and true validity) store in positive cache and remove neg entry
type Fetcher func() (interface{}, bool, error)

// fetchReq stores a fetch request for async refresh.
type fetchReq struct {
	key string
	f   Fetcher
}

// negCacheEntry stores an invalid fetch result for negative caching
type negCacheEntry struct {
	x   interface{}
	err error
}

// Option is the type of option passed to the constructor.
type Option func(c *Cache)

// WithSize sets the max cache size (in number of objets). When the cache is full,
// a portion of the least recently used objets will be evicted.
// Default: 5000
func WithSize(n int32) Option {
	return func(c *Cache) {
		c.posSize = n
	}
}

// WithPruneSize sets the number of entries to evict when the positive cache is full.
// Default: posSize/20
func WithPruneSize(n int32) Option {
	return func(c *Cache) {
		c.posPruneSize = n
	}
}

// WithTTL sets the duration upon which an object will be deemed expired.
// In this case, a background refresh will occur.
// Default: 60 * time.Second
func WithTTL(t time.Duration) Option {
	return func(c *Cache) {
		c.posTTL = t
	}
}

// WithNegSize sets the size of the negative cache, used for storing errors and invalid objects.
// Default: 500
func WithNegSize(n int32) Option {
	return func(c *Cache) {
		c.negSize = n
	}
}

// WithNegPruneSize sets the number of entries to evict when the negative cache is full.
// Default: negSize/20
func WithNegPruneSize(n int32) Option {
	return func(c *Cache) {
		c.negPruneSize = n
	}
}

// WithNegTTL sets the duration of a negative cache entry.
// Default: 5 * time.Second
func WithNegTTL(t time.Duration) Option {
	return func(c *Cache) {
		c.negTTL = t
	}
}

// WithStaleFetchers sets the number of fetchers in the pool for async fetch of stale entries.
// Default: 3
func WithStaleFetchers(n int) Option {
	return func(c *Cache) {
		c.staleFetchers = n
	}
}

// WithStaleQueueSize sets the size of the chan for queuing items.
// Default: 1000
func WithStaleQueueSize(n int) Option {
	return func(c *Cache) {
		c.staleQueueSize = n
	}
}

// WithFetchers sets the max number of concurrent fetches.
// Default: 100
func WithFetchers(n int) Option {
	return func(c *Cache) {
		c.maxFetchers = n
	}
}

// WithStale allows to decide globally if expired items are served.
// If false, WithStaleValidator will be ignored.
// Default: true
func WithStale(useStale bool) Option {
	return func(c *Cache) {
		c.canUseStale = useStale
	}
}

// WithStaleValidator allows to use a function to decide if a stale item
// will be served. The duration is the extra time after expiration.
// Default: nil (WithStale() will decide if stale items are served)
func WithStaleValidator(f func(interface{}, time.Duration) bool) Option {
	return func(c *Cache) {
		c.staleValidator = f
	}
}

// New builds a cache given some options.
func New(opts ...Option) (*Cache, error) {
	c := &Cache{
		posSize:        5000,
		posTTL:         60 * time.Second,
		posPruneSize:   0,
		negSize:        500,
		negTTL:         5 * time.Second,
		negPruneSize:   0,
		staleFetchers:  3,
		staleQueueSize: 1000,
		maxFetchers:    100,
		canUseStale:    true,
	}

	for _, o := range opts {
		o(c)
	}

	c.fetchLimiter = make(chan struct{}, c.maxFetchers)

	if c.posPruneSize == 0 {
		c.posPruneSize = c.posSize/20 + 1
	}
	if c.negPruneSize == 0 {
		c.negPruneSize = c.negSize/20 + 1
	}

	c.posCache = ccache.New(ccache.Configure().
		MaxSize(int64(c.posSize)).ItemsToPrune(uint32(c.posPruneSize)))
	c.negCache = ccache.New(ccache.Configure().
		MaxSize(int64(c.posSize)).ItemsToPrune(uint32(c.posPruneSize)))

	// for cache locking
	c.fetching = make(map[string]struct{})
	c.fetchCond = sync.NewCond(&c.fetchLock)

	// for async stale fetch
	c.fetchQueue = make(chan fetchReq, c.staleQueueSize)
	c.queued = make(map[string]struct{})

	for i := 0; i < c.staleFetchers; i++ {
		go c.staleFetcher()
	}
	return c, nil
}

// Fetch returns an object given its cache key and a Fetcher function for fetching it if
// it expired or missing.
//
// This function will usually be implemented as a closure in order to capture
// the various parameters needed for fetching from the backend the entry
// corresponding to the cache key.
//
// Entries are first looked up in the positive cache, then negative, then fetched.
//
// An asynchronous fetch will happen if the entry is stale.
func (c *Cache) Fetch(key string, f Fetcher) (interface{}, error) {
	atomic.AddUint64(&c.requests, 1)
	item, cached, err := c.tryCache(key, f)
	if cached {
		atomic.AddUint64(&c.hits, 1)
		// fresh or stale
		return item, err
	}

	// entry not in cache
	c.fetchLock.Lock()
	_, fetching := c.fetching[key]
	if !fetching {
		// nobody is fetching it yet, let's do it
		c.fetching[key] = struct{}{}
		c.fetchLock.Unlock()
		c.fetchLimiter <- struct{}{}
		atomic.AddUint64(&c.newFetches, 1)
		item, err = c.cacheItem(key, f)
		<-c.fetchLimiter
		c.endFetch(key)
		c.fetchCond.Broadcast()
		return item, err
	}
	// wait for the fetcher to finish
	for {
		c.fetchCond.Wait()
		_, ok := c.fetching[key]
		if !ok {
			break
		}
	}
	c.fetchLock.Unlock()

	// get the hopefully newly cached entry
	item, cached, err = c.tryCache(key, f)
	if cached {
		return item, err
	}
	// last resort (if too small a cache)
	c.fetchLimiter <- struct{}{}
	item, _, err = f()
	<-c.fetchLimiter
	return item, err
}

// tryCache tries to find the given key in the positive and negative caches.
// If an element is expired, it will be queued for async fetch and its stale
// version will be returned immediately.
// The boolean in the return value indicates if the key has been found in cache.
func (c *Cache) tryCache(key string, f Fetcher) (interface{}, bool, error) {
	item := c.posCache.Get(key)
	if item != nil {
		valid := true
		if item.Expired() {
			// stale item, let's enqueue a refresh
			c.enqueueFetch(key, f)
			if !c.useStale(item) {
				valid = false
			}
		}
		if valid { // not expired or can use stale
			return item.Value(), true, nil
		}
		// if cannot use stale
		return nil, false, nil
	}

	item = c.negCache.Get(key)
	if item != nil {
		if item.Expired() {
			// stale negative, remove it from cache
			c.negCache.Delete(key)
		}
		ne := item.Value().(*negCacheEntry)
		return ne.x, true, ne.err
	}
	return nil, false, nil
}

// enqueueFetch puts a fetch request in the queue.
// It does nothing if the request is already in the queue, or if the queue is full.
func (c *Cache) enqueueFetch(key string, f Fetcher) {
	c.queuedLock.Lock()
	_, ok := c.queued[key]
	if !ok {
		select {
		case c.fetchQueue <- fetchReq{key, f}:
			c.queued[key] = struct{}{}
		default:
			// drop request on full queue instead of blocking
		}
	}
	c.queuedLock.Unlock()
}

// staleFetcher grabs a fetch request from the chan and executes it.
// Requests are guaranteed to be unique in the queue (using "queued" map): no need to protect
// against multiple concurrent fetches, no need to use cache locking.
func (c *Cache) staleFetcher() {
	for fr := range c.fetchQueue {
		// fetch it
		atomic.AddUint64(&c.staleFetches, 1)
		_, _ = c.cacheItem(fr.key, fr.f)
		c.endQueuing(fr.key)
	}
}

// cacheItem fetches an object using the supplied closure and stores it
// in cache using the supplied key. Depending on the error and validity returned by
// the closure, the entry will end in either the positive or the negative cache.
// if error not nil, store in negative cache (but keep positive entry);
// else if validity is false, store in negative cache and delete positive entry;
// else (error nil and validity true) store in positive cache and remove neg entry
func (c *Cache) cacheItem(key string, f Fetcher) (interface{}, error) {
	item, valid, err := f()

	if err != nil {
		c.negCache.Set(key, &negCacheEntry{item, err}, c.negTTL)
	} else if !valid {
		c.negCache.Set(key, &negCacheEntry{item, err}, c.negTTL)
		c.posCache.Delete(key)
	} else {
		c.posCache.Set(key, item, c.posTTL)
		c.negCache.Delete(key)
	}
	return item, err
}

// endQueuing marks an item as not being in the fetch queue anymore.
func (c *Cache) endQueuing(key string) {
	c.queuedLock.Lock()
	delete(c.queued, key)
	c.queuedLock.Unlock()
}

// endFetch marks an item as not being fetched anymore.
func (c *Cache) endFetch(key string) {
	c.fetchLock.Lock()
	delete(c.fetching, key)
	c.fetchLock.Unlock()
}

// useStale decides if we serve a stale item. The behaviour can be modified globally
// with the option NoStale(), or per-item using a supplied function(interface{}, time.Duration)
// where the duration is the extra time after cache expiration), with option StaleValidator()
func (c *Cache) useStale(item *ccache.Item) bool {
	if !c.canUseStale {
		return false
	}
	if c.staleValidator != nil {
		// ccache gives a negative TTL for expired items, inverse it
		return c.staleValidator(item.Value(), -item.TTL())
	}
	return true
}

// Hits returns the number of cache hits since start.
func (c *Cache) Hits() uint64 {
	return atomic.LoadUint64(&c.hits)
}

// Requests returns the number of cache requests (hits and misses) since start.
func (c *Cache) Requests() uint64 {
	return atomic.LoadUint64(&c.requests)
}

// NewFetches returns the number of fetches for items not in cache, since start.
func (c *Cache) NewFetches() uint64 {
	return atomic.LoadUint64(&c.newFetches)
}

// StaleFetches returns the number of fetches for expired items, since start.
func (c *Cache) StaleFetches() uint64 {
	return atomic.LoadUint64(&c.staleFetches)
}
