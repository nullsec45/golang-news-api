package repository

import (
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"context"
	"gorm.io/gorm"
	"github.com/nullsec45/golang-news-api/internal/core/domain/model"
	"github.com/gofiber/fiber/v2/log"
	"strings"
)

type ContentRepository interface {
	GetContents(ctx context.Context) ([]entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	UpdateContent(ctx context.Context, req *entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
}

type contentRepository struct {
	db *gorm.DB
}

func (c *contentRepository) GetContents(ctx context.Context) ([]entity.ContentEntity, error) {
	var modelContents []model.Content

	err := c.db.Order("created_at desc").Preload("User","Category").Find(&modelContents).Error
	if err != nil {
		code = "[REPOSITORY] GetContents - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resps := []entity.ContentEntity{}
	for _, val := range modelContents {
		tags := strings.Split(val.Tags, ",")
		resp := entity.ContentEntity{
			ID : val.ID,
			Title:val.Title,
			Excerpt:val.Excerpt,
			Description:val.Description,
			Image:val.Image,
			Tags:tags,
			Status:val.Status,
			CategoryID:val.CategoryID,
			CreatedByID:val.CreatedByID,
			CreatedAt:val.CreatedAt,
			Category:entity.CategoryEntity{
				ID:val.Category.ID,
				Title:val.Category.Title,
				Slug:val.Category.Slug,
			},
			User:entity.UserEntity{
				ID:val.User.ID,
				Name:val.User.Name,
			},
		}

		resps	= append(resps, resp)
	}

	return resps, nil
}

func (c *contentRepository) CreateContent(ctx context.Context, req entity.ContentEntity) error {	
	panic("implement me")
}

func (c *contentRepository) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	panic("implement me")
}

func (c *contentRepository) UpdateContent(ctx context.Context, req *entity.ContentEntity) error {
	panic("implement me")
}		

func (c *contentRepository) DeleteContent(ctx context.Context, id int64) error {
	panic("implement me")
}

func NewContentRepository(db *gorm.DB) *contentRepository {
	return &contentRepository{db: db}
}