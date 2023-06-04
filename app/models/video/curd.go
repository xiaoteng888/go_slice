package video

import (
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

// 获取未切片视频
func GetMp4() ([]Video, error) {
	var videos []Video
	result := model.DB.Table(TableName).Where("slice_status IN ?", []int{STATUS_INVALID, STATUS_FAILED}).Find(&videos)
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
	return video, nil
}
