package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-bridging-service/config"
	"github.com/hexcraft-biz/base-bridging-service/controllers"
	"github.com/hexcraft-biz/feature"
)

func LoadBridging(e *gin.Engine, cfg config.ConfigInterface) {

	bridgingV1 := feature.New(e, "/bridging/v1")

	ec := controllers.NewEndpoints(cfg)
	tc := controllers.NewTopics(cfg)
	etrc := controllers.NewEndpointTopicRels(cfg)

	bridgingV1.GET("/endpoints", ec.List())
	bridgingV1.GET("/endpoints/:id", ec.GetOne())
	bridgingV1.GET("/endpoints/:id/topics", ec.GetTopics())
	bridgingV1.POST("/endpoints", ec.Create())
	bridgingV1.DELETE("/endpoints/:id", ec.Delete())

	bridgingV1.GET("/topics", tc.List())
	bridgingV1.GET("/topics/:id", tc.GetOne())
	bridgingV1.POST("/topics", tc.Create())
	bridgingV1.DELETE("/topics/:id", tc.Delete())

	bridgingV1.GET("/endpoint-topic-rels/:id", etrc.GetOne())
	bridgingV1.POST("/endpoint-topic-rels", etrc.Create())
	bridgingV1.DELETE("/endpoint-topic-rels/:id", etrc.Delete())
}
