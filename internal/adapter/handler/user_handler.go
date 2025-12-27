package handler

import (
	"github.com/nullsec45/golang-news-api/internal/adapter/handler/request"
	"github.com/nullsec45/golang-news-api/internal/adapter/handler/response"
	"github.com/nullsec45/golang-news-api/internal/core/service"
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	validatorLib "github.com/nullsec45/golang-news-api/lib/validator"
)

type UserHandler interface {
	GetUserByID(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
	RegisterUser(c *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService:userService}
}

func (u *userHandler) GetUserByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] GetUserByID - 1"
		log.Errorw(code,err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized Access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	user, err := u.userService.GetUserByID(c.Context(), int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] GetUserByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}
	
	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Meta.Message = "Success Get User By ID"
	resp := response.UserResponse{
		ID:user.ID,
		Name:user.Name,
		Email:user.Email,
	}
	defaultSuccessResponse.Data=resp

	return c.JSON(defaultSuccessResponse)
}


func (u *userHandler) UpdatePassword(c *fiber.Ctx)  error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code := "[HANDLER] UpdatePassword - 1"
		log.Errorw(code,err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="Unauthorized Access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.UpdatePasswordRequest
	if err = c.BodyParser(&req); err != nil {
		code := "[HANDLER] UpdatePassword - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message="invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code := "[HANDLER] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqUpdatePassword := entity.UpdatePasswordEntity {
		CurrentPassword:req.CurrentPassword,
		NewPassword:req.NewPassword,
		ConfirmPassword:req.ConfirmPassword,
	}

	err = u.userService.UpdatePassword(c.Context(), reqUpdatePassword, int64(claims.UserID))
	if err != nil {
		code := "[HANDLER] UpdatePassword - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()
		
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status=true
	defaultSuccessResponse.Meta.Message="Update password Successfully"
	defaultSuccessResponse.Data=nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (u *userHandler) RegisterUser(c *fiber.Ctx) error {
	req := request.RegisterRequest{}
	resp := response.Meta{}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] RegisterUser - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if req.Role == "" {
        req.Role = "viewer"
    }	

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] RegisterUser - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqRegister := entity.RegisterUserEntity {
		Name:req.Name,
		Email:req.Email,
		Role:req.Role,
		Password:req.Password,
	}

	err := u.userService.RegisterUser(c.Context(), reqRegister)
	
	if err != nil {
		code = "[HANDLER] RegisterUser - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status=false
		errorResp.Meta.Message=err.Error()

		if err.Error() == "Email already exists" {
			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	resp.Status=true
	resp.Message="Register successfull"

	return c.JSON(resp)

}
