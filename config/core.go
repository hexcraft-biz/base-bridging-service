package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type ConfigInterface interface {
	GetDB() *sqlx.DB
	GetRedis() *redis.Client
	GetTrustProxy() string
	GetGcpProjectID() string
}
