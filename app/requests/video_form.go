package requests

import (
	"path"
	"strings"

	"github.com/thedevsaddam/govalidator"

	"goblog/app/models/video"
)

// ValidateVideoForm 验证表单，返回 errs 长度等于零即通过
func ValidateVideoForm(data video.Video) map[string][]string {

	// 1. 定制认证规则
	rules := govalidator.MapData{
		"name": []string{"required", "max_cn:40"},
	}

	// 2. 定制错误消息
	messages := govalidator.MapData{
		"name": []string{
			"required:影片名为必填项",
			"max_cn:影片名长度需小于 40",
		},
	}

	// 3. 配置初始化
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid", // 模型中的 Struct 标签标识符
		Messages:      messages,
	}

	// 4. 开始验证
	errs := govalidator.New(opts).ValidateStruct()
	// 5. 自己写一个验证上传视频
	//检查图片后缀
	ext := strings.ToLower(path.Ext(data.UpVideo.Filename))

	if ext != ".mp4" {
		errs["url"] = append(errs["url"], "视频只能是MP4")
	}

	return errs
}
