package file

import (
	"bytes"
	"context"
	"fmt"
	"goblog/app/models/video"
	"goblog/pkg/app"
	"goblog/pkg/helpers"
	"goblog/pkg/logger"
	pkgs3 "goblog/pkg/s3"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cheggaaa/pb"
	"github.com/gogf/gf/util/gconv"
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

// 切片
func Slice(inputVideo string, _video video.Video) error {

	// 定义输入视频文件名、输出目录名和切片长度
	//inputVideo := "./storage/" + name + ".mp4"
	start := strings.LastIndex(inputVideo, "/") + 1
	end := strings.Index(inputVideo, ".mp4")
	name := inputVideo[start:end]
	outputDir := "./storage/movie/" + name
	segmentLength := 90 //时长秒
	// 检查原始视频文件是否存在
	url := "." + inputVideo
	fmt.Print(url)
	_, err := os.Stat(url)
	if os.IsNotExist(err) {
		fmt.Print("视频文件不存在", url)
		return err
	}
	// 创建切片输出目录
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		//logger.LogError(err)
		return err
	}
	// 获取视频时长
	fmt.Println("获取视频时长...")
	//D:/ffmpeg/ffmpeg-master-latest-win64-gpl-shared/bin/ffprobe
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", url)
	output, err := cmd.Output()
	if err != nil {
		//logger.LogError(err)
		os.Exit(1)
		return err
	}
	duration := string(output)
	fmt.Printf("视频时长为：%s", duration)
	_video.Length = duration

	// 切片视频
	fmt.Println("开始切片视频...")
	//D:/ffmpeg/ffmpeg-master-latest-win64-gpl-shared/bin/ffmpeg
	cmd = exec.Command("ffmpeg", "-i", url, "-codec", "copy", "-vbsf", "h264_mp4toannexb", "-map", "0", "-f", "segment", "-segment_list", outputDir+"/playlist.m3u8", "-segment_time", gconv.String(segmentLength), outputDir+"/output_%03d.ts")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		logger.LogError(err)
		os.Exit(1)
		return err
	}
	fmt.Println("切片完成！")
	_video.Status = video.STATUS_SUCCESS
	// 删除原始文件
	os.Remove(url)

	// 显示切片文件信息
	files, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		//logger.LogError(err)
		return err
	}
	fmt.Println("切片文件列表：")
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fmt.Println(file.Name())
	}

	// 显示进度条
	bar := pb.StartNew(len(files))
	for i := 0; i < len(files); i++ {
		// 模拟处理每个切片文件
		bar.Increment()
	}
	bar.Finish()

	err = UptoS3(outputDir, _video)
	return err
}

// UptoS3 上传到上亚马逊s3
func UptoS3(dirPath string, _video video.Video) error {

	if pkgs3.S3Client == nil {
		fmt.Print("初始化 S3 客户端----------")
		// 初始化 S3 客户端
		pkgs3.InitS3()
	}
	bucket := "pianpian"
	// prefix := "xj/IMG_3987/"
	// directory := "/storage/movie/IMG_3987/"
	fmt.Print("output路径：", dirPath)
	// 上传所有文件
	dirName := path.Base(dirPath)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// 读取文件内容
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// 上传文件到 S3
			key := filepath.Join("xj/"+dirName, info.Name())
			// 获取文件扩展名
			ext := filepath.Ext(info.Name())
			if ext == ".m3u8" {
				_video.Url = "/" + key
				_video.Update()
			}

			_, err = pkgs3.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    &key,
				Body:   bytes.NewReader(data),
			})
			if err != nil {
				return err
			}

			fmt.Printf("上传 %s 到 S3://%s/%s\n", path, bucket, key)
		}

		return nil
	})
	if err != nil {
		//logger.LogError(err)
		return err
	}

	fmt.Println("所有文件上传成功")
	return nil
}
