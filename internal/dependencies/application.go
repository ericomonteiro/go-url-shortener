package dependencies

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
)

type ShortenerApp struct {
	DB    *sql.DB
	Redis *redis.Client
}
