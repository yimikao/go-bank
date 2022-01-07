package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretkey string
}

const minSecretKeyLength = 32

func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeyLength {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeyLength)
	}
	return &JWTMaker{secretkey: secretkey}, nil
}

func (m *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	claims, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //raw token: *jwtToken
	return jwtToken.SignedString([]byte(m.secretkey))             //token that will later be parsed

}

func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {

	//closure to verify token header to make sure that the signing algorithm matches with  one normally use to sign tokens
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC) //i used HS256, which is an instance of the SigningMethodHMAC struct
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretkey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
