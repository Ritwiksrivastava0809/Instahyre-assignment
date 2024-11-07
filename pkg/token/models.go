package token

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(phoneNumber string, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	ID          uuid.UUID `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

type JWTMAKER struct {
	secretKey string
}
