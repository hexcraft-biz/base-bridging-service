package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-bridging-service/config"
	"github.com/hexcraft-biz/base-bridging-service/controllers"
	"github.com/hexcraft-biz/feature"
)

func LoadCommon(e *gin.Engine, cfg config.ConfigInterface) {
	c := controllers.NewCommon(cfg)
	e.NoRoute(c.NotFound())

	commonV1 := feature.New(e, "/healthcheck/v1")
	commonV1.GET("/ping", c.Ping())
}
