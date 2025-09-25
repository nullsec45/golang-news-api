package service

import (
	"github.com/nullsec45/golang-news-api/internal/adapter/repository"
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"github.com/nullsec45/golang-news-api/lib/conv"
	"context"
	"github.com/gofiber/fiber/v2/log"
)

type CategoryService interface {
	GetCategories(ctx context.Context)([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64)(*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func (c *categoryService) GetCategories (ctx context.Context)([]entity.CategoryEntity, error) {
	results, err := c.categoryRepository.GetCategories(ctx)
	if err != nil {
		code = "[SERVICE] GetCategories - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return results, nil
}

func (c *categoryService) GetCategoryByID(ctx context.Context, id int64)(*entity.CategoryEntity, error) {
	panic("kiw")
}

func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	slug := conv.GenerateSlug(req.Title)
	req.Slug=slug

	err = c.categoryRepository.CreateCategory(ctx, req)
	if err != nil {
		code="[SERVICE] CreateCategory - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *categoryService) EditCategoryByID(ctx context.Context, req entity.CategoryEntity) error {
	panic("kiw")
}

func (c *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	panic("kiw")
}

func NewCategoryService(categoryRepo repository.CategoryRepository)CategoryService{
	return &categoryService{categoryRepository:categoryRepo}
}