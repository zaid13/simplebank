package token

import (
	"time"
)

//maker is the interface that will manage the tokens
type Maker interface {

	//createToken creats a new token for a username for a specifc time duration
	CreateToken(username string, duration time.Duration) (string, error)

	//verifyToken chcks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
