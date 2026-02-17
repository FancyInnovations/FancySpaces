package spaces

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/dgraph-io/ristretto/v2"
)

var (
	sizeOfUser = int64(reflect.TypeFor[InternalSpace]().Size())
	ttl        = 30 * time.Minute
)

type spacesCache struct {
	cache *ristretto.Cache[string, *InternalSpace]
}

func newSpacesCache() *spacesCache {
	cache, err := ristretto.NewCache(&ristretto.Config[string, *InternalSpace]{
		NumCounters: 100 * 10,         // x10 of expected number of elements when full
		MaxCost:     64 * 1024 * 1024, // 64 MB
		BufferItems: 64,               // keep 64
	})
	if err != nil {
		slog.Error("Failed to create space cache", sloki.WrapError(err))
		panic(err)
	}

	return &spacesCache{
		cache: cache,
	}
}

func (c *spacesCache) GetByID(id string) (*InternalSpace, error) {
	if id == "" {
		return nil, ErrSpaceNotFound
	}

	user, found := c.cache.Get(id)
	if !found {
		return nil, ErrSpaceNotFound
	}

	return user, nil
}

func (c *spacesCache) UpsertSpace(s *InternalSpace) {
	c.cache.SetWithTTL(s.ID, s, sizeOfUser, ttl)
}

func (c *spacesCache) Invalidate(s *InternalSpace) {
	c.cache.Del(s.ID)
}
