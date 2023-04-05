package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"

	//"github.com/golang/protobuf/ptypes/duration"
	"time"
)

const minSecretKeySize = 32

//JW maker is a JSON web token maker
type JWTMaker struct {
	secretKey string
}

//new JWT maker creates a new  JWT Maker
func NewJWTMaker(secret string) (Maker, error) {
	if len(secret) < minSecretKeySize {
		return nil, fmt.Errorf("the secreyt key must of %d length ", minSecretKeySize)
	}
	return &JWTMaker{
		secretKey: secret,
	}, nil
}

//createToken creats a new token for a username for a specifc time duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, e := jwtToken.SignedString([]byte(maker.secretKey))
	if e != nil {

	}
	return t, e

}

//verifyToken chcks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {

			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {

			return nil, ErrInvalidToken
		}

		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {

		return nil, ErrInvalidToken

	}

	return payload, nil

}
