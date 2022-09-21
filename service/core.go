package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-bridging-service/config"
	"github.com/hexcraft-biz/base-bridging-service/features"
)

func New(cfg config.ConfigInterface) *gin.Engine {

	engine := gin.Default()
	engine.SetTrustedProxies([]string{cfg.GetTrustProxy()})

	// base features
	features.LoadCommon(engine, cfg)
	// Bridging features
	features.LoadBridging(engine, cfg)

	return engine
}
