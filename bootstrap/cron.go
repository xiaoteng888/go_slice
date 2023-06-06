package bootstrap

import (
	"fmt"
	"sync"

	"goblog/app/http/controllers"
	"goblog/pkg/logger"
	"time"

	cronV3 "github.com/robfig/cron/v3"
)

var mutex sync.Mutex

func SetupCron() {

	c := cronV3.New(cronV3.WithSeconds()) //精确到秒
	vc := new(controllers.VideosController)
	// 扫描文件夹
	go c.AddFunc("@every 300s", func() {
		defer func() {
			if err := recover(); err != nil {
				logger.LogError(err.(error)) // 记录

				fmt.Printf("Recovered from panic: %v\n", err.(error))

			}
		}()
		fmt.Println("\n定时任务-扫描视频：每300秒执行一次", time.Now().Format("2006-01-02 15:04:05"))
		// 获取互斥锁
		mutex.Lock()
		vc.SaveToMysql()
		// 释放互斥锁
		mutex.Unlock()
	})

	// 执行切片
	go c.AddFunc("@every 600s", func() {
		defer func() {
			if err := recover(); err != nil {
				logger.LogError(err.(error)) // 记录

				fmt.Printf("Recovered from panic: %v\n", err.(error))

			}
		}()
		fmt.Println("\n定时任务-切片上传S3：每600秒执行一次", time.Now().Format("2006-01-02 15:04:05"))
		// 获取互斥锁
		mutex.Lock()
		vc.DoSlice()
		// 释放互斥锁
		mutex.Unlock()
	})
	c.Start()

}
