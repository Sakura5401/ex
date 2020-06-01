package models

import (
	"testing"
)

func TestCreateBankDetailFail(t *testing.T) {
	cases := []struct {
		name       string
		bankDetail BankDetail
	}{
		{"BankIdentificationCode10", BankDetail{
			BankIdentificationCode: "12345678910",
		},
		},
		{"LettersBankIdentificationCode", BankDetail{
			BankIdentificationCode: "12345678a",
		},
		},
		{"SpacesBankIdentificationCode", BankDetail{
			BankIdentificationCode: "         ",
		},
		},
		{"CorrAccount19", BankDetail{
			CorrAccount: "1234567891234567891",
		},
		},
		{"LettersCorrAccount", BankDetail{
			CorrAccount: "123456789123456789aa",
		},
		},
		{"SpacesCorrAccount", BankDetail{
			CorrAccount: " 123456789123456789 ",
		},
		},
		{"LettersCheckingAccount", BankDetail{
			CheckingAccount: "123456789123456789aa",
		},
		},
		{"SpacesCheckingAccount", BankDetail{
			CheckingAccount: " 123456789123456789 ",
		},
		},
		{"DetailName#", BankDetail{
			DetailName: "aaa#",
		},
		},
		{"DetailName_ ", BankDetail{
			DetailName: "aaa_ ",
		},
		},
		{"onlySpacesDetailName", BankDetail{
			DetailName: "    ",
		},
		},
		{"onlyNumbersDetailName", BankDetail{
			DetailName: "123",
		},
		},
		{"BankName:", BankDetail{
			BankName: "aaa: ",
		},
		},

		{"BankName!", BankDetail{
			BankName: "aaa!",
		},
		},
		{"onlySpacesBankName", BankDetail{
			BankName: "    ",
		},
		},
		{"onlyNumbersBankName", BankDetail{
			BankName: "123",
		},
		},
		{"City", BankDetail{
			City: "aaa%",
		},
		},
		{"City&", BankDetail{
			City: "aaa&",
		},
		},
		{"onlySpacesCity", BankDetail{
			City: "    ",
		},
		},
		{"onlyNumbersCity", BankDetail{
			City: "123",
		},
		},
		{"onlySpacesAddress", BankDetail{
			Address: "    ",
		},
		},
		{"onlyNumbersAddress", BankDetail{
			Address: "123",
		},
		},
		{"Address`", BankDetail{
			Address: "`aaa`",
		},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dropDB := setUp()
			defer dropDB()
			errorCreate := db.Create(&tc.bankDetail).Error //nolint
			if errorCreate == nil {
				t.Error(errorCreate)
			}
		})
	}
}

func TestCreateBankDetailSuccess(t *testing.T) {
	cases := []struct {
		name       string
		bankDetail BankDetail
	}{
		{"BankIdentificationCode", BankDetail{
			BankIdentificationCode: "123456789",
			CorrAccount:            "12345678912345678912",
			CheckingAccount:        "12345678912345678912",
		},
		},
		{"CorrAccount", BankDetail{
			BankIdentificationCode: "123456789",
			CorrAccount:            "12345678912345678912",
			CheckingAccount:        "12345678912345678912",
		},
		},
		{"CheckingAccount", BankDetail{
			BankIdentificationCode: "123456789",
			CorrAccount:            "12345678912345678912",
			CheckingAccount:        "12345678912345678912",
		},
		},
		{"DetailName", BankDetail{
			BankIdentificationCode: "123456789",
			CorrAccount:            "12345678912345678912",
			CheckingAccount:        "12345678912345678912",
			DetailName:             "Реквизиты 1",
		},
		},
		{"BankName", BankDetail{
			BankIdentificationCode: "123456789",
			CorrAccount:            "12345678912345678912",
			CheckingAccount:        "12345678912345678912",
			BankName:               "Акрополь",
		},
		},
		{"City", BankDetail{
			BankIdentificationCode: "123456789",
			CorrAccount:            "12345678912345678912",
			CheckingAccount:        "12345678912345678912",
			City:                   "г Санкт-Петербург",
		},
		},
		{"Address", BankDetail{
			BankIdentificationCode: "123456789",
			CorrAccount:            "12345678912345678912",
			CheckingAccount:        "12345678912345678912",
			Address:                "109428, ГОРОД МОСКВА, ПРОСПЕКТ РЯЗАНСКИЙ, ДОМ 22, КОРПУС 2, ИНВ. №IX, КОМ.19.20.34-40",
		},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dropDB := setUp()
			defer dropDB()
			errorCreate := db.Create(&tc.bankDetail).Error //nolint
			if errorCreate != nil {
				t.Error(errorCreate)
			}
		})
	}
}
