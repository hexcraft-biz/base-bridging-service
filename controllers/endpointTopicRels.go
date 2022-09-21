package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/hexcraft-biz/base-bridging-service/config"
	"github.com/hexcraft-biz/base-bridging-service/models"
	"github.com/hexcraft-biz/controller"
)

type EndpointTopicRels struct {
	*controller.Prototype
	Config config.ConfigInterface
}

func NewEndpointTopicRels(cfg config.ConfigInterface) *EndpointTopicRels {
	return &EndpointTopicRels{
		Prototype: controller.New("endpointTopicRels", cfg.GetDB()),
		Config:    cfg,
	}
}

func (ctrl *EndpointTopicRels) NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
	}
}

type TargetEndpointTopicRel struct {
	ID string `uri:"id" binding:"required"`
}

func (ctrl *EndpointTopicRels) GetOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		var targetETR TargetEndpointTopicRel
		if err := c.ShouldBindUri(&targetETR); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewEndpointTopicRelsTableEngine(ctrl.DB).GetByID(targetETR.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusOK, entityRes)
				return
			}
		}
	}
}

type createEndpointTopicRelParams struct {
	EndpointId string `json:"endpointId" binding:"required"`
	TopicId    string `json:"topicId" binding:"required"`
}

func (ctrl *EndpointTopicRels) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params createEndpointTopicRelParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewEndpointTopicRelsTableEngine(ctrl.DB).Insert(params.EndpointId, params.TopicId); err != nil {
			if myErr, ok := err.(*mysql.MySQLError); ok && myErr.Number == 1062 {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": http.StatusText(http.StatusConflict)})
				return
			} else if ok && myErr.Number == 1452 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		} else {
			if entityRes, err := models.NewEndpointsTableEngine(ctrl.DB).GetByID(params.EndpointId); err == nil {
				ctx := context.Background()
				ctrl.Config.GetRedis().Del(ctx, entityRes.Path).Result()
			}

			c.AbortWithStatusJSON(http.StatusCreated, entityRes)
			return
		}
	}
}

func (ctrl *EndpointTopicRels) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, dbEngine := new(TargetEndpointTopicRel), models.NewEndpointTopicRelsTableEngine(ctrl.DB)

		if err := c.ShouldBindUri(req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
			return
		}

		if rel, err := dbEngine.GetByID(req.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else if rel == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
			return
		} else if _, err := dbEngine.DeleteByID(req.ID); err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1451 {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": err.Error()})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		} else {
			if entityRes, err := models.NewEndpointsTableEngine(ctrl.DB).GetByID(rel.EndpointID.String()); err == nil {
				ctx := context.Background()
				ctrl.Config.GetRedis().Del(ctx, entityRes.Path).Result()
			}

			c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": http.StatusText(http.StatusNoContent)})
			return
		}
	}
}
