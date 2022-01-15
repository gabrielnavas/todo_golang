package models

import (
	"errors"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

const DurationTimeDefault = time.Hour * 24 * 7

type TokenMaker interface {
	CreateToken(payload *TokenPayload) (string, error)
	VerifyToken(token string) (*TokenPayload, error)
}

type TokenPayload struct {
	UserId      int64
	LevelAccess LevelAccess
	IssuedAt    time.Time
	ExpiredAt   time.Time
}

func NewTokenPayload(userId int64, levelAccess LevelAccess, duration time.Duration) *TokenPayload {
	payload := &TokenPayload{
		userId,
		levelAccess,
		time.Now(),
		time.Now().Add(duration),
	}

	return payload
}

func (payload *TokenPayload) Valid() error {
	// TODO: precisa do time.Now(), não poderia ser time.After(...) ??
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
