package userService

import (
	"GinGateway/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct {
	userID   int64
	username string
}

func UserRouter() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery(), gin.Logger())
	e.POST("/api/login", func(context *gin.Context) {
		err := middleware.ForwardHandler(context.Writer, context.Request, "user")
		if err != nil {
			if err.Error() == "not exist" {
				context.Status(http.StatusNotFound)
			} else {
				context.Status(http.StatusServiceUnavailable)
			}
		}
	})

	e.Use(middleware.Authorize())
	e.Use(middleware.TokenBucketLimiter())
	e.Use(middleware.RouteForward("user"))
	return e
}
