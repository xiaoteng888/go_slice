package video

import (
	"fmt"
	"goblog/pkg/config"
	"goblog/pkg/logger"
	"goblog/pkg/model"
)

// Create 创建文章，通过 article.ID 来判断是否创建成功
func (video *Video) Create() (err error) {
	if err = model.DB.Table(TableName).Create(&video).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

// 修改视频
func (video *Video) Update() (rowsAffected int64, err error) {
	result := model.DB.Table(TableName).Save(&video)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

// 修改视频
func (video *Video) Update1() (rowsAffected int64, err error) {
	result := model.DB1.Table(TableName).Save(&video)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}

// 获取未切片视频
func GetMp4() ([]Video, error) {
	var videos []Video
	result := model.DB.Table(TableName).Where("win_no = ?", config.Env("WIN_NO", 0)).Where("slice_status IN ?", []int{STATUS_INVALID}).Find(&videos)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return videos, err
	}
	return videos, nil
}

// 根据名称获取视频
func Get(name string) (Video, error) {
	var video Video
	result := model.DB.Table(TableName).Where("video_name = ?", name).Find(&video)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return video, err
	}
	fmt.Print(666)
	return video, nil
}

// 获取上传失败视频
func GetYestedayMp4() ([]Video, error) {
	var videos []Video
	// 获取昨天的起始时间和结束时间
	//yesterdayStart := time.Now().AddDate(0, 0, -7).Format("2006-01-02 00:00:00")
	//yesterdayEnd := time.Now().AddDate(0, 0, -3).Format("2006-01-02 23:59:59")
	result := model.DB.Table(TableName).Where("win_no = ?", config.Env("WIN_NO", 0)).Where("slice_status = ?", STATUS_FAILED).Find(&videos)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return videos, err
	}
	return videos, nil
}

// 获取正在切片视频
func GetDoMp4() ([]Video, error) {
	var videos []Video
	result := model.DB.Table(TableName).Where("win_no = ?", config.Env("WIN_NO", 0)).Where("slice_status = ?", STATUS_ON).Find(&videos)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return videos, err
	}
	return videos, nil
}
