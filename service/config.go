package service

import "github.com/gofrs/uuid"

type Config struct {
	VerificationToken string
	AccessToken       string
	BotID             uuid.UUID
	Prefix            string
	Origin            string
}
