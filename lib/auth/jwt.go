package auth

import (
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nullsec45/golang-news-api/config"
	"github.com/nullsec45/golang-news-api/internal/core/domain/entity"
)

type Jwt interface {
	GenerateToken(data *entity.JwtData) (string, int64, error)
	VerifyAccessToken(token string)(*entity.JwtData, error)
}

type Options struct {
	signingKey string
	issuer string
}

func (o *Options) GenerateToken(data *entity.JwtData) (string,int64,error) {
	now := time.Now().Local()
	expiresAt := now.Add(time.Hour * 24)
	data.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	data.RegisteredClaims.Issuer = o.issuer
	data.RegisteredClaims.NotBefore=jwt.NewNumericDate(now)
	acToken := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	accessToken, err := acToken.SignedString([]byte(o.signingKey))
	if err != nil {
		return "", 0, err
	}

	return accessToken, expiresAt.Unix(), nil
}

func (o *Options) VerifyAccessToken(token string) (*entity.JwtData, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token)(interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		}
		return []byte(o.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		claim, ok := parsedToken.Claims.(jwt.MapClaims)

		if!ok || !parsedToken.Valid {
			return nil, err
		}

		jwtData := &entity.JwtData{
			UserID: claim["user_id"].(float64),
		}

		return jwtData, nil
	}

	return nil, fmt.Errorf("Token is not valid")
}

func NewJwt(cfg *config.Config) Jwt {
	// return &jwtImpl{
	// 	options:Options{
	// 		signingKey:cfg.App.JwtSecretKey,
	// 		issuer: cfg.App.JwtIssuer,
	// 	}
	// }
	opt := new(Options)
	opt.signingKey=cfg.App.JwtSecretKey
	opt.issuer=cfg.App.JwtIssuer

	return opt
}