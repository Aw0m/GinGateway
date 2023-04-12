package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"

	"GinGateway/constant"
	"GinGateway/middleware/limiter"
	"GinGateway/middleware/router"
	"GinGateway/service/discovery"
)

func Init() {
	flag.IntVar(&port, "port", 8081, "the port to start gin")
	flag.Parse()
	// 初始化服务注册和服务发现中心的配置
	discovery.InitDiscovery("config.yaml")
}

var (
	g    errgroup.Group
	port int
)

func main() {
	Init()
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.Use(limiter.TokenBucketLimiter(constant.BucketLockType_CAS))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Use(router.RouteForward("work", constant.HttpType_HTTP))

	g.Go(func() error {
		return r.Run(fmt.Sprintf(":%d", port))
	})

	//g.Go(func() error {
	//	return r.Run(fmt.Sprintf(":%d", port))
	//})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
