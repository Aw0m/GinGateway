package noteService

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wxprojectApiGateway/middleware"
)

func NoteRouter() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	//e.Use(middleware.Authorize())
	e.Use(middleware.TokenBucketLimiter())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server note",
			},
		)
	})
	return e
}
