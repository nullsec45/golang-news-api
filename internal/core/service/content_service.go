package service

import (
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"github.com/nullsec45/golang-news-api/internal/adapter/repository"
	"context"
	"github.com/nullsec45/golang-news-api/internal/adapter/cloudflare"
	"github.com/nullsec45/golang-news-api/config"
	"github.com/gofiber/fiber/v2/log"
)

type ContentService interface {
	GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	UpdateContent(ctx context.Context, req *entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
	UploadImageR2(ctx context.Context, req entity.FileUploadEntity)(string, error)
}

type contentService struct {
	contentRepo repository.ContentRepository
	cfg               *config.Config
	r2                cloudflare.CloudflareR2Adapter
}

func (c *contentService) GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error) {
	results, totalData, totalPages, err := c.contentRepo.GetContents(ctx, query)
	if err != nil {
		code = "[SERVICE] GetContents - 1"
		log.Errorw(code, err)
		return nil, 0,0, err
	}	

	return results, totalData, totalPages, nil
}

func (c *contentService) CreateContent(ctx context.Context, req entity.ContentEntity) error {	
	err = c.contentRepo.CreateContent(ctx, req)
	
	if err != nil {
		code = "[SERVICE] CreateContent - 1"
		log.Errorw(code, err)
		return  err
	}	

	return  nil
}

func (c *contentService) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	result, err := c.contentRepo.GetContentByID(ctx, id)
	if err != nil {
		code = "[SERVICE] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func (c *contentService) UpdateContent(ctx context.Context, req *entity.ContentEntity) error {
	err = c.contentRepo.UpdateContent(ctx, req)
	if err != nil {
		code = "[SERVICE] UpdateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}		

func (c *contentService) DeleteContent(ctx context.Context, id int64) error {
	err = c.contentRepo.DeleteContent(ctx, id)
	if err != nil {
		code = "[SERVICE] DeleteContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *contentService) UploadImageR2(ctx context.Context, req entity.FileUploadEntity) (string, error) {
	urlImage, err := c.r2.UploadImage(&req)
	if err != nil {
		code = "[SERVICE] UploadImageR2 - 1"
		log.Errorw(code, err)
		return "", err
	}
	return urlImage, nil
}


func NewContentService(repo repository.ContentRepository, cfg *config.Config, r2 cloudflare.CloudflareR2Adapter) ContentService {
	return &contentService{
		contentRepo: repo,
		cfg:         cfg,
		r2:          r2,
	}
}