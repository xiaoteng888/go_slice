package controllers

import (
	"fmt"
	"goblog/app/models/video"
	"goblog/app/requests"
	files "goblog/pkg/file"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
	"sync"
	"time"

	"github.com/gogf/gf/util/gconv"
)

type VideosController struct {
	BaseController
}

// Create 上传视频页面
func (*VideosController) Create(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	n, ok := query["n"]
	data := view.D{
		"videoTypes":    video.VideoTypes,
		"countries":     video.Countries,
		"shootingTypes": video.ShootingTypes,
		"subtitleTypes": video.SubtitleTypes,
	}
	if ok {
		data["Notice"] = n
	}
	view.Render(w, data, "videos.create", "videos._form_field")
}

// Store 视频上传操作
func (*VideosController) Store(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	// 1. 初始化数据
	//currentUser := auth.User()
	fmt.Print("----------1")
	_video := video.Video{
		Name:         r.PostFormValue("name"),
		Description:  r.PostFormValue("description"),
		Country:      gconv.Int64(r.PostFormValue("country")),
		VideoType:    gconv.Int64(r.PostFormValue("video_type")),
		ShootingType: gconv.Int64(r.PostFormValue("shooting_type")),
		SubtitleType: gconv.Int64(r.PostFormValue("subtitle_type")),
		Number:       r.PostFormValue("number"),
		Producer:     r.PostFormValue("producer"),
		Actor:        r.PostFormValue("actor"),
		PublishTime:  r.PostFormValue("publish_time"),
	}

	file, handler, err := r.FormFile("uploadFile")

	if err != nil {
		fmt.Print("----------2")
		data := view.D{
			"Video": _video,
			"err":   "没有选择文件",
		}
		view.Render(w, data, "videos.create", "videos._form_field")
		return
	}
	_video.UpVideo = handler
	fmt.Print("----------3")
	// 2. 表单验证
	errors := requests.ValidateVideoForm(_video)
	fmt.Print("----------4")
	// 3. 检测错误
	if len(errors) == 0 {
		video, err := files.SaveUploadVideo(r, handler, file)
		fmt.Print("----------5")
		if err != nil {
			logger.LogError(err)
			fmt.Print("----------6")
			data := view.D{
				"Video": _video,
				"err":   err,
			}
			view.Render(w, data, "videos.create", "videos._form_field")
			return
		}
		fmt.Print("----------7")
		//上传成功，开始切片
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			// 将视频文件进行切片
			defer wg.Done()
			if err := files.Slice(video, _video); err != nil {
				logger.LogError(err)
			}
		}()
		time.Sleep(1 * time.Second) // 等待1秒钟，确保goroutine有足够的时间来执行
		wg.Wait()
		fmt.Print("----------8")
		indexURL := route.Name2URL("videos.create")
		http.Redirect(w, r, indexURL+"?n=1", http.StatusFound)
	} else {
		fmt.Print("----------9")
		view.Render(w, view.D{
			"Video":  _video,
			"Errors": errors,
		}, "videos.create", "videos._form_field")
	}
}
