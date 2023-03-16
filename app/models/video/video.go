package video

import (
	"goblog/app/models"
	"mime/multipart"
)

// Video 视频模型
type Video struct {
	models.BaseModel

	Name string `gorm:"type:varchar(255);not null;" valid:"name"`
	Url  string `gorm:"type:varchar(255);not null;" valid:"url"`

	UpVideo *multipart.FileHeader `gorm:"-" valid:"up_video" form:"up_video"`
}

// CreatedAtDate 创建日期
func (video Video) CreatedAtDate() string {
	return video.CreatedAt.Format("2006-01-02")
}
