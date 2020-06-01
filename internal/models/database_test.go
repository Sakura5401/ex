package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.safecrow.ru/safecrow/gateway-requisites/v2/utils"

	_ "github.com/lib/pq"
)

var db = func() *gorm.DB { //nolint
	settings := utils.GetEnvs()
	db, err := gorm.Open("postgres", settings.PostgresURL)
	if err != nil {
		panic(err)
	}

	return db
}()

func setUp() func() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "test_" + defaultTableName
	}

	db.CreateTable(&BankDetail{})

	return func() {
		db.DropTable(&BankDetail{})
	}
}
