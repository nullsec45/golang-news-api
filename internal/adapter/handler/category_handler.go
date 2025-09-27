package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nullsec45/golang-news-api/internal/core/service"
	"github.com/nullsec45/golang-news-api/internal/adapter/handler/request"
	"github.com/nullsec45/golang-news-api/internal/adapter/handler/response"
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"github.com/nullsec45/golang-news-api/lib/conv"
	validatorLib "github.com/nullsec45/golang-news-api/lib/validator"

	"github.com/gofiber/fiber/v2/log"
)

var defaultSuccessResponse response.DefaultSuccessResponse


type CategoryHandler interface {
	GetCategories(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx)  error
	EditCategoryByID(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}


type categoryHandler struct {
	categoryService service.CategoryService
}

func (ch *categoryHandler) GetCategories(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetCategories  - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse {
			ID:result.ID,
			Title:result.Title,
			Slug:result.Slug,
			CreatedByName:result.User.Name,
		}
		categoryResponses=append(categoryResponses, categoryResponse)
	}

	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Data=categoryResponses
	defaultSuccessResponse.Pagination=nil
	defaultSuccessResponse.Meta.Message="Categories fetched Successfully"

	return c.JSON(defaultSuccessResponse)
}

func (ch *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("id")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message= err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.categoryService.GetCategoryByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponse := response.SuccessCategoryResponse{
		ID:result.ID,
		Title:result.Title,
		Slug:result.Slug,
		CreatedByName:result.User.Name,
	}
	
	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Data=categoryResponse
	defaultSuccessResponse.Pagination=nil
	defaultSuccessResponse.Meta.Message="Category fetched Successfully"

	return c.JSON(defaultSuccessResponse)
}

func (ch *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] CreateCategory - 1"
		log.Errorw(code,err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code="[HANDLER] CreateCategory - 2"
		log.Errorw(code,err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err=validatorLib.ValidateStruct(req); err != nil {
		code="[HANDLER] CreateCategory - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		Title:req.Title,
		User:entity.UserEntity{
			ID:int64(userID),
		},
	}

	err = ch.categoryService.CreateCategory(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateCategory - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Data=nil
	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Pagination=nil
	defaultSuccessResponse.Meta.Message="Category created successfully"
	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (ch *categoryHandler) EditCategoryByID(c *fiber.Ctx) error {
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] EditCategory - 1"
		log.Errorw(code,err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code="[HANDLER] EditCategory - 2"
		log.Errorw(code,err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err=validatorLib.ValidateStruct(req); err != nil {
		code="[HANDLER] EditCategory - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	idParam := c.Params("id")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] EditCategory - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message= err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		ID:id,
		Title:req.Title,	
		User:entity.UserEntity{
			ID:int64(userID),
		},
	}
	
	err = ch.categoryService.EditCategoryByID(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] EditCategoryByID - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Data=nil
	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Pagination=nil
	defaultSuccessResponse.Meta.Message="Category updated successfully"
	return c.JSON(defaultSuccessResponse)
}

func (ch *categoryHandler) DeleteCategory(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] DeleteCategory - 1"
		log.Errorw(code,err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("id")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] DeleteCategory - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message= err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = ch.categoryService.DeleteCategory(c.Context(), id)

	if err != nil {
		code = "[HANDLER] DeleteCategory - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Data=nil
	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Pagination=nil
	defaultSuccessResponse.Meta.Message="Category deleted successfully"

	return c.JSON(defaultSuccessResponse)
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService:categoryService}
}