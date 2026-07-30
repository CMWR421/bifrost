package lrucache

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fillWith returns a Loader serving a fixed value and index key.
func fillWith(value string, indexKey string) Loader[string] {
	return func() (string, string, error) {
		return value, indexKey, nil
	}
}

func TestCache_HitAndMiss(t *testing.T) {
	c := New[string](4)
	_, ok := c.Get("k1")
	assert.False(t, ok)

	v, err := c.Fill("k1", fillWith("v1", "id1"))
	require.NoError(t, err)
	assert.Equal(t, "v1", v)

	got, ok := c.Get("k1")
	assert.True(t, ok)
	assert.Equal(t, "v1", got)
	assert.Equal(t, 1, c.Len())
}

func TestCache_CapacityEvictionAndIndexCleanup(t *testing.T) {
	c := New[string](2)
	for i := 1; i <= 3; i++ {
		_, err := c.Fill(fmt.Sprintf("k%d", i), fillWith(fmt.Sprintf("v%d", i), fmt.Sprintf("id%d", i)))
		require.NoError(t, err)
	}
	assert.Equal(t, 2, c.Len())

	// k1 was least recently used and must be gone, along with its index
	// entry: evicting by its old index key must not disturb survivors.
	_, ok := c.Get("k1")
	assert.False(t, ok)
	c.EvictByIndex("id1")
	assert.Equal(t, 2, c.Len())

	_, ok = c.Get("k2")
	assert.True(t, ok)
	_, ok = c.Get("k3")
	assert.True(t, ok)
}

func TestCache_ValidatorRejectionIsMissAndRemoval(t *testing.T) {
	type timedValue struct {
		value     string
		expiresAt time.Time
	}
	c := New(4, WithValidator(func(v timedValue) bool {
		return time.Now().Before(v.expiresAt)
	}))
	_, err := c.Fill("k1", func() (timedValue, string, error) {
		return timedValue{value: "v1", expiresAt: time.Now().Add(-time.Minute)}, "id1", nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, c.Len())

	_, ok := c.Get("k1")
	assert.False(t, ok, "rejected entry must be reported as a miss")
	assert.Equal(t, 0, c.Len(), "rejected entry must be removed")
}

func TestCache_EvictExactKey(t *testing.T) {
	c := New[string](4)
	_, err := c.Fill("k1", fillWith("v1", "id1"))
	require.NoError(t, err)

	c.Evict("k1")
	_, ok := c.Get("k1")
	assert.False(t, ok)
	assert.Equal(t, 0, c.Len())

	// Evicting an absent key is a no-op, not a panic.
	c.Evict("never-set")
}

func TestCache_EvictByIndex(t *testing.T) {
	c := New[string](4)
	_, err := c.Fill("k1", fillWith("v1", "id1"))
	require.NoError(t, err)
	_, err = c.Fill("k2", fillWith("v2", "id2"))
	require.NoError(t, err)

	c.EvictByIndex("id1")
	_, ok := c.Get("k1")
	assert.False(t, ok, "entry holding the evicted index key must be gone")
	_, ok = c.Get("k2")
	assert.True(t, ok, "unrelated entry must survive")

	// Unknown and empty index keys are no-ops.
	c.EvictByIndex("unknown")
	c.EvictByIndex("")
	assert.Equal(t, 1, c.Len())
}

func TestCache_EvictWhere(t *testing.T) {
	c := New[string](8)
	for _, k := range []string{"a:1", "a:2", "b:1"} {
		_, err := c.Fill(k, fillWith("v", ""))
		require.NoError(t, err)
	}

	c.EvictWhere(func(key string) bool { return key[0] == 'a' })
	_, ok := c.Get("a:1")
	assert.False(t, ok)
	_, ok = c.Get("a:2")
	assert.False(t, ok)
	_, ok = c.Get("b:1")
	assert.True(t, ok)

	// A nil predicate is a no-op, not a panic.
	c.EvictWhere(nil)
	assert.Equal(t, 1, c.Len())
}

func TestCache_Flush(t *testing.T) {
	c := New[string](4)
	for i := 1; i <= 3; i++ {
		_, err := c.Fill(fmt.Sprintf("k%d", i), fillWith("v", fmt.Sprintf("id%d", i)))
		require.NoError(t, err)
	}

	c.Flush()
	assert.Equal(t, 0, c.Len())
	for i := 1; i <= 3; i++ {
		_, ok := c.Get(fmt.Sprintf("k%d", i))
		assert.False(t, ok)
	}
}

func TestCache_InflightDedup(t *testing.T) {
	c := New[string](4)
	var calls atomic.Int32
	release := make(chan struct{})
	loader := func() (string, string, error) {
		calls.Add(1)
		<-release
		return "v1", "id1", nil
	}

	var wg sync.WaitGroup
	results := make([]string, 2)
	for i := range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := c.Fill("k1", loader)
			// assert, not require: require.NoError calls t.FailNow(), which
			// per the testing package's own contract must only be called
			// from the goroutine running the test.
			assert.NoError(t, err)
			results[i] = v
		}()
	}
	// Let both goroutines reach the fill before releasing the leader.
	assert.Eventually(t, func() bool { return calls.Load() == 1 }, time.Second, time.Millisecond)
	close(release)
	wg.Wait()

	assert.Equal(t, int32(1), calls.Load(), "loader must run exactly once")
	assert.Equal(t, []string{"v1", "v1"}, results)
}

