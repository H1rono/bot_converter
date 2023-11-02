package repository

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"git.trap.jp/toki/bot_converter/model"
)

func (repo *GormRepository) CreateConverter(creatorID, channelID uuid.UUID, secret string) (*model.Converter, error) {
	if creatorID == uuid.Nil || channelID == uuid.Nil {
		return nil, ErrNilID
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	c := &model.Converter{
		ID:        id,
		CreatorID: creatorID,
		ChannelID: channelID,
	}
	if len(secret) > 0 {
		c.Secret = sql.NullString{
			String: secret,
			Valid:  true,
		}
	}

	if err := repo.db.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (repo *GormRepository) GetConverter(id uuid.UUID) (*model.Converter, error) {
	if id == uuid.Nil {
		return nil, ErrNilID
	}

	var c model.Converter
	if err := repo.db.Where(&model.Converter{ID: id}).First(&c).Error; err != nil {
		return nil, convertError(err)
	}
	return &c, nil
}

func (repo *GormRepository) GetConverterConfig(id uuid.UUID) (*model.Config, error) {
	if id == uuid.Nil {
		return nil, ErrNilID
	}

	var c model.Config
	if err := repo.db.Where(&model.Config{ConverterID: id}).First(&c).Error; err != nil {
		return nil, convertError(err)
	}
	return &c, nil
}

func (repo *GormRepository) GetConverterByCreatorID(creatorID uuid.UUID) ([]*model.Converter, error) {
	if creatorID == uuid.Nil {
		return nil, ErrNilID
	}

	var cs []*model.Converter
	if err := repo.db.Where(&model.Converter{CreatorID: creatorID}).Find(&cs).Error; err != nil {
		return nil, convertError(err)
	}
	return cs, nil
}

func (repo *GormRepository) DeleteConverter(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrNilID
	}

	return repo.db.Delete(&model.Converter{ID: id}).Error
}
