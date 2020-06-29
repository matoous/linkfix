package callonce_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"

	"github.com/matoous/linkfix/internal/callonce"
)

func TestReturnsSameValue(t *testing.T) {
	var g callonce.Group

	val1, err := g.Do(context.Background(), "key1", func(ctx context.Context) (interface{}, error) {
		return "value1", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "value1", val1)

	var called int64

	val2, err := g.Do(context.Background(), "key1", func(ctx context.Context) (interface{}, error) {
		atomic.AddInt64(&called, 1)
		return "value2", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "value1", val2, "subsequent call returns the same value")
	assert.Equal(t, int64(0), atomic.LoadInt64(&called))
}

func TestDeduplicatesRequests(t *testing.T) {
	var g callonce.Group
	var eg errgroup.Group

	var called int64

	eg.Go(func() error {
		val, err := g.Do(context.Background(), "key1", func(ctx context.Context) (interface{}, error) {
			atomic.AddInt64(&called, 1)
			// Unfortunately there is no way to synchronize this function with the second one
			// since only one of them gets called, so we need to resort to Sleep().
			time.Sleep(10 * time.Millisecond)
			return "value1", nil
		})
		assert.NoError(t, err)
		assert.Equal(t, "value1", val)
		assert.Equal(t, int64(1), atomic.LoadInt64(&called))
		return nil
	})

	eg.Go(func() error {
		val, err := g.Do(context.Background(), "key1", func(ctx context.Context) (interface{}, error) {
			atomic.AddInt64(&called, 1)
			// Unfortunately there is no way to synchronize this function with the second one
			// since only one of them gets called, so we need to resort to Sleep().
			time.Sleep(10 * time.Millisecond)
			return "value1", nil
		})
		assert.NoError(t, err)
		assert.Equal(t, "value1", val, "subsequent call returns the same value")
		assert.Equal(t, int64(1), atomic.LoadInt64(&called))
		return nil
	})

	_ = eg.Wait()
}

func TestDifferentKeys(t *testing.T) {
	var g callonce.Group
	var eg errgroup.Group

	c := make(chan struct{})

	eg.Go(func() error {
		val1, err := g.Do(context.Background(), "key1", func(ctx context.Context) (interface{}, error) {
			<-c
			return "value1", nil
		})
		assert.NoError(t, err)
		assert.Equal(t, "value1", val1)
		return nil
	})

	eg.Go(func() error {
		val2, err := g.Do(context.Background(), "key2", func(ctx context.Context) (interface{}, error) {
			close(c)
			return "value2", nil
		})
		assert.NoError(t, err)
		assert.Equal(t, "value2", val2, "subsequent call returns the same value")
		return nil
	})

	_ = eg.Wait()
}
