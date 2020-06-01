package internal

import (
	"errors"
	"regexp"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	pb "gitlab.safecrow.ru/safecrow/gateway-requisites/v2/api/proto"
	"gitlab.safecrow.ru/safecrow/gateway-requisites/v2/internal/models"
)

// SetBankDetail создание банковских реквизитов.
// Первые банковские реквизиты = основной платёжный метод, если ранее не было создано ни банковских карт, ни реквизитов.
func SetBankDetail(db *gorm.DB, request *pb.SetBankDetailRequest) (*pb.SetBankDetailResponse, error) {
	userID, _ := uuid.FromString(request.UserID)
	userType := 0 //now it's default legal entity = 0

	if !ValidateCheckingAccount(request.CheckingAccount, userType) {
		return nil, errors.New("CheckingAccount failed validation = " + request.CheckingAccount)
	}

	detail, err := IsExistDefaultDetail(db, userID)

	if err != nil {
		log.Error().Msg("Error in IsExistDefaultDetail")
		return nil, err
	}

	card, err := IsExistDefaultCard(db, userID)

	if err != nil {
		log.Error().Msg("Error in IsExistDefaultCard")
		return nil, err
	}

	if detail == nil && card == nil {
		log.Info().Msgf("Record'll be create and it'll be default method of payment")

		model := models.BankDetail{
			UserID:                 userID,
			DetailName:             request.DetailName,
			BankIdentificationCode: request.BankIdentificationCode,
			BankName:               request.BankName,
			City:                   request.City,
			Address:                request.Address,
			CorrAccount:            request.CorrAccount,
			CheckingAccount:        request.CheckingAccount,
			DefaultDetail:          true,
		}

		err := db.Create(&model).Error

		if err != nil {
			log.Error().Msg("Error when create requisite")
			return nil, err
		}

		res := &pb.SetBankDetailResponse{
			RequisiteID:   model.UUID.String(),
			DefaultDetail: model.DefaultDetail,
		}

		log.Info().
			Str("requisiteID", res.RequisiteID).
			Msg("Create requisites")

		return res, nil
	}

	model := models.BankDetail{
		UserID:                 userID,
		DetailName:             request.DetailName,
		BankIdentificationCode: request.BankIdentificationCode,
		BankName:               request.BankName,
		City:                   request.City,
		Address:                request.Address,
		CorrAccount:            request.CorrAccount,
		CheckingAccount:        request.CheckingAccount,
		DefaultDetail:          false,
	}

	err = db.Create(&model).Error
	if err != nil {
		log.Error().Msg("Error when create requisite")
		return nil, err
	}

	res := &pb.SetBankDetailResponse{
		RequisiteID:   model.UUID.String(),
		DefaultDetail: model.DefaultDetail,
	}

	log.Info().
		Str("requisiteID", res.RequisiteID).
		Msg("Create requisites")

	return res, nil
}

// GetBankDetails получение всех банковских реквизитов
func GetBankDetails(db *gorm.DB, request *pb.GetBankDetailsRequest) (*pb.GetBankDetailsResponse, error) {
	userID, _ := uuid.FromString(request.UserID)

	var rows []models.BankDetail

	db.Model(&models.BankDetail{}).Where("user_id = ?", userID).Find(&rows)

	var resData []*pb.BankDetailInfo

	for _, data := range rows {
		row := &pb.BankDetailInfo{
			RequisiteID:            data.UUID.String(),
			DetailName:             data.DetailName,
			BankIdentificationCode: data.BankIdentificationCode,
			BankName:               data.BankName,
			City:                   data.City,
			Address:                data.Address,
			CorrAccount:            data.CorrAccount,
			CheckingAccount:        data.CheckingAccount,
			DefaultDetail:          data.DefaultDetail,
		}
		resData = append(resData, row)
	}

	res := &pb.GetBankDetailsResponse{
		BankDetailInfo: resData,
	}

	return res, nil
}

