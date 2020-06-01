package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

// BankDetail структура банковких реквизитов
type BankDetail struct {
	UUID                   uuid.UUID `gorm:"type:uuid;primary_key"`
	UserID                 uuid.UUID `gorm:"type:uuid" validate:"required"`
	DetailName             string    // Название реквизитов
	BankIdentificationCode string    `validate:"required"` // БИК
	BankName               string
	City                   string
	Address                string
	CorrAccount            string `validate:"required"` // Корреспонденский счёт
	CheckingAccount        string `validate:"required"` // Расчётный счёт
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DefaultDetail          bool `gorm:"default:false"` // Основной платёжный метод?
	Removed                bool `gorm:"default:false"`
}

// BeforeCreate func
func (bankDetail *BankDetail) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UUID", uuid.Must(uuid.NewV4())) //nolint:golint,errcheck
	return nil
}

// BeforeSave func
func (bankDetail *BankDetail) BeforeSave() error {
	if bankDetail.BankIdentificationCode != "" {
		err := bankDetailValid("bank_identification_code", bankDetail.BankIdentificationCode)
		if err != nil {
			return err
		}
	}

	if bankDetail.CorrAccount != "" {
		err := bankDetailValid("corr_account", bankDetail.CorrAccount)
		if err != nil {
			return err
		}
	}

	if bankDetail.CheckingAccount != "" {
		err := bankDetailValid("checking_account", bankDetail.CheckingAccount)
		if err != nil {
			return err
		}
	}

	if bankDetail.DetailName != "" {
		err := bankDetailValid("detail_name", bankDetail.DetailName)
		if err != nil {
			return err
		}
	}

	if bankDetail.BankName != "" {
		err := bankDetailValid("bank_name", bankDetail.BankName)
		if err != nil {
			return err
		}
	}

	if bankDetail.City != "" {
		err := bankDetailValid("city", bankDetail.City)
		if err != nil {
			return err
		}
	}

	if bankDetail.Address != "" {
		err := bankDetailValid("address", bankDetail.Address)
		if err != nil {
			return err
		}
	}

	return nil
}

func bankDetailValid(fieldName string, value string) error {
	switch fieldName {
	case "bank_identification_code", "corr_account", "checking_account":
		err := bankDetailLengthValid(fieldName, value)

		if err != nil {
			return err
		}

		if notNumbersValid(value) {
			return notNumbersError(fieldName, value)
		}
	case "detail_name", "bank_name", "city", "address":
		if !(lettersValid(value) && lettersSymbolsNumbersValid(value)) {
			return failedValidationError(fieldName, value)
		}
	}

	return nil
}

func bankDetailLengthValid(fieldName string, value string) error {
	lenBIC := 9
	lenCorrAccount := 20
	lenCheckingAccount := 20

	switch fieldName {
	case "bank_identification_code":
		if len(value) != lenBIC {
			return wrongLengthError(fieldName, value)
		}
	case "corr_account":
		if len(value) != lenCorrAccount {
			return wrongLengthError(fieldName, value)
		}
	case "checking_account":
		if len(value) != lenCheckingAccount {
			return wrongLengthError(fieldName, value)
		}
	}

	return nil
}
