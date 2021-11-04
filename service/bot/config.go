package bot

import "github.com/gofrs/uuid"

type Config struct {
	AccessToken string
	BotID       uuid.UUID
	Prefix      string
	Origin      string
}
