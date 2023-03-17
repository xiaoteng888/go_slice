package s3

import (
	"context"
	pkgconfig "goblog/pkg/config"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gogf/gf/util/gconv"
)

// s3Client 对象
var S3Client *s3.Client

// Initialize 初始化S3
func Initialize() {
	InitS3()
}
func InitS3() {
	// 使用您的 AWS access key 和 secret key 来创建 AWS 配置对象
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			gconv.String(pkgconfig.Env("AWS_ACCESS_KEY")),
			gconv.String(pkgconfig.Env("AWS_SECRET_ACCESS_KEY")),
			"",
		)),
	)
	if err != nil {
		panic(err)
	}

	// 使用 AWS 配置对象来创建 S3 客户端
	S3Client = s3.NewFromConfig(cfg)

	// 现在，您可以使用 s3Client 对象来调用 S3 API 了
	// 例如，使用 ListBuckets API 列出 S3 中的所有存储桶
	// result, err := s3Client.ListBuckets(context.TODO(), nil)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, bucket := range result.Buckets {
	// 	fmt.Println(*bucket.Name)
	// }
}
