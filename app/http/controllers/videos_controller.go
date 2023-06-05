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
	"path/filepath"

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
		VideoName:    r.PostFormValue("name"),
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
		_video.UpUrl = video
		_video.Update()
		fmt.Print("----------7")
		//上传成功，开始切片
		// go func() {
		// 	// 将视频文件进行切片
		// 	if err := files.Slice(video, _video); err != nil {
		// 		logger.LogError(err)
		// 	}
		// }()

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

// Slice视频切片页面
func (*VideosController) Slice(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "videos.slice")
}

func (*VideosController) SaveToMysql() {
	// 把文件夹里的视频创建到表里
	PathToMysql()
}

// DoSlice 执行视频切片操作
func (*VideosController) DoSlice() {

	// 获取未切片的视频
	videos, err := video.GetMp4()

	if err != nil {
		logger.LogError(err)
		//w.WriteHeader(http.StatusInternalServerError)
		fmt.Print("500 服务器内部错误")
	} else {
		if len(videos) == 0 {
			fmt.Print("暂无视频可切片 \n")
			return
		}

		for _, v := range videos {
			v.SliceStatus = 2
			v.Update()
		}

		for _, _video := range videos {
			if _video.UpUrl == "" {
				fmt.Print("视频不存在 \n")
				return
			}

			// 将视频文件进行切片
			fmt.Print("开始切片---- 视频名：", _video.VideoName, "视频位置：", _video.UpUrl, "\n")
			err := files.Slice(_video.UpUrl, _video)

			if err != nil {
				logger.LogError(err)
				fmt.Print("切片报错", err, "\n")
			} else {
				fmt.Print("视频名:", _video.VideoName, "视频位置：", _video.UpUrl, "切片完成\n")
			}
		}
	}

}

// 循环处理文件路径
func PathToMysql() {
	root := "public/uploads/movies"
	uproot := "public/uploads/upmovies"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// 先查询数据库有无这个视频
			_one, _ := video.Get(info.Name())
			if _one.ID > 0 {
				//查到就不执行
				fmt.Println("上个扫描任务正在执行：影片---", _one.VideoName, "---已经存在或者正在切片上传S3")
				return nil
			}
			fmt.Print(111)
			// 先把视频转移，再存入数据库
			rootFile, err := os.Open(path)
			fmt.Print(222)
			if err != nil {
				return err
			}
			fmt.Print(333)
			uprootFile := filepath.Join(uproot, info.Name())
			targetFile, err := os.Create(uprootFile)
			if err != nil {
				return err
			}
			fmt.Print(444)
			defer targetFile.Close()
			fmt.Print(555)
			_, err = io.Copy(targetFile, rootFile)
			if err != nil {
				return err
			}
			fmt.Print(666)
			rootFile.Close()
			// 删除原路径视频
			fmt.Println("原路径视频：", path)
			err = os.Remove(path)
			if err != nil {
				fmt.Print(err)
				return err
			}
			fmt.Println(uprootFile)
			_video := video.Video{
				UpUrl:     "/" + filepath.ToSlash(uprootFile),
				VideoName: info.Name(),
			}
			_video.Update()
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
