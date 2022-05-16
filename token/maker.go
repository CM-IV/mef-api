package token

import "time"

//Interface for managing Paseto Tokens
type Maker interface {
	//Creates new token for a specific username & duration
	CreateToken(username string, duration time.Duration) (string, error)

	//Checks if token is valid or not
	VerifyToken(token string) (*Payload, error)
}
