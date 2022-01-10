package token

import (
	"errors"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        int64     `json:"id"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func NewPayload(userId int64, duration time.Duration) *Payload {
	payload := &Payload{
		userId,
		time.Now(),
		time.Now().Add(duration),
	}

	return payload
}

func (payload *Payload) Valid() error {
	// TODO: precisa do time.Now(), não poderia ser time.After(...) ??
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
