package controllers

import (
	"fmt"
	"goblog/app/models/video"
	"goblog/app/requests"
	files "goblog/pkg/file"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"io"
	"net/http"
	"os"
)

type VideosController struct {
	BaseController
}

// Create 上传视频页面
func (*VideosController) Create(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	n, ok := query["n"]
	data := view.D{}
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
	_video := video.Video{
		Name: r.PostFormValue("name"),
	}

	file, handler, err := r.FormFile("uploadFile")

	if err != nil {
		data := view.D{
			"Video": _video,
			"err":   "没有选择文件",
		}
		view.Render(w, data, "videos.create", "videos._form_field")
		return
	}
	_video.UpVideo = handler

	// 2. 表单验证
	errors := requests.ValidateVideoForm(_video)

	// 3. 检测错误
	if len(errors) == 0 {
		video, err := files.SaveUploadVideo(r, handler, file)
		if err != nil {
			logger.LogError(err)

			data := view.D{
				"Video": _video,
				"err":   err,
			}
			view.Render(w, data, "videos.create", "videos._form_field")
			return
		}
		_video.Url = video
		_video.Update()

		indexURL := route.Name2URL("videos.create")
		fmt.Print(indexURL + "?n=1")
		http.Redirect(w, r, indexURL+"?n=1", http.StatusFound)
	} else {
		view.Render(w, view.D{
			"Video":  _video,
			"Errors": errors,
		}, "videos.create", "videos._form_field")
	}
}

// 处理 /upload  逻辑
func (*VideosController) Upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法

	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		data := view.D{

			"err": "没有选择文件",
		}
		view.Render(w, data, "videos.create", "videos._form_field")
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("/storage/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

}
