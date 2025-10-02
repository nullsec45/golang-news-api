package handler

import (
	"github.com/nullsec45/golang-news-api/internal/adapter/handler/request"
	"github.com/nullsec45/golang-news-api/internal/adapter/handler/response"
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"github.com/nullsec45/golang-news-api/internal/core/service"
	"github.com/nullsec45/golang-news-api/lib/conv"
	validatorLib "github.com/nullsec45/golang-news-api/lib/validator"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"fmt"
	"os"
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
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetContentByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("id")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()
	
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := coh.contentService.GetContentByID(c.Context(), contentID)

	if err != nil {
		code = "[HANDLER] GetContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()
	
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Meta.Message="Contents fetched Successfully"

	respContent := response.ContentResponse {
		ID : result.ID,
		Title: result.Title,
		Excerpt: result.Excerpt,
		Description: result.Description,
		Image:result.Image,
		Tags: result.Tags,
		Status: result.Status,
		CategoryID: result.CategoryID,
		CreatedByID: result.CreatedByID,
		CreatedAt: result.CreatedAt.Local().String(),
		CategoryName:result.Category.Title,
		Author: result.User.Name,
	}


	defaultSuccessResponse.Data=respContent
	defaultSuccessResponse.Pagination=nil
	return c.JSON(defaultSuccessResponse)
}

func (coh *contentHandler) CreateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] CreateContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.ContentRequest
	if err = c.BodyParser(&req); err != nil {
		code := "[HANDLER] CreateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code := "[HANDLER] CreateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:req.Title,
		Excerpt:req.Excerpt,
		Description:req.Description,
		Image:req.Image,
		Tags:tags,
		Status:req.Status,
		CategoryID:req.CategoryID,
		CreatedByID:int64(userID),
	}

	err = coh.contentService.CreateContent(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] CreateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()
		
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Meta.Message="Content created Successfully"
	defaultSuccessResponse.Data=nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (coh *contentHandler) UpdateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] UpdateContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.ContentRequest
	if err = c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code := "[HANDLER] UpdateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	idParam := c.Params("id")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] UpdateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()
	
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		ID:contentID,
		Title:req.Title,
		Excerpt:req.Excerpt,
		Description:req.Description,
		Image:req.Image,
		Tags:tags,
		Status:req.Status,
		CategoryID:req.CategoryID,
		CreatedByID:int64(userID),
	}

	err = coh.contentService.UpdateContent(c.Context(), &reqEntity)
	if err != nil {
		code := "[HANDLER] UpdateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()
		
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Meta.Message="Content updated Successfully"
	defaultSuccessResponse.Data=nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (coh *contentHandler) DeleteContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] DeleteContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("id")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] DeleteContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()
	
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = coh.contentService.DeleteContent(c.Context(), contentID)
	if err != nil {
		code = "[HANDLER] DeleteContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()
	
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Meta.Message="Content deleted Successfully"
	defaultSuccessResponse.Data=nil

	return c.JSON(defaultSuccessResponse)

}

func (coh *contentHandler) UploadImageR2(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] UploadImageR2 - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.FileUploadRequest
	file, err := c.FormFile("image")
	if err != nil {
		code := "[HANDLER] UploadImageR2 - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Invalid request body"		
	
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = c.SaveFile(file, fmt.Sprintf("./temp/content/%s", file.Filename)); err != nil {
		code := "[HANDLER] UploadImageR2 - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	req.Image = fmt.Sprintf("./temp/content/%s", file.Filename)
	reqEntity := entity.FileUploadEntity{
		Name:fmt.Sprintf("%d-%d", int64(claims.UserID), time.Now().UnixNano()),
		Path:req.Image,
	}

	urlImage, err := coh.contentService.UploadImageR2(c.Context(), reqEntity)
	if err != nil {
		code := "[HANDLER] UploadImageR2 - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message= err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if req.Image != "" {
		err = os.Remove(req.Image)
		if err != nil {
			code := "[HANDLER] UploadImageR2 - 5"
			log.Errorw(code, err)
			errorResp.Meta.Status=false
			errorResp.Meta.Message= err.Error()
	
			return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
		}
	}


	urlImageResp := map[string]interface{}{
		"urlimage":urlImage,
	}

	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Meta.Message="Success Upload Image"
	defaultSuccessResponse.Data=urlImageResp

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}