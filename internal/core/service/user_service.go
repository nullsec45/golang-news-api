package service

import (
	"github.com/nullsec45/golang-news-api/internal/adapter/repository"
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"github.com/nullsec45/golang-news-api/lib/conv"
	"context"

	"github.com/gofiber/fiber/v2/log"
	"errors"
	// "fmt"
)

type UserService interface {
	UpdatePassword(ctx context.Context, req entity.UpdatePasswordEntity, id int64) error
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	result, err := u.userRepo.GetUserByID(ctx, id)
	if err != nil {
		code := "[SERVICE] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}
	return result, nil
}

// UpdatePassword implements UserService.
func (u *userService) UpdatePassword(ctx context.Context, req entity.UpdatePasswordEntity, id int64) error {
	result, err := u.userRepo.GetUserByIDWithPassword(ctx, id)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 1"
		log.Errorw(code, err)
		return err
	}

	if checkPass := conv.CheckPasswordHash(req.CurrentPassword, result.Password); !checkPass {
		code = "[SERVICE] UpdatePassword - 2"
		err = errors.New("Failed update password, current password invalid.")
		log.Errorw(code, "Invalid Password")
		return err
	}

	// if req.NewPassword != req.ConfirmPassword {
	// 	code = "[SERVICE] UpdatePassword - 3"
	// 	err = errors.New("Failed update password, new password and confirm password don't match!.")
	// 	log.Errorw(code, "Invalid Password")
	// 	return  err
	// }

	password, err := conv.HashPassword(req.ConfirmPassword)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 4"
		log.Errorw(code, err)
		return err
	}

	err = u.userRepo.UpdatePassword(ctx, password, id)
	if err != nil {
		code := "[SERVICE] UpdatePassword - 5"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}
