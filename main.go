package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/hexcraft-biz/base-bridging-service/handlers"
	"github.com/hexcraft-biz/base-bridging-service/service"
	"github.com/hexcraft-biz/env"
	envRedis "github.com/hexcraft-biz/env/redis"
	"github.com/hexcraft-biz/feature"
	"github.com/jmoiron/sqlx"
)

type TestStruct struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
}

func main() {
	// Prepare your config implemented from ConfigInterface.
	cfg, _ := Load()
	cfg.DBOpen(false)
	cfg.InitRedis()

	// New base-bridging-service.
	engine := service.New(cfg)

	// Example for set up route and pubsub handler.
	testV1 := feature.New(engine, "/test/v1")
	testV1.GET("/ping", func(c *gin.Context) {
		c.Set("publishData", TestStruct{
			Count: 1,
			Name:  "ok",
		})
		c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
	}, handlers.GcpPubsubPublish(cfg))

	// Then Run Gin Engine.
	engine.Run(":" + cfg.Env.AppPort)
}

//================================================================
// Env implemented from env pkg
//================================================================
type Env struct {
	*env.Prototype
	GcpProjectID string
}

func FetchEnv() (*Env, error) {
	if e, err := env.Fetch(); err != nil {
		return nil, err
	} else {
		return &Env{
			Prototype:    e,
			GcpProjectID: os.Getenv("GCP_PROJECT_ID"),
		}, nil
	}
}

//================================================================
// Config implement ConfigInterface
//================================================================
type Config struct {
	*Env
	DB    *sqlx.DB
	Redis *redis.Client
}

func Load() (*Config, error) {
	e, err := FetchEnv()
	if err != nil {
		return nil, err
	}

	return &Config{Env: e}, nil
}

func (cfg *Config) InitRedis() error {
	var err error

	cfg.Redis, err = envRedis.NewRedisClient()
	return err
}

func (cfg *Config) DBOpen(init bool) error {
	var err error

	cfg.DBClose()
	cfg.DB, err = cfg.MysqlConnectWithMode(init)

	return err
}

func (cfg *Config) DBClose() {
	if cfg.DB != nil {
		cfg.DB.Close()
	}
}

//================================================================
// ConfigInterface Functions
//================================================================
func (cfg *Config) GetDB() *sqlx.DB {
	return cfg.DB
}

func (cfg *Config) GetRedis() *redis.Client {
	return cfg.Redis
}

func (cfg *Config) GetTrustProxy() string {
	return cfg.Env.TrustProxy
}

func (cfg *Config) GetGcpProjectID() string {
	return cfg.Env.GcpProjectID
}
