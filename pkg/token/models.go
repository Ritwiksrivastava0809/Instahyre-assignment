package token

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(phoneNumber string, userID uint, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	ID          uuid.UUID `json:"id"`
	UserID      uint      `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

type JWTMAKER struct {
	secretKey string
}
