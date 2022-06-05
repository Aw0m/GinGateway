package main

import (
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
	"wxprojectApiGateway/noteService"
	"wxprojectApiGateway/service/discovery"
	"wxprojectApiGateway/userService"
	"wxprojectApiGateway/workService"
)

var (
	g errgroup.Group
)

func init() {
	// 初始化服务注册和服务发现中心的配置
	discovery.InitDiscovery("service/discovery/config.yaml")
}

func main() {
	server01 := &http.Server{
		Addr:         ":8080",
		Handler:      userService.UserRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8081",
		Handler:      noteService.NoteRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server03 := &http.Server{
		Addr:         ":8082",
		Handler:      workService.WorkRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	g.Go(func() error {
		return server01.ListenAndServe()
	})

	g.Go(func() error {
		return server02.ListenAndServe()
	})
	g.Go(func() error {
		return server03.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
