package initializers

import (
	"go_auth/config"

	"github.com/dgraph-io/ristretto"
	"github.com/gin-contrib/cache/persistence"
)

var Store *persistence.InMemoryStore
var Cache *ristretto.Cache

func InitCache() {
	var err error

	Store = persistence.NewInMemoryStore(config.CacheTTL)

	Cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: config.CacheSize,
		MaxCost:     config.MaxCost,
		BufferItems: config.BufferItems,
	})
	if err != nil {
		panic(err)
	}
}
