package tokenjwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"

	"api/modules/users/models"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (models.TokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(payload *models.TokenPayload) (string, error) {
	signingMethods := jwt.SigningMethodHS256 /*Método de assinatura*/
	jwtToken := jwt.NewWithClaims(signingMethods, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*models.TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC) //Método de assinatura
		if !ok {
			return nil, models.ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &models.TokenPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, models.ErrExpiredToken) {
			return nil, models.ErrExpiredToken
		}
		return nil, models.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*models.TokenPayload)
	if !ok {
		return nil, models.ErrInvalidToken
	}
	return payload, nil
}
