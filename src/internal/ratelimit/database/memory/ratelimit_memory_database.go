package memory

import (
	"reflect"
	"time"

	"github.com/dgraph-io/ristretto/v2"
)

var (
	ttl        = 15 * time.Minute
	sizeOfTime = int64(reflect.TypeOf(time.Time{}).Size())
)

type DB struct {
	tokens      *ristretto.Cache[string, int]
	refillTimes *ristretto.Cache[string, time.Time]
}

func NewDB() *DB {
	tokens, err := ristretto.NewCache(&ristretto.Config[string, int]{
		NumCounters: 500 * 10,         // x10 of expected number of elements when full
		MaxCost:     16 * 1024 * 1024, // 16 MB
		BufferItems: 64,               // keep 64
	})
	if err != nil {
		panic(err)
	}

	refillTimes, err := ristretto.NewCache(&ristretto.Config[string, time.Time]{
		NumCounters: 500 * 10,         // x10 of expected number of elements when full
		MaxCost:     16 * 1024 * 1024, // 16 MB
		BufferItems: 64,               // keep 64
	})
	if err != nil {
		panic(err)
	}

	return &DB{
		tokens:      tokens,
		refillTimes: refillTimes,
	}
}

func (db *DB) GetTokens(client string) (int, error) {
	count, exists := db.tokens.Get(client)
	if !exists {
		return 0, nil
	}

	return count, nil
}

func (db *DB) SetTokens(client string, tokens int) error {
	db.tokens.SetWithTTL(client, tokens, 4, ttl)
	return nil
}

func (db *DB) GetLastRefill(client string) (time.Time, error) {
	t, exists := db.refillTimes.Get(client)
	if !exists {
		return time.Now().Add((-5) * time.Minute), nil // return a time in the past to trigger immediate refill
	}

	return t, nil
}

func (db *DB) SetLastRefill(client string, t time.Time) error {
	db.refillTimes.SetWithTTL(client, t, sizeOfTime, ttl)
	return nil
}
