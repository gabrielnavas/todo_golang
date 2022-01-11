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
	CreateToken(userId int64, duration time.Duration) (string, error)
	VerifyToken(token string) (*TokenPayload, error)
}

type TokenPayload struct {
	ID        int64     `json:"id"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func NewTokenPayload(userId int64, duration time.Duration) *TokenPayload {
	payload := &TokenPayload{
		userId,
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
