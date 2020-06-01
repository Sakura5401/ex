package handlers

import (
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	pb "gitlab.safecrow.ru/safecrow/gateway-requisites/v2/api/proto"
	"gitlab.safecrow.ru/safecrow/gateway-requisites/v2/internal"
)

// HandlerSetBankDetail хендлер создания банковских реквизитов
func HandlerSetBankDetail(message []byte, db *gorm.DB) []byte {
	var request pb.SetBankDetailRequest

	err := proto.Unmarshal(message, &request)

	if err != nil {
		log.Error().Err(err)

		got, _ := proto.Marshal(&pb.SetBankDetailResponse{
			Header: &pb.Header{
				Msg:  err.Error(),
				Code: 400,
			},
		})

		return got
	}

	res, err := internal.SetBankDetail(db, &request)

	if err != nil {
		log.Error().Err(err)

		got, _ := proto.Marshal(&pb.SetBankDetailResponse{
			Header: &pb.Header{
				Msg:  err.Error(),
				Code: 400,
			},
		})

		return got
	}

	res.Header = &pb.Header{Code: 200}
	got, _ := proto.Marshal(res)

	return got
}

// HandlerGetBankDetails хендлер получения всех банковских реквизитов
func HandlerGetBankDetails(message []byte, db *gorm.DB) []byte {
	var request pb.GetBankDetailsRequest

	err := proto.Unmarshal(message, &request)

	if err != nil {
		log.Error().Err(err)

		got, _ := proto.Marshal(&pb.GetBankDetailsResponse{
			Header: &pb.Header{
				Msg:  err.Error(),
				Code: 400,
			},
		})

		return got
	}

	res, err := internal.GetBankDetails(db, &request)

	if err != nil {
		log.Error().Err(err)

		got, _ := proto.Marshal(&pb.GetBankDetailsResponse{
			Header: &pb.Header{
				Msg:  err.Error(),
				Code: 400,
			},
		})

		return got
	}

	res.Header = &pb.Header{Code: 200}
	got, _ := proto.Marshal(res)

	return got
}

// HandlerGetBankDetail хендлер получения банковских реквизитов
func HandlerGetBankDetail(message []byte, db *gorm.DB) []byte {
	var request pb.GetBankDetailRequest

	err := proto.Unmarshal(message, &request)

	if err != nil {
		log.Error().Err(err)

		got, _ := proto.Marshal(&pb.GetBankDetailResponse{
			Header: &pb.Header{
				Msg:  err.Error(),
				Code: 400,
			},
		})

		return got
	}

	res, err := internal.GetBankDetail(db, &request)

	if err != nil {
		log.Error().Err(err)

		got, _ := proto.Marshal(&pb.GetBankDetailResponse{
			Header: &pb.Header{
				Msg:  err.Error(),
				Code: 400,
			},
		})

		return got
	}

	res.Header = &pb.Header{Code: 200}
	got, _ := proto.Marshal(res)

	return got
}

// HandlerUpdateDefaultPayments хендлер назначения основного платежного метода
func HandlerUpdateDefaultPayments(message []byte, db *gorm.DB) []byte {
	var request pb.UpdateDefaultPaymentsRequest

	err := proto.Unmarshal(message, &request)

	if err != nil {
		log.Error().Err(err)

		got, _ := proto.Marshal(&pb.UpdateDefaultPaymentsResponse{
			Header: &pb.Header{
				Msg:  err.Error(),
				Code: 400,
			},
		})

		return got
	}

	res, err := internal.UpdateDefaultPayments(db, &request)

	if err != nil {
		log.Error().Err(err)

		got, _ := proto.Marshal(&pb.UpdateDefaultPaymentsResponse{
			Header: &pb.Header{
				Msg:  err.Error(),
				Code: 400,
			},
		})

		return got
	}

	res.Header = &pb.Header{Code: 200}
	got, _ := proto.Marshal(res)

	return got
}
