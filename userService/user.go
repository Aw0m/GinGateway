package userService

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wxprojectApiGateway/middleware"
)

type User struct {
	userID   int64
	username string
}

func UserRouter() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.POST("/api/login", func(context *gin.Context) {
		err := middleware.ForwardHandler(context.Writer, context.Request, "user")
		if err.Error() == "not exist" {
			context.Status(http.StatusNotFound)
		} else if err != nil {
			context.Status(http.StatusServiceUnavailable)
		}
	})

	e.Use(middleware.Authorize())
	e.Use(middleware.RouteForward("user"))
	return e
}
