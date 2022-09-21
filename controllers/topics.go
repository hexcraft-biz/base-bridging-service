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

type Topics struct {
	*controller.Prototype
	Config config.ConfigInterface
}

func NewTopics(cfg config.ConfigInterface) *Topics {
	return &Topics{
		Prototype: controller.New("topics", cfg.GetDB()),
		Config:    cfg,
	}
}

func (ctrl *Topics) NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
	}
}

type ListTopicsParams struct {
	Limit  uint64 `form:"limit,default=20"`
	Offset uint64 `form:"offset,default=0"`
}

func (ctrl *Topics) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		listParams := new(ListTopicsParams)
		if err := c.ShouldBindQuery(listParams); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
			return
		}
		pg := model.NewPagination(listParams.Offset, listParams.Limit)

		if endpoints, err := models.NewTopicsTableEngine(ctrl.DB).List(pg); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusOK, endpoints)
			return
		}
	}
}

type TargetTopic struct {
	ID string `uri:"id" binding:"required"`
}

func (ctrl *Topics) GetOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		var targetTopic TargetTopic
		if err := c.ShouldBindUri(&targetTopic); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewTopicsTableEngine(ctrl.DB).GetByID(targetTopic.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				if absRes, absErr := entityRes.GetAbsTopic(); absErr != nil {
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

type createTopicParams struct {
	Name string `json:"name" binding:"required,min=5,max=256"`
}

func (ctrl *Topics) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params createTopicParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewTopicsTableEngine(ctrl.DB).Insert(params.Name); err != nil {
			if myErr, ok := err.(*mysql.MySQLError); ok && myErr.Number == 1062 {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": http.StatusText(http.StatusConflict)})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		} else {
			if absRes, absErr := entityRes.GetAbsTopic(); absErr != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": absErr.Error()})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusCreated, absRes)
				return
			}
		}
	}
}

func (ctrl *Topics) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, dbEngine := new(TargetTopic), models.NewTopicsTableEngine(ctrl.DB)

		if err := c.ShouldBindUri(req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
			return
		}

		if topic, err := dbEngine.GetByID(req.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else if topic == nil {
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
			etrs, _ := models.NewEndpointTopicRelsTableEngine(ctrl.DB).GetByTopicID(topic.ID.String())
			ctx := context.Background()
			for i := 0; i < len(etrs); i++ {
				ctrl.Config.GetRedis().Del(ctx, etrs[i].Path).Result()
			}

			c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": http.StatusText(http.StatusNoContent)})
			return
		}
	}
}
