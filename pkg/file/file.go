package file

import (
	"fmt"
	"goblog/pkg/app"
	"goblog/pkg/helpers"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func SaveUploadVideo(r *http.Request, file *multipart.FileHeader, uploadfile multipart.File) (string, error) {
	var video string
	//确保目录存在，不存在创建
	publicPath := "public"
	dirName := fmt.Sprintf("/uploads/movies/%s/", app.TimenowInTimezone().Format("2006/01/02"))
	os.MkdirAll(publicPath+dirName, 0777)
	// 保存文件
	fileName := randomNameFromUploadFile(file)
	// public/uploads/movies/2021/12/22/nFDacgaWKpWWOmOt.png
	videoPath := publicPath + dirName + fileName
	savefile, err := os.OpenFile(videoPath, os.O_WRONLY|os.O_CREATE, 0777)
	//t, err := os.Create(avatarPath)
	if err != nil {
		return video, err
	}

	if _, err := io.Copy(savefile, uploadfile); err != nil {
		return video, err
	}
	defer savefile.Close()
	defer uploadfile.Close()
	return "/" + videoPath, nil
}

func randomNameFromUploadFile(file *multipart.FileHeader) string {
	return helpers.RandomString(16) + filepath.Ext(file.Filename)
}
