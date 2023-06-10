package file

import (
	"bytes"
	"context"
	"fmt"
	"goblog/app/models/video"
	"goblog/pkg/app"
	"goblog/pkg/helpers"
	pkgs3 "goblog/pkg/s3"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gogf/gf/util/gconv"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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
	end := strings.Index(strings.ToLower(inputVideo), ".mp4")
	if end == -1 {
		end = strings.Index(strings.ToLower(inputVideo), ".ts")
	}
	name := inputVideo[start:end]
	fmt.Print(name)
	outputDir := "./storage/movie/" + name
	//检查输出目录是否存在
	_, err := os.Stat(outputDir)
	if err == nil {
		err := os.RemoveAll(filepath.Join(outputDir))
		if err != nil {
			fmt.Println("删除目录出错:", err)
			return err
		}
	}

	//segmentLength := 20 //时长秒
	// 检查原始视频文件是否存在
	url := "." + inputVideo
	fmt.Print(url)
	_, err = os.Stat(url)
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

	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", url)
	output, err := cmd.Output()
	if err != nil {
		//logger.LogError(err)
		os.Exit(1)
		return err
	}
	duration := string(output)
	formattedDuration, err := formatDuration(duration)
	if err != nil {
		fmt.Println("转换出错:", err)
		return err
	}

	fmt.Println("格式化后的时长:", formattedDuration)

	_video.MovieLength = formattedDuration
	// 切片视频
	fmt.Println("开始切片视频...")

	// cmd = exec.Command("ffmpeg", "-i", url, "-c:v", "libx264", "-crf", "30", "-c:a", "copy", "-map", "0", "-f", "segment", "-segment_list", outputDir+"/playlist.m3u8", "-segment_time", gconv.String(segmentLength), outputDir+"/output_%03d.ts")

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Println(err)
	// 	logger.LogError(err)
	// 	os.Exit(1)
	// 	return err
	// }
	seconds, _ := strconv.ParseFloat(strings.TrimSpace(duration), 64)
	err = sliceVideo(url, outputDir, seconds)
	if err != nil {
		fmt.Println(err, gconv.Float64(output))
		return err
	}
	fmt.Println("切片完成！")

	// 显示切片文件信息
	files, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		//logger.LogError(err)
		return err
	}
	fmt.Println("切片文件列表：")
	dirName := path.Base(outputDir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// 获取文件扩展名
		ext := filepath.Ext(file.Name())
		if ext == ".m3u8" {

			_video.Url = "/xj/" + dirName + "/" + file.Name()
			_video.Update()
		}

		fmt.Println(file.Name())
	}

	// 显示进度条
	//bar := pb.StartNew(len(files))
	bucket := "yellowbook-media"
	if pkgs3.S3Client == nil {
		fmt.Print("初始化 S3 客户端----------\n")
		// 初始化 S3 客户端
		pkgs3.InitS3()
	}
	// 设置并发任务的最大数量
	maxConcurrency := 3
	var wg1 sync.WaitGroup
	// 创建通道，用于控制并发任务的数量
	concurrencyCh := make(chan struct{}, maxConcurrency)
	errCh := make(chan error, len(files))
	for i, info := range files {
		wg1.Add(1)
		_info := info
		// 启动一个协程执行任务
		go func(_info fs.DirEntry) {

			defer wg1.Done()

			// 控制并发任务的数量
			concurrencyCh <- struct{}{}
			defer func() { <-concurrencyCh }()

			// 模拟处理每个切片文件
			//bar.Increment()
			progress := float64(i+1) / float64(len(files)) * 100
			//fmt.Print("进度", progress, "%", "\n")
			// 上传文件到 S3
			key := path.Join("xj", dirName, _info.Name())
			data, err := os.Open(outputDir + "/" + _info.Name())
			//data, err := ioutil.ReadFile(info.Name())
			if err != nil {
				_video.SliceStatus = video.STATUS_FAILED
				_video.Update()
				errCh <- fmt.Errorf("上传S3修改状态失败: %s", err)

			}

			_, err = pkgs3.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    &key,
				Body:   data, //bytes.NewReader(data),
			})
			if err != nil {
				_video.SliceStatus = video.STATUS_FAILED
				_video.Update()
				errCh <- fmt.Errorf("上传S3修改状态失败: %s", err)

			}
			fmt.Printf("上传 %s/%s 到 S3://%s/%s\n", outputDir, _info.Name(), bucket, key)
			data.Close()
			// 计算进度百分比并发送消息
			fmt.Printf("进度百分之：%.2f\n", progress)
		}(_info)
	}
	// 等待所有任务完成
	go func() {
		wg1.Wait()
		close(errCh)
	}()

	// 检查错误通道，如果有错误则返回第一个错误
	for err := range errCh {
		return err
	}
	_video.SliceStatus = video.STATUS_SUCCESS
	_video.Update()
	//上传成功删除视频
	os.Remove(url)
	_, err = os.Stat(outputDir)
	if err == nil {
		err = os.RemoveAll(filepath.Join(outputDir))
		if err != nil {
			fmt.Println("删除目录出错:", err)
			return err
		}
	}
	return nil
}

