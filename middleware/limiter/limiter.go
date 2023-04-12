package limiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"GinGateway/constant"
)

func TokenBucketLimiter(buckLockType int) gin.HandlerFunc {
	return func(context *gin.Context) {
		var b TokenBucket
		switch buckLockType {
		case constant.BucketLockType_CAS:
			b = GetBucketCAS(1000, 1000, time.NewTicker(time.Millisecond))
		case constant.BucketLockType_Mutex:
			b = GetBucketMutex(1000, 1000, time.NewTicker(time.Millisecond))
		}
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
