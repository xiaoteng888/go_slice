package bootstrap

import (
	"fmt"

	"goblog/app/http/controllers"
	"goblog/pkg/logger"
	"time"

	cronV3 "github.com/robfig/cron/v3"
)

func SetupCron() {

	c := cronV3.New(cronV3.WithSeconds()) //精确到秒
	vc := new(controllers.VideosController)
	go c.AddFunc("@every 300s", func() {
		defer func() {
			if err := recover(); err != nil {
				logger.LogError(err.(error)) // 记录

				fmt.Printf("Recovered from panic: %v\n", err.(error))

			}
		}()
		fmt.Println("\n定时任务-扫描切片：每300秒执行一次", time.Now().Format("2006-01-02 15:04:05"))
		vc.DoSlice()
	})

	c.Start()

}
