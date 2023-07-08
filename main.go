package main

import (
	"embed"
	"fmt"
	"goblog/app/http/middlewares"
	"goblog/app/models/video"
	"goblog/bootstrap"
	"goblog/config"
	c "goblog/pkg/config"
	"goblog/pkg/logger"
	"net/http"
)

//go:embed resources/views/articles/*
//go:embed resources/views/auth/*
//go:embed resources/views/categories/*
//go:embed resources/views/layouts/*
//go:embed resources/views/videos/*
var tplFS embed.FS

//go:embed public/*
var staticFS embed.FS

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {

	// 初始化 SQL
	bootstrap.SetupDB()

	// 初始化模板加载
	bootstrap.SetupTemplate(tplFS)

	// 初始化路由绑定
	router := bootstrap.SetupRoute(staticFS)
	// 初始化切片
	videos, err := video.GetDoMp4()
	if err != nil {
		logger.LogError(err)
		fmt.Print("500 服务器内部错误")
	}
	if len(videos) > 0 {
		for _, v := range videos {
			v.SliceStatus = 0
			v.Update()
		}
	}
	// 定时任务
	bootstrap.SetupCron()

	err = http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)
}
