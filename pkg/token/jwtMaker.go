package token

import (
	"fmt"
	"time"

	"spam-search/pkg/constants"
	errorlogs "spam-search/pkg/constants/errorlogs"

	"github.com/dgrijalva/jwt-go"
)

func NewJWTMAKER(secretKey string) (Maker, error) {
	if len(secretKey) < constants.MinSecretKeyLen {
		return nil, fmt.Errorf(errorlogs.InvalidKeySize, constants.MinSecretKeyLen)
	}
	return &JWTMAKER{secretKey}, nil
}

func (maker *JWTMAKER) CreateToken(phoneNumber string, userID uint, duration time.Duration) (string, error) {
	payload, err := NewPayload(phoneNumber, userID, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMAKER) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf(constants.InvalidToken)
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && verr.Errors == constants.JWTValidationErrorExpired {
			return nil, fmt.Errorf(constants.ExipredToken)
		}
		return nil, fmt.Errorf(constants.InvalidToken)
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf(constants.InvalidToken)
	}

	return payload, nil
}
