package repository

import (
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"github.com/nullsec45/golang-news-api/internal/core/domain/model"
	"context"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"errors"
)

type UserRepository interface {
	UpdatePassword(ctx context.Context, newPass string, id int64) error
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
	GetUserByIDWithPassword(ctx context.Context, id int64) (*entity.UserEntityWithPassword, error)

	RegisterUser(ctx context.Context, req entity.RegisterUserEntity) error
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	var modelUser model.User
	err = u.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		code := "[REPOSITORY] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:    id,
		Name:  modelUser.Name,
		Email: modelUser.Email,
	}, nil
}

func (u *userRepository) GetUserByIDWithPassword(ctx context.Context, id int64) (*entity.UserEntityWithPassword, error) {
	var modelUser model.User
	err := u.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		code := "[REPOSITORY] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.UserEntityWithPassword{
		ID:    id,
		Name:  modelUser.Name,
		Email: modelUser.Email,
		Password: modelUser.Password,
	}, nil
}

func (u *userRepository) UpdatePassword(ctx context.Context, newPass string, id int64) error {
	err = u.db.Model(&model.User{}).Where("id = ?", id).Update("password", newPass).Error
	if err != nil {
		code := "[REPOSITORY] UpdatePassword - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (u *userRepository) RegisterUser(ctx context.Context, req entity.RegisterUserEntity) error {
	var existingUser model.User

	err := u.db.WithContext(ctx).
        Select("id").
        Where("email = ?", req.Email).
        First(&existingUser).Error

	if err == nil {
		err=errors.New("Email already exists")
	}	

	if err != gorm.ErrRecordNotFound {
        log.Errorw("[REPOSITORY] RegisterUser - 2 (Check Email)", err)
        return err
    }

		
	modelUser := model.User{
		Name:req.Name,
		Email:req.Email,
		Role:req.Role,
		Password:req.Password,
	}

	err = u.db.Create(&modelUser).Error
	if err != nil {
		code := "[REPOSITORY] RegisterUser - 3"
		log.Errorw(code, err)
		return err
	}

	if err != gorm.ErrRecordNotFound {
        log.Errorw("[REPOSITORY] RegisterUser - 4 (Register User)", err)
        return err
    }

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
