package token

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"time"
)

type PastoMaker struct {
	pasto       *paseto.V2
	symetricKey []byte
}

func CreatePasetoMakerInstance(symetricKey string) (Maker, error) {
	if len(symetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Error the the size must be greater than : %d", chacha20poly1305.KeySize)
	}

	maker := &PastoMaker{
		pasto:       paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}
	return maker, nil

}

//createToken creats a new token for a username for a specifc time duration
func (maker *PastoMaker) CreateToken(username string, duration time.Duration) (string, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {

	}
	return maker.pasto.Encrypt(maker.symetricKey, payload, nil)
}

//verifyToken chcks if the token is valid or not
func (maker *PastoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.pasto.Decrypt(token, maker.symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err

	}
	return payload, nil

}
