package video

import (
	"goblog/app/models"
	"mime/multipart"
)

// Video 视频模型
type Video struct {
	models.BaseModel

	Name   string `gorm:"type:varchar(255);not null;" valid:"name"`
	Url    string `gorm:"type:varchar(255);not null;" valid:"url"`
	Status int64  `gorm:"not null;default:0;index;comment:是否切片:1是,0否" valid:"status"`

	UpVideo *multipart.FileHeader `gorm:"-" valid:"up_video" form:"up_video"`
}

// CreatedAtDate 创建日期
func (video Video) CreatedAtDate() string {
	return video.CreatedAt.Format("2006-01-02")
}
