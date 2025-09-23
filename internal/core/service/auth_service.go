package service

import (
	"github.com/nullsec45/golang-news-api/config"
	"github.com/nullsec45/golang-news-api/internal/adapter/repository"
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
	"context"
	"time"
	"github.com/nullsec45/golang-news-api/lib/auth"
	"github.com/nullsec45/golang-news-api/lib/conv"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"errors"
)

var err error
var code string

type AuthService interface {
	GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error)
}

type authService struct {
	authRepository repository.AuthRepository
	cfg *config.Config
	jwtToken auth.Jwt
}

func (a *authService) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error) {
	result, err := a.authRepository.GetUserByEmail(ctx, req)
	if err != nil {
		code = "[SERVICE] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if checkPass := conv.CheckPasswordHash(req.Password, result.Password); !checkPass {
		code = "[SERVICE] GetUserByEmail - 2"
		err = errors.New("invalid password")
		log.Errorw(code, "Invalid Password")
		return nil, err
	}

	jwtData := entity.JwtData {
		UserID:float64(result.ID),
		RegisteredClaims:jwt.RegisteredClaims{
			NotBefore:jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			ID:string(result.ID),
		},
	}

	accessToken, expiresAt, err := a.jwtToken.GenerateToken(&jwtData)
	if err != nil {
		code = "[SERVICE] GetUserByEmail - 3"
		log.Errorw(code, err)
		return nil, err
	}

	resp := entity.AccessToken{
		AccessToken:accessToken,
		ExpiresAt:expiresAt,
	}
	
	return &resp, nil
}

func NewAuthService(authRepository repository.AuthRepository, cfg *config.Config, jwtToken auth.Jwt) AuthService {
	return &authService{
		authRepository:authRepository,
		cfg:cfg,
		jwtToken:jwtToken,
	}
}