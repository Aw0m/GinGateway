package discovery

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"time"
)

func GetService(serviceName string) (string, bool) {
	for {
		// 获得当前所有的 url
		serviceURLs := rds.SMembers(ctx, serviceName).Val()
		if len(serviceURLs) == 0 {
			return "", false
		}
		rand.Seed(time.Now().Unix())
		url := serviceURLs[rand.Intn(len(serviceURLs))]
		// 判断该url是否还存在，即该服务是否挂掉了
		if _, err := rds.Get(ctx, url).Result(); err != nil {
			if err == redis.Nil {
				// 没找到该url的key则说明挂了
				log.Printf("该url: %s 已挂！\n", url)
				rds.SRem(ctx, serviceName, url)
			} else {
				log.Println("error!", err.Error())
			}
			continue
		}
		// 存在则直接返回url
		fmt.Printf("最终URL为： %s;  ", url)
		return url, true
	}
}
