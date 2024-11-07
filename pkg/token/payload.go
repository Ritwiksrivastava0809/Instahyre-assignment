package token

import (
	"fmt"
	"spam-search/pkg/constants"
	"time"

	"github.com/google/uuid"
)

func NewPayload(userID uint, duration time.Duration) (*Payload, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	Payload := &Payload{
		ID:        token,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return Payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return fmt.Errorf(constants.ExipredToken)
	}

	return nil
}
