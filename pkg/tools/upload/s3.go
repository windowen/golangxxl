package upload

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"liveJob/pkg/common/config"
	"liveJob/pkg/zlogger"
)

// 默认 ACL
const defaultAcl = "private"

// Uploader 结构体包含上传所需的所有配置
type Uploader struct {
	bucket   string // 储存桶
	acl      string // 权限
	metadata map[string]*string
	s3Client *s3.S3
}

// NewUploader 创建一个新的 Uploader，并应用 Option 配置
func NewUploader(opts ...Option) *Uploader {
	s3Cfg := config.Config.S3

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(s3Cfg.AccessKeyID, s3Cfg.SecretAccessKey, ""),
		Region:           aws.String(s3Cfg.Region),
		Endpoint:         aws.String(s3Cfg.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
	}

	u := &Uploader{
		bucket:   s3Cfg.Bucket,
		acl:      defaultAcl,
		metadata: nil,
		s3Client: s3.New(sess),
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

// CheckFile 文件格式、大小检查
func (u *Uploader) CheckFile(file *multipart.FileHeader) ([]byte, error) {
	// todo upload_exceed_limit
	src, err := file.Open()
	if err != nil {
		return nil, err
	}

	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			zlogger.Errorf("failed to close file stream, err: %v", err)
		}
	}(src)

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (u *Uploader) Upload(ctx context.Context, fileName string, fileData []byte) error {
	_, err := u.s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:   aws.String(u.bucket),
		Key:      aws.String(fileName),
		Body:     bytes.NewReader(fileData),
		ACL:      aws.String(u.acl),
		Metadata: u.metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// Option 是一个函数类型，用于配置 Uploader
type Option func(*Uploader)

// WithBucket 设置上传的S3存储桶
func WithBucket(bucket string) Option {
	return func(u *Uploader) {
		u.bucket = bucket
	}
}

// WithACL 设置文件的访问控制列表（ACL）
func WithACL(acl string) Option {
	return func(u *Uploader) {
		u.acl = acl
	}
}

// WithMetadata 设置文件的元数据
func WithMetadata(metadata map[string]*string) Option {
	return func(u *Uploader) {
		u.metadata = metadata
	}
}
