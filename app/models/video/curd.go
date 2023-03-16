package video

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
)

// Create 创建文章，通过 article.ID 来判断是否创建成功
func (video *Video) Create() (err error) {
	if err = model.DB.Create(&video).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

// 修改视频
func (video *Video) Update() (rowsAffected int64, err error) {
	result := model.DB.Table("videos").Save(&video)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}
	return result.RowsAffected, nil
}
