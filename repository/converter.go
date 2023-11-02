package repository

import (
	"github.com/gofrs/uuid"

	"git.trap.jp/toki/bot_converter/model"
)

type ConverterRepository interface {
	// CreateConverter Converterを作成します
	//
	// secretに空文字列を指定した場合は、secretを持たないconverterが作成されます
	CreateConverter(creatorID, channelID uuid.UUID, secret string) (*model.Converter, error)
	// GetConverter Converterを取得します
	GetConverter(id uuid.UUID) (*model.Converter, error)
	// GetConverterConfig Converter Configを取得します
	GetConverterConfig(id uuid.UUID) (*model.Config, error)
	// GetConverterByCreatorID 指定されたユーザーによって作られたconverter全てを取得します。
	GetConverterByCreatorID(creatorID uuid.UUID) ([]*model.Converter, error)
	// DeleteConverter Converterを削除します
	DeleteConverter(id uuid.UUID) error
}