func TestCache_ErrorSharedNotCached(t *testing.T) {
	c := New[string](4)
	sentinel := errors.New("load failed")
	var calls atomic.Int32
	failing := func() (string, string, error) {
		calls.Add(1)
		return "", "", sentinel
	}

	_, err := c.Fill("k1", failing)
	assert.ErrorIs(t, err, sentinel)
	assert.Equal(t, 0, c.Len(), "errors must not be cached")

	// A second call retries the loader instead of serving the error.
	_, err = c.Fill("k1", failing)
	assert.ErrorIs(t, err, sentinel)
	assert.Equal(t, int32(2), calls.Load())
}

func TestCache_GenerationGuardDiscardsStaleFill(t *testing.T) {
	c := New[string](4)
	started := make(chan struct{})
	release := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		v, err := c.Fill("k1", func() (string, string, error) {
			close(started)
			<-release
			return "stale", "id1", nil
		})
		// The caller still receives the loaded value; only the install is
		// discarded. assert, not require: require.NoError calls
		// t.FailNow(), which per the testing package's own contract must
		// only be called from the goroutine running the test.
		assert.NoError(t, err)
		assert.Equal(t, "stale", v)
	}()

	<-started
	c.Evict("k1") // lands while the load is in flight
	close(release)
	wg.Wait()

	_, ok := c.Get("k1")
	assert.False(t, ok, "a fill racing an eviction must not install its result")
}

func TestCache_UpsertRebindsIndex(t *testing.T) {
	c := New[string](4)
	_, err := c.Fill("k1", fillWith("v1", "id-old"))
	require.NoError(t, err)

	// Re-fill the same key with a new index key after evicting: the old
	// index mapping must not linger.
	c.Evict("k1")
	_, err = c.Fill("k1", fillWith("v2", "id-new"))
	require.NoError(t, err)

	c.EvictByIndex("id-old")
	got, ok := c.Get("k1")
	assert.True(t, ok, "stale index key must not evict the rebound entry")
	assert.Equal(t, "v2", got)

	c.EvictByIndex("id-new")
	_, ok = c.Get("k1")
	assert.False(t, ok)
}

func TestCache_LoaderPanicReleasesWaitersWithError(t *testing.T) {
	c := New[string](4)
	started := make(chan struct{})

	var waiterErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-started
		_, waiterErr = c.Fill("k1", fillWith("never-used", ""))
	}()

	func() {
		defer func() {
			r := recover()
			require.NotNil(t, r, "the leader must re-panic")
		}()
		_, _ = c.Fill("k1", func() (string, string, error) {
			close(started)
			// Give the waiter a moment to join the inflight call; if it
			// races past cleanup instead, it just runs its own fresh fill,
			// which this test also accepts below.
			time.Sleep(20 * time.Millisecond)
			panic("boom")
		})
	}()
	wg.Wait()

	if waiterErr != nil {
		assert.Contains(t, waiterErr.Error(), "panicked")
	}
	// Either way the key must not be wedged: a fresh fill succeeds.
	v, err := c.Fill("k1", fillWith("v-after", ""))
	require.NoError(t, err)
	assert.Equal(t, "v-after", v)
}

func TestCache_NilSafety(t *testing.T) {
	var c *Cache[string]
	_, ok := c.Get("k1")
	assert.False(t, ok)
	c.Evict("k1")
	c.EvictByIndex("id1")
	c.EvictWhere(func(string) bool { return true })
	c.Flush()
	assert.Equal(t, 0, c.Len())

	v, err := c.Fill("k1", fillWith("v1", "id1"))
	require.NoError(t, err)
	assert.Equal(t, "v1", v, "nil cache must degrade to calling the loader directly")
}

func TestNew_PanicsOnNonPositiveCapacity(t *testing.T) {
	assert.Panics(t, func() { New[string](0) })
	assert.Panics(t, func() { New[string](-1) })
}

func TestEncodeDecodeKey_RoundTrip(t *testing.T) {
	parts := []string{"user", "alice", "client-1"}
	key := EncodeKey(parts...)
	got, ok := DecodeKey(key, len(parts))
	require.True(t, ok)
	assert.Equal(t, parts, got)
}

// TestEncodeDecodeKey_NoCollisionOnEmbeddedDelimiter pins the reason this
// codec exists over a plain separator join: a value that happens to
// contain the separator (or, here, digits and a colon shaped like a length
// prefix) must not let one tuple's key collide with a different tuple's.
func TestEncodeDecodeKey_NoCollisionOnEmbeddedDelimiter(t *testing.T) {
	keyA := EncodeKey("user", "a\x00b", "c")
	keyB := EncodeKey("user", "a", "b\x00c")
	assert.NotEqual(t, keyA, keyB, "distinct tuples must not alias to the same key")

	gotA, ok := DecodeKey(keyA, 3)
	require.True(t, ok)
	assert.Equal(t, []string{"user", "a\x00b", "c"}, gotA)

	gotB, ok := DecodeKey(keyB, 3)
	require.True(t, ok)
	assert.Equal(t, []string{"user", "a", "b\x00c"}, gotB)
}

func TestDecodeKey_RejectsMalformedInput(t *testing.T) {
	tests := []struct {
		name string
		key  string
		n    int
	}{
		{"empty string", "", 3},
		{"no colon", "abc", 1},
		{"non-numeric length", "x:abc", 1},
		{"negative length", "-1:a", 1},
		{"length exceeds remaining bytes", "10:ab", 1},
		{"trailing garbage after all parts", "1:a1:btrailing", 2},
		{"too few parts", "4:user", 2},
		{"negative part count", "1:a", -1},
		{"non-canonical length prefix with leading plus", "+1:a", 1},
		{"non-canonical length prefix with leading zero", "01:a", 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, ok := DecodeKey(tc.key, tc.n)
			assert.False(t, ok)
		})
	}
}
