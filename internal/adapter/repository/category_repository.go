package repository

import (
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"context"
	"github.com/nullsec45/golang-news-api/internal/core/domain/model"
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"fmt"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context)([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64)(*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *gorm.DB
}

func (c *categoryRepository) GetCategories(ctx context.Context)([]entity.CategoryEntity, error) {
	var modelCategories []model.Category

	err = c.db.Order("created_at DESC").Preload("User").Find(&modelCategories).Error
	if err != nil {
		code = "[REPOSITORY] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if len(modelCategories) == 0 {
		code="[REPOSITORY] GetCategories - 2"
		err=errors.New("data not found")
		log.Errorw(code, err)
		return nil, err
	}

	var resps []entity.CategoryEntity
	for _, val := range modelCategories {
		resps = append(resps, entity.CategoryEntity{
			ID: val.ID,
			Title:val.Title,
			Slug:val.Slug,
			User:entity.UserEntity{
				ID:val.User.ID,
				Name:val.User.Name,
				Email:val.User.Email,
				Password:val.User.Password,
			},
		})
	}

	return resps, nil
}

func (c *categoryRepository) GetCategoryByID(ctx context.Context, id int64)(*entity.CategoryEntity, error) {
	panic("kiw")
}

func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	var countSlug int64
	err = c.db.Table("categories").Where("slug=?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code = "[REPOSITORY] CreateCategory - 2"
		log.Errorw(code,err)
		return err
	}
	countSlug=countSlug + 1
	slug := fmt.Sprintf("%s-%d", req.Slug, countSlug)
	modelCategory := model.Category{
		Title:req.Title,
		Slug:slug,
		CreatedByID:req.User.ID,
	}

	err=c.db.Create(&modelCategory).Error
	if err != nil {
		code="[REPOSITORY] CreateCategory - 1"
		log.Errorw(code,err)
		return err
	}

	return nil
}

func (c *categoryRepository) EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error {
	panic("kiw")
}

func (c *categoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	panic("kiw")
}


func NewCategoryRepository(db *gorm.DB) CategoryRepository{
	return &categoryRepository{
		db:db,
	}
}