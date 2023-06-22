package bootstrap

import (
	"fmt"
	"sync"

	"goblog/app/http/controllers"
	"goblog/pkg/logger"
	"time"

	cronV3 "github.com/robfig/cron/v3"
)

var mutex1 sync.Mutex
var mutex2 sync.Mutex
var mutex3 sync.Mutex

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
		// 获取互斥锁1
		mutex1.Lock()
		defer mutex1.Unlock()
		vc.SaveToMysql()

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
		// 获取互斥锁2
		mutex2.Lock()
		defer mutex2.Unlock()
		vc.DoSlice()
	})

	go c.AddFunc("0 20 * * *", func() {
		defer func() {
			if err := recover(); err != nil {
				logger.LogError(err.(error)) // 记录

				fmt.Printf("Recovered from panic: %v\n", err.(error))

			}
		}()
		fmt.Println("\n定时任务-切片昨天切片中的视频：每天8点执行一次", time.Now().Format("2006-01-02 15:04:05"))
		// 获取互斥锁3
		mutex3.Lock()
		defer mutex3.Unlock()
		vc.DoYestedaySlice()
	})
	c.Start()
	// 程序运行保持不退出
	select {}
}
