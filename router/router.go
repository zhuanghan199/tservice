package router

import (
	"net/http"
	"tservice/common/logger"
	. "tservice/services/auth"

	"github.com/gin-gonic/gin"
)

func Route() {
	r := gin.Default()
	api := r.Group("/v1")
	{
		api.POST("/auth", gin.WrapH(AuthHandler()))
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "request not found")
	})
	logger.Debugln("listen 8080")
	r.Run(":8080")
}
