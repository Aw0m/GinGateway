package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func TokenBucketLimiter() gin.HandlerFunc {
	return func(context *gin.Context) {
		b := GetBucket(1000, 1000, time.NewTicker(time.Millisecond))
		if ok := b.GetToken(); !ok {
			context.Status(http.StatusServiceUnavailable)
			context.Abort()
			log.Println("请求被拒绝")
		} else {
			context.Next()
			log.Println("请求成功")
		}
	}
}
