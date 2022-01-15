package usecases

import (
	"api/modules/users/models"
	"time"
)

type TokenManager interface {
	CreateToken(userId int64, levelAccess models.LevelAccess, duration time.Duration) (string, error)
	VerifyToken(token string) (*models.TokenPayload, error)
}

type TokenManagerJwt struct {
	tokenMaker models.TokenMaker
}

func NewTokenManager(tokenMaker models.TokenMaker) TokenManager {
	return &TokenManagerJwt{tokenMaker}
}

func (manager *TokenManagerJwt) CreateToken(userId int64, levelAccess models.LevelAccess, duration time.Duration) (string, error) {
	payload := models.NewTokenPayload(userId, levelAccess, duration)
	token, err := manager.tokenMaker.CreateToken(payload)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (manager *TokenManagerJwt) VerifyToken(token string) (*models.TokenPayload, error) {
	payload, err := manager.tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
