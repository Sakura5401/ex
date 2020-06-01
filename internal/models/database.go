package models

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

func Connection(postgresURL string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", postgresURL)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&BankDetail{})
	db.AutoMigrate(&BankCard{})
	db.AutoMigrate(&LegalEntityInfo{})
	db.AutoMigrate(&IndividualInfo{})

	return db, err
}

func CreateTableIfNotExists(db *gorm.DB) error {
	if !db.HasTable(&BankDetail{}) {
		err := db.AutoMigrate(&BankDetail{}).Error
		if err != nil {
			return err
		}
	}

	if !db.HasTable(&BankCard{}) {
		err := db.AutoMigrate(&BankCard{}).Error
		if err != nil {
			return err
		}
	}

	if !db.HasTable(&LegalEntityInfo{}) {
		err := db.AutoMigrate(&LegalEntityInfo{}).Error
		if err != nil {
			return err
		}
	}

	if !db.HasTable(&IndividualInfo{}) {
		err := db.AutoMigrate(&IndividualInfo{}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func ModifyColumnIfTableExist(db *gorm.DB) error {
	if db.HasTable(&BankCard{}) {
		err := db.Model(&BankCard{}).DropColumn("card_number").Error
		if err != nil {
			log.Info().Msg("BankCard doesn't have a column card_number")
		}
	}

	if db.HasTable(&BankCard{}) {
		err := db.Model(&BankCard{}).DropColumn("external_id").Error
		if err != nil {
			log.Info().Msg("BankCard doesn't have a column external_id")
		}
	}

	if db.HasTable(&BankCard{}) {
		err := db.Model(&BankCard{}).DropColumn("expiration_date").Error
		if err != nil {
			log.Info().Msg("BankCard doesn't have a column expiration_date")
		}
	}

	if db.HasTable(&IndividualInfo{}) {
		err := db.Model(&IndividualInfo{}).DropColumn("login").Error
		if err != nil {
			log.Info().Msg("BankCard doesn't have a column login")
		}
	}

	if db.HasTable(&IndividualInfo{}) {
		err := db.Model(&IndividualInfo{}).DropColumn("pass").Error
		if err != nil {
			log.Info().Msg("BankCard doesn't have a column pass")
		}
	}

	return nil
}

func DropTablesIfExists(db *gorm.DB) error {
	if db.HasTable(&BankDetail{}) {
		err := db.DropTable(&BankDetail{}).Error
		if err != nil {
			return err
		}
	}

	if db.HasTable(&BankCard{}) {
		err := db.DropTable(&BankCard{}).Error
		if err != nil {
			return err
		}
	}

	if db.HasTable(&LegalEntityInfo{}) {
		err := db.DropTable(&LegalEntityInfo{}).Error
		if err != nil {
			return err
		}
	}

	if db.HasTable(&IndividualInfo{}) {
		err := db.DropTable(&IndividualInfo{}).Error
		if err != nil {
			return err
		}
	}

	return nil
}
