package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-bridging-service/config"
	"github.com/hexcraft-biz/controller"
)

type Common struct {
	*controller.Prototype
}

func NewCommon(cfg config.ConfigInterface) *Common {
	return &Common{
		Prototype: controller.New("common", cfg.GetDB()),
	}
}

func (ctrl *Common) NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
	}
}

func (ctrl *Common) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
	}
}
