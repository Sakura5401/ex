package models

import (
	"fmt"
	"regexp"
)

func notNumbersValid(value string) bool {
	notNumbersMatch := "[\\D]+"
	notNumbers := regexp.MustCompile(notNumbersMatch)

	return notNumbers.MatchString(value)
}

func lettersValid(value string) bool {
	lettersMatch := "[а-яА-ЯёЁa-zA-Z]+"
	letters := regexp.MustCompile(lettersMatch)

	return letters.MatchString(value)
}

func companyNameValid(value string) bool {
	//companyNameMatch := "^([а-яА-ЯёЁa-zA-Z\\d.,/()?&=:;`+-№@«»%\\\"]\\s?)+$"
	companyNameMatch := "\\S+"
	companyName := regexp.MustCompile(companyNameMatch)

	return companyName.MatchString(value)
}

func lettersSymbolsValid(value string) bool {
	lettersSymbolsMatch := "^([а-яА-ЯёЁa-zA-Z-.]\\s?)+$"
	lettersSymbols := regexp.MustCompile(lettersSymbolsMatch)

	return lettersSymbols.MatchString(value)
}

func lettersSymbolsNumbersValid(value string) bool {
	lettersSymbolsNumbersMatch := "^([а-яА-ЯёЁa-zA-Z\\d-.,/()№«»\"]\\s?)+$"
	lettersSymbolsNumbers := regexp.MustCompile(lettersSymbolsNumbersMatch)

	return lettersSymbolsNumbers.MatchString(value)
}

func phoneNumberValid(value string) bool {
	phoneNumberMatch := "^[+7][\\d]+$"
	phoneNumber := regexp.MustCompile(phoneNumberMatch)

	return phoneNumber.MatchString(value)
}

func emailValid(value string) bool {
	emailMatch := "^.+@[a-zA-Z][a-zA-Z\\.]+$"
	email := regexp.MustCompile(emailMatch)

	return email.MatchString(value)
}

func notNumbersError(fieldName string, value string) error {
	return fmt.Errorf("%s must contain only numbers, %s=%s", fieldName, fieldName, value)
}

func wrongLengthError(fieldName string, value string) error {
	return fmt.Errorf("%s has wrong length, length=%d, %s=%s", fieldName, len(value), fieldName, value)
}

func failedValidationError(fieldName string, value string) error {
	return fmt.Errorf("%s failed validation, %s=%s", fieldName, fieldName, value)
}
