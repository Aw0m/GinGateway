package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func TokenBucketLimiter() gin.HandlerFunc {
	return func(context *gin.Context) {
		b := GetBucket(1000, 1000, time.NewTicker(time.Millisecond))
		if ok := b.GetToken(); !ok {
			context.Status(http.StatusServiceUnavailable)
			context.Abort()
			fmt.Println("限流中间件：请求被拒绝")
		} else {
			context.Next()
			fmt.Println("限流中间件：请求成功")
		}
	}
}
