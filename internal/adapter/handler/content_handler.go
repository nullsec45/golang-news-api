package handler

import (
	// "github.com/nullsec45/golang-news-api/internal/adapter/handler/request"
	"github.com/nullsec45/golang-news-api/internal/adapter/handler/response"
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"github.com/nullsec45/golang-news-api/internal/core/service"
	// "github.com/nullsec45/golang-news-api/lib/conv"
	// validatorLib "github.com/nullsec45/golang-news-api/lib/validator"
	// "fmt"
	// "os"
	// "strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentByID(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	UpdateContent(c *fiber.Ctx) error
	DeleteContent(c *fiber.Ctx) error
	UploadImageR2(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{contentService: contentService}
} 

func (coh *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetContents - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := coh.contentService.GetContents(c.Context())
	if err != nil {
		code = "[HANDLER] GetContents  - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Meta.Message="Contents fetched Successfully"

	respContents := []response.ContentResponse{}
	for _, content := range results {
		respContent := response.ContentResponse {
			ID : content.ID,
			Title: content.Title,
			Excerpt: content.Excerpt,
			Description: content.Description,
			Image:content.Image,
			Tags: content.Tags,
			Status: content.Status,
			CategoryID: content.CategoryID,
			CreatedByID: content.CreatedByID,
			CreatedAt: content.CreatedAt.Format(time.RFC3339),
			CategoryName:content.Category.Title,
			Author: content.User.Name,
		}

		respContents = append(respContents, respContent)
	}

	defaultSuccessResponse.Data=results
	defaultSuccessResponse.Pagination=nil
	return c.JSON(defaultSuccessResponse)
}

func (coh *contentHandler) GetContentByID(c *fiber.Ctx) error {
	panic("implement me")
}

func (coh *contentHandler) CreateContent(c *fiber.Ctx) error {
	panic("implement me")
}

func (coh *contentHandler) UpdateContent(c *fiber.Ctx) error {
	panic("implement me")
}

func (coh *contentHandler) DeleteContent(c *fiber.Ctx) error {
	panic("implement me")
}

func (coh *contentHandler) UploadImageR2(c *fiber.Ctx) error {
	panic("implement me")
}