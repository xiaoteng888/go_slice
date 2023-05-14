package video

import (
	"goblog/app/models"
	"mime/multipart"
)

const TableName = "video_movies_data"

// Video 视频模型
type Video struct {
	models.BaseModel

	VideoName    string `gorm:"type:varchar(255);not null;" valid:"video_name"`
	Description  string `gorm:"type:varchar(255);" valid:"description"`
	Url          string `gorm:"type:varchar(255);not null" valid:"url"`
	MovieLength  string `gorm:"type:varchar(255);not null" valid:"movie_length"`
	Country      int64  `gorm:"default:9;index;comment:0:日本/1:韩国/5:国产/8:欧美/9:其他地区;" valid:"country"`
	VideoType    int64  `gorm:"default:0;index;comment:0:有码/1:无码;" valid:"video_type"`
	ShootingType int64  `gorm:"default:0;index;comment:0:专业拍摄 1:偷拍 2:自拍 3:业务拍摄;" valid:"shooting_type"`
	SubtitleType int64  `gorm:"default:0;comment:字幕类型(0:无 1:中文 2:英文 3:中英文 4:其他);" valid:"subtitle_type"`
	Number       string `gorm:"type:varchar(100);comment:番号;" valid:"number"`
	Producer     string `gorm:"type:varchar(255);comment:制作商;" valid:"producer"`
	Actor        string `gorm:"type:varchar(255);comment:主演;" valid:"actor"`
	PublishTime  string `gorm:"type:varchar(255);comment:发行时间;" valid:"publish_time"`
	SliceStatus  int64  `gorm:"not null;default:0;index;comment:是否切片:1是,0否" valid:"slice_status"`
	UpUrl        string `gorm:"type:varchar(255);not null" valid:"up_url"`

	UpVideo *multipart.FileHeader `gorm:"-" valid:"up_video" form:"up_video"`
}

const (
	STATUS_SUCCESS = 1
	STATUS_INVALID = 0
	STATUS_ON      = 2
	STATUS_FAILED  = 3
)

var VideoTypes = map[string]string{
	"0": "有码",
	"1": "无码",
}

var Countries = map[string]string{
	"0": "日本",
	"1": "韩国",
	"5": "中国大陆",
	"6": "中国台湾",
	"8": "欧美",
	"9": "其他地区",
}

var ShootingTypes = map[string]string{
	"0": "专业拍摄",
	"1": "偷拍",
	"2": "自拍",
	"3": "业务拍摄",
}

var SubtitleTypes = map[string]string{
	"0": "无字幕",
	"1": "中文字幕",
	"4": "其他字幕",
}

// CreatedAtDate 创建日期
func (video Video) CreatedAtDate() string {
	return video.CreatedAt.Format("2006-01-02")
}
