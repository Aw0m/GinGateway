package workService

import (
	"GinGateway/middleware"
	"GinGateway/service/discovery"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WorkRouter() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery(), gin.Logger())

	e.GET("/", func(c *gin.Context) {

		res, _ := discovery.GetService("work")
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": http.StatusOK,
				"res":  res,
			},
		)
		c.Abort()
	})
	e.Use(middleware.Authorize())
	e.Use(middleware.TokenBucketLimiter())
	e.Use(middleware.RouteForward("work"))
	return e
}
