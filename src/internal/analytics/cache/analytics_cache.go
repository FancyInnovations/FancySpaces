package cache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/fancyinnovations/fancyspaces/internal/analytics"
)

var (
	ttl = 10 * time.Minute
)

type Cache struct {
	versionDownloadCounts    *ristretto.Cache[string, uint64]
	allVersionDownloadCounts *ristretto.Cache[string, map[string]uint64]
}

func NewCache() *Cache {
	versionDownloadCounts, err := ristretto.NewCache(&ristretto.Config[string, uint64]{
		NumCounters: 100 * 10,         // x10 of expected number of elements when full
		MaxCost:     16 * 1024 * 1024, // 16 MB
		BufferItems: 64,               // keep 64
	})
	if err != nil {
		panic(err)
	}

	allVersionDownloadCounts, err := ristretto.NewCache(&ristretto.Config[string, map[string]uint64]{
		NumCounters: 100 * 10,         // x10 of expected number of elements when full
		MaxCost:     16 * 1024 * 1024, // 16 MB
		BufferItems: 64,               // keep 64
	})
	if err != nil {
		panic(err)
	}

	return &Cache{
		versionDownloadCounts:    versionDownloadCounts,
		allVersionDownloadCounts: allVersionDownloadCounts,
	}
}

func (c *Cache) GetDownloadCountForVersion(spaceID, versionID string) (error, uint64) {
	key := spaceID + ":" + versionID

	count, found := c.versionDownloadCounts.Get(key)
	if !found {
		return analytics.ErrNotInCache, 0
	}

	return nil, count
}

func (c *Cache) GetDownloadCountForVersions(spaceID string) (error, map[string]uint64) {
	counts, found := c.allVersionDownloadCounts.Get(spaceID)
	if !found {
		return analytics.ErrNotInCache, nil
	}

	return nil, counts
}

func (c *Cache) SetDownloadCountForVersion(spaceID, versionID string, count uint64) {
	key := spaceID + ":" + versionID

	c.versionDownloadCounts.SetWithTTL(key, count, 4, ttl)
}

func (c *Cache) SetDownloadCountForVersions(spaceID string, counts map[string]uint64) {
	c.allVersionDownloadCounts.SetWithTTL(spaceID, counts, int64(len(counts)*4), ttl)
}