// GetBankDetail получение банковских реквизитов
func GetBankDetail(db *gorm.DB, request *pb.GetBankDetailRequest) (*pb.GetBankDetailResponse, error) {
	requisiteID, err := uuid.FromString(request.RequisiteID)
	if err != nil {
		log.Error().Msg("Error when convert requisiteID from string to uuid")
		return &pb.GetBankDetailResponse{}, err
	}

	var data models.BankDetail

	db.Where("uuid = ?", requisiteID).First(&data)

	res := &pb.GetBankDetailResponse{
		RequisiteID:            data.UUID.String(),
		DetailName:             data.DetailName,
		BankIdentificationCode: data.BankIdentificationCode,
		BankName:               data.BankName,
		City:                   data.City,
		Address:                data.Address,
		CorrAccount:            data.CorrAccount,
		CheckingAccount:        data.CheckingAccount,
		DefaultDetail:          data.DefaultDetail,
	}

	return res, nil
}

// ValidateCheckingAccount валидация расчетного счета банковских реквизитов.
// Для юридических лиц расчетный счет начинается с цифр 40802 или 407,
// для физических - не может начинаться с цифр 40802 или 407.
func ValidateCheckingAccount(checkingAccount string, userType int) bool {
	match, _ := regexp.MatchString("^(40802|407)", checkingAccount) //

	if userType == 0 {
		return match // юрики и ИП + расчётный начинается с цифр 40802 или 407
	}

	return !match // физики + расчётный НЕ начинается с цифр 40802 или 407
}

// UpdateDefaultPayments назначить основным платёжным методом.
// Основной платёжный метод один - либо банковские реквизиты, либо банковская карта.
// Если ранее уже был назначен платёжный метод, то заменяем на новый.
func UpdateDefaultPayments(db *gorm.DB, request *pb.UpdateDefaultPaymentsRequest) (
	*pb.UpdateDefaultPaymentsResponse, error) {
	userID, err := uuid.FromString(request.UserID)
	if err != nil {
		log.Error().Err(err).Msgf("Error when convert `request.UserID` from string to uuid: %s", request.UserID)
		return nil, err
	}

	infoID, err := uuid.FromString(request.RequisiteID)
	if err != nil {
		log.Error().Err(err).Msgf("Error when convert `request.RequisiteID` from string to uuid: %s", request.RequisiteID)
		return nil, err
	}

	detail, err := IsExistBankDetail(db, infoID)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	card, err := IsExistBankCard(db, infoID)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if detail == nil && card == nil {
		return nil, errors.New("Not found UUID = " + infoID.String())
	}

	oldDefaultDetail, err := IsExistDefaultDetail(db, userID)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	oldDefaultCard, err := IsExistDefaultCard(db, userID)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if oldDefaultDetail != nil {
		oldDefaultDetail.DefaultDetail = false
		oldDefaultDetail.UpdatedAt = time.Time{}
		err = db.Save(&oldDefaultDetail).Error

		if err != nil {
			log.Error().Msg("Error when update old BankDetail")
			return nil, err
		}
	}

	if oldDefaultCard != nil {
		oldDefaultCard.DefaultCard = false
		oldDefaultCard.UpdatedAt = time.Time{}
		err = db.Save(&oldDefaultCard).Error

		if err != nil {
			log.Error().Msg("Error when update old BankCard")
			return nil, err
		}
	}

	res := &pb.UpdateDefaultPaymentsResponse{}

	if detail != nil {
		detail.DefaultDetail = true
		detail.UpdatedAt = time.Time{}
		err = db.Save(&detail).Error

		if err != nil {
			log.Error().Msg("Error when update new BankDetail")
			return nil, err
		}

		res = &pb.UpdateDefaultPaymentsResponse{
			DefaultPayment: detail.DefaultDetail,
		}
	}

	if card != nil {
		card.DefaultCard = true
		card.UpdatedAt = time.Time{}

		err = db.Save(&card).Error
		if err != nil {
			log.Error().Msg("Error when update new BankCard")
			return nil, err
		}

		res = &pb.UpdateDefaultPaymentsResponse{
			DefaultPayment: card.DefaultCard,
		}
	}

	log.Info().
		Str("UUID", infoID.String()).
		Msg("Update default payment")

	return res, nil
}
