package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/hexcraft-biz/base-bridging-service/config"
	"github.com/hexcraft-biz/base-bridging-service/models"
	"github.com/hexcraft-biz/controller"
	"github.com/hexcraft-biz/model"
)

type Endpoints struct {
	*controller.Prototype
	Config config.ConfigInterface
}

func NewEndpoints(cfg config.ConfigInterface) *Endpoints {
	return &Endpoints{
		Prototype: controller.New("endpoints", cfg.GetDB()),
		Config:    cfg,
	}
}

func (ctrl *Endpoints) NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
	}
}

type ListEndpointsParams struct {
	Limit  uint64 `form:"limit,default=20"`
	Offset uint64 `form:"offset,default=0"`
}

func (ctrl *Endpoints) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		listParams := new(ListEndpointsParams)
		if err := c.ShouldBindQuery(listParams); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
			return
		}
		pg := model.NewPagination(listParams.Offset, listParams.Limit)

		if endpoints, err := models.NewEndpointsTableEngine(ctrl.DB).List(pg); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusOK, endpoints)
			return
		}
	}
}

type TargetEndpoint struct {
	ID string `uri:"id" binding:"required"`
}

func (ctrl *Endpoints) GetOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		var targetEndpoint TargetEndpoint
		if err := c.ShouldBindUri(&targetEndpoint); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewEndpointsTableEngine(ctrl.DB).GetByID(targetEndpoint.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				if absRes, absErr := entityRes.GetAbsEndpoint(); absErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": absErr.Error()})
					return
				} else {
					c.AbortWithStatusJSON(http.StatusOK, absRes)
					return
				}
			}
		}
	}
}

func (ctrl *Endpoints) GetTopics() gin.HandlerFunc {
	return func(c *gin.Context) {
		var targetEndpoint TargetEndpoint
		if err := c.ShouldBindUri(&targetEndpoint); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewEndpointTopicRelsTableEngine(ctrl.DB).GetByEndpointID(targetEndpoint.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusOK, entityRes)
			return
		}
	}
}

type createEndpointParams struct {
	Path string `json:"path" binding:"required,min=5,max=128"`
}

func (ctrl *Endpoints) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params createEndpointParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewEndpointsTableEngine(ctrl.DB).Insert(params.Path); err != nil {
			if myErr, ok := err.(*mysql.MySQLError); ok && myErr.Number == 1062 {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": http.StatusText(http.StatusConflict)})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		} else {
			if absRes, absErr := entityRes.GetAbsEndpoint(); absErr != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": absErr.Error()})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusCreated, absRes)
				return
			}
		}
	}
}

func (ctrl *Endpoints) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, dbEngine := new(TargetEndpoint), models.NewEndpointsTableEngine(ctrl.DB)

		if err := c.ShouldBindUri(req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
			return
		}

		if endpoint, err := dbEngine.GetByID(req.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else if endpoint == nil {
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
			ctx := context.Background()
			ctrl.Config.GetRedis().Del(ctx, endpoint.Path).Result()

			c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": http.StatusText(http.StatusNoContent)})
			return
		}
	}
}
