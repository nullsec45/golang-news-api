package cloudflare

import (
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"fmt"
	"os"
	"github.com/gofiber/fiber/v2/log"
	"github.com/nullsec45/golang-news-api/config"
	"context"
)

var code string
var err error

type CloudflareR2Adapter interface {
	UploadImage(req *entity.FileUploadEntity) (string, error)
}

type cloudflareR2Adapter struct {
	Client *s3.Client
	Bucket string
	BaseURL string
}

func (c * cloudflareR2Adapter) UploadImage(req *entity.FileUploadEntity) (string, error) {
	openedFile, err := os.Open(req.Path)
	if err != nil {
		code = "[CLOUDFLARE R2] - UploadImage - 1"
		log.Errorw(code, err)
		return "", err
	}

	defer openedFile.Close()

	_, err = c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(c.Bucket),
		Key:         aws.String(req.Name),
		Body:        openedFile,
		ContentType: aws.String("image/jpeg"),
	})

	if err != nil {
		code = "[CLOUDFLARE R2] UploadImage - 2"
		log.Errorw(code, err)
		return "", err
	}
	return fmt.Sprintf("%s/%s", c.BaseURL, req.Name), nil
}


func NewCloudflareR2Adapter(client *s3.Client, cfg *config.Config) CloudflareR2Adapter{
	clientBase := s3.NewFromConfig(cfg.LoadAWSConfig(), func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2.AccountID))
	})

	return &cloudflareR2Adapter{
		Client: clientBase,
		Bucket: cfg.R2.Name,
		BaseURL: cfg.R2.PublicURL,
	}
}