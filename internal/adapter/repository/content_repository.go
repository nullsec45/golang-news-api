package repository

import (
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	err := c.db.Order("created_at desc").Preload(clause.Associations).Find(&modelContents).Error
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
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:req.Title,
		Excerpt:req.Excerpt,
		Description:req.Description,
		Image:req.Image,
		Tags:tags,
		Status:req.Status,
		CategoryID:req.CategoryID,
		CreatedByID:req.CreatedByID,
	}

	err = c.db.Create(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] CreateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *contentRepository) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	var modelContent model.Content

	err = c.db.Where("id = ?", id).Preload(clause.Associations).First(&modelContent).Error

	if err != nil {
		code = "[REPOSITORY] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	tags := strings.Split(modelContent.Tags, ",")
	resp := entity.ContentEntity{
		ID : modelContent.ID,
		Title:modelContent.Title,
		Excerpt:modelContent.Excerpt,
		Description:modelContent.Description,
		Image:modelContent.Image,
		Tags:tags,
		Status:modelContent.Status,
		CategoryID:modelContent.CategoryID,
		CreatedByID:modelContent.CreatedByID,
		CreatedAt:modelContent.CreatedAt,
		Category:entity.CategoryEntity{
			ID:modelContent.Category.ID,
			Title:modelContent.Category.Title,
			Slug:modelContent.Category.Slug,
		},
		User:entity.UserEntity{
			ID:modelContent.User.ID,
			Name:modelContent.User.Name,
		},
	}

	return &resp,nil
}

func (c *contentRepository) UpdateContent(ctx context.Context, req *entity.ContentEntity) error {
		tags := strings.Join(req.Tags, ",")
		modelContent := model.Content{
			Title:req.Title,
			Excerpt:req.Excerpt,
			Description:req.Description,
			Image:req.Image,
			Tags:tags,
			Status:req.Status,
			CategoryID:req.CategoryID,
			CreatedByID:req.CreatedByID,
		}

		err = c.db.Where("id = ?", req.ID).Updates(modelContent).Error

		if err != nil {
			code = "[REPOSITORY] UpdateContent - 1"
			log.Errorw(code, err)
			return err
		}
		
		return nil
}		

func (c *contentRepository) DeleteContent(ctx context.Context, id int64) error {
	err = c.db.Where("id= ?", id).Delete(&model.Content{}).Error
	if err != nil {
		code = "[REPOSITORY] DeleteContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewContentRepository(db *gorm.DB) *contentRepository {
	return &contentRepository{db: db}
}