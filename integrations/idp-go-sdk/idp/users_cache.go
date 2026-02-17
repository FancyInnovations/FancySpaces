package idp

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/dgraph-io/ristretto/v2"
)

var (
	sizeOfUser = int64(reflect.TypeFor[User]().Size())
	ttl        = 30 * time.Minute
)

type usersCache struct {
	idCache    *ristretto.Cache[string, *User]
	emailCache *ristretto.Cache[string, *User]
}

func newUsersCache() *usersCache {
	idCache, err := ristretto.NewCache(&ristretto.Config[string, *User]{
		NumCounters: 100 * 10,         // x10 of expected number of elements when full
		MaxCost:     64 * 1024 * 1024, // 64 MB
		BufferItems: 64,               // keep 64
	})
	if err != nil {
		slog.Error("Failed to create ID cache", sloki.WrapError(err))
		panic(err)
	}

	emailCache, err := ristretto.NewCache(&ristretto.Config[string, *User]{
		NumCounters: 100 * 10,         // x10 of expected number of elements when full
		MaxCost:     64 * 1024 * 1024, // 64 MB
		BufferItems: 64,               // keep 64
	})
	if err != nil {
		slog.Error("Failed to create email cache", sloki.WrapError(err))
		panic(err)
	}

	return &usersCache{
		idCache:    idCache,
		emailCache: emailCache,
	}
}

func (c *usersCache) GetByID(id string) (*User, error) {
	if id == "" {
		return nil, ErrUserNotFound
	}

	user, found := c.idCache.Get(id)
	if !found {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (c *usersCache) GetByEmail(email string) (*User, error) {
	if email == "" {
		return nil, ErrUserNotFound
	}

	user, found := c.emailCache.Get(email)
	if !found {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (c *usersCache) UpsertUser(u *User) {
	c.idCache.SetWithTTL(u.ID, u, sizeOfUser, ttl)
	c.emailCache.SetWithTTL(u.Email, u, sizeOfUser, ttl)
}

func (c *usersCache) Invalidate(u *User) {
	c.idCache.Del(u.ID)
	c.emailCache.Del(u.Email)
}