// UptoS3 上传到上亚马逊s3
func UptoS3(dirPath string) error {

	if pkgs3.S3Client == nil {
		fmt.Print("初始化 S3 客户端----------")
		// 初始化 S3 客户端
		pkgs3.InitS3()
	}
	bucket := "yellowbook-media"
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

// 将秒数转换为00：00：00
func formatDuration(durationStr string) (string, error) {
	duration, err := strconv.ParseFloat(strings.TrimSpace(durationStr), 64)
	if err != nil {
		return "", err
	}

	totalSeconds := int(duration)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds), nil
}

// 分段切片
func sliceVideo(inputVideo, outputDir string, video_length float64) error {

	segmentCount := 5 // 分段数量

	// 设置并发任务的最大数量
	maxConcurrency := 3

	// 创建等待组，用于等待所有任务完成
	var wg sync.WaitGroup

	// 创建通道，用于控制并发任务的数量
	concurrencyCh := make(chan struct{}, maxConcurrency)

	// 计算每个分段的时长
	durationPerSegment := video_length / gconv.Float64(segmentCount)

	errCh := make(chan error, segmentCount)

	// 并发切片任务
	for i := 0; i < segmentCount; i++ {
		wg.Add(1)

		// 启动一个协程执行任务
		go func(segmentIndex int) {
			defer wg.Done()

			// 控制并发任务的数量
			concurrencyCh <- struct{}{}
			defer func() { <-concurrencyCh }()

			// 计算切片起始时间和结束时间
			startTime := durationPerSegment * float64(segmentIndex)
			endTime := startTime + durationPerSegment
			// 执行切片命令
			//cmd := exec.Command("ffmpeg", "-i", inputVideo, "-ss", fmt.Sprintf("%.2f", startTime), "-to", fmt.Sprintf("%.2f", endTime), "-c", "copy", "-f", "segment", "-segment_list", fmt.Sprintf("%s/playlist_%d.m3u8", outputDir, segmentIndex), fmt.Sprintf("%s/output_%d.ts", outputDir, segmentIndex))

			cmd := exec.Command("ffmpeg", "-i", inputVideo, "-ss", fmt.Sprintf("%.2f", startTime), "-to", fmt.Sprintf("%.2f", endTime), "-c:v", "libx264", "-crf", "30", "-c:a", "copy", "-map", "0", "-f", "segment", "-segment_list", outputDir+"/playlist.m3u8", "-segment_time", gconv.String(20), outputDir+"/output_"+gconv.String(segmentIndex)+"%03d.ts")

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				errCh <- fmt.Errorf("切片视频分段 %d 出错: %s", segmentIndex, err)
				return
			}

			fmt.Printf("视频分段 %d 切片完成\n", segmentIndex)
		}(i)
	}

	// 等待所有任务完成
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// 检查错误通道，如果有错误则返回第一个错误
	for err := range errCh {
		return err
	}

	fmt.Println("所有切片任务已完成")
	return nil
}
