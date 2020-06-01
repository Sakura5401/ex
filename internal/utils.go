package internal

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"gitlab.safecrow.ru/safecrow/gateway-requisites/v2/internal/models"
)

// IsExistDefaultDetail получить банковские реквизиты, которые назначены основным платёжным методом, если существуют.
// Возвращает либо банковские реквизиты, либо nil.
func IsExistDefaultDetail(db *gorm.DB, userID uuid.UUID) (*models.BankDetail, error) {
	var detail models.BankDetail

	err := db.Where("user_id = ? and default_detail is true", userID).First(&detail).Error

	if gorm.IsRecordNotFoundError(err) {
		log.Info().Msg("Default bank details is not exist")
		return nil, nil
	}

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &detail, nil
}

// IsExistDefaultCard получить банковскую карту, которая назначена основным платёжным методом, если существует.
// Возвращает либо банковскую карту, либо nil.
func IsExistDefaultCard(db *gorm.DB, userID uuid.UUID) (*models.BankCard, error) {
	var card models.BankCard

	err := db.Where("user_id = ? and default_card is true", userID).First(&card).Error

	if gorm.IsRecordNotFoundError(err) {
		log.Info().Msg("Default card is not exist")
		return nil, nil
	}

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &card, nil
}

// IsExistBankDetail получить банковские реквизиты, если существуют.
// Возвращает либо банковские реквизиты, либо nil.
func IsExistBankDetail(db *gorm.DB, infoID uuid.UUID) (*models.BankDetail, error) {
	var detail models.BankDetail

	err := db.Where("UUID = ?", infoID).First(&detail).Error

	if gorm.IsRecordNotFoundError(err) {
		log.Info().Msg("Bank details is not exist")
		return nil, nil
	}

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &detail, nil
}

// IsExistBankCard получить банковскую карту, если существует.
// Возвращает либо банковскую карту, либо nil.
func IsExistBankCard(db *gorm.DB, infoID uuid.UUID) (*models.BankCard, error) {
	var card models.BankCard

	err := db.Where("UUID = ?", infoID).First(&card).Error

	if gorm.IsRecordNotFoundError(err) {
		log.Info().Msg("Bank card is not exist")
		return nil, nil
	}

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &card, nil
}
