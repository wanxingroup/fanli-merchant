package merchant

import (
	"context"
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	idCreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/jinzhu/gorm"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/validate"
)

const passwordMinLength = 8
const passwordMaxLength = 50

func (_ Controller) RegisterUser(_ context.Context, req *protos.RegisterUserRequest) (reply *protos.RegisterUserReply, err error) {

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	if !validate.CheckMobile(req.GetMobile()) {

		return &protos.RegisterUserReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeMobileFormatInvalid,
				Message: constant.ErrorMessageMobileFormatInvalid,
				Stack:   nil,
			},
		}, nil
	}

	if checkMobileExist(req.GetMobile()) {

		return &protos.RegisterUserReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeMobileExist,
				Message: constant.ErrorMessageMobileExist,
				Stack:   nil,
			},
		}, nil
	}

	if len(req.Password) < passwordMinLength {

		return &protos.RegisterUserReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodePasswordLengthNotEnough,
				Message: constant.ErrorMessagePasswordLengthNotEnough,
				Stack:   nil,
			},
		}, nil
	}
	if len(req.Password) > passwordMaxLength {

		return &protos.RegisterUserReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodePasswordLengthTooLarge,
				Message: constant.ErrorMessagePasswordLengthTooLarge,
				Stack:   nil,
			},
		}, nil
	}

	merchantId, err := registerUser(req.Mobile, req.Password, req.Name, req.Area, req.ManagementCentre, req.NetworkStation, req.IsRebate, model.MerchantStatus(req.Status))
	if err != nil {

		log.GetLogger().WithError(err).Error("register user error")
		return &protos.RegisterUserReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeRegisterUserFailed,
				Message: constant.ErrorMessageRegisterUserFailed,
				Stack:   nil,
			},
		}, nil
	}

	return &protos.RegisterUserReply{
		MerchantId: merchantId,
	}, nil
}

func registerUser(mobile string, password string, name string, area string, managementCentre string, networkStation string, isRebate bool, status model.MerchantStatus) (uint64, error) {

	salt, encodedPassword := hashPassword(password)
	record := &model.Merchant{
		MerchantId:       idCreator.NextID(),
		Mobile:           mobile,
		Password:         encodedPassword,
		Salt:             salt,
		Name:             name,
		Area:             area,
		ManagementCentre: managementCentre,
		NetworkStation:   networkStation,
		IsRebate:         isRebate,
		Status:           status,
	}

	err := database.GetDB(constant.DatabaseConfigKey).Create(record).Error
	if err != nil {

		log.GetLogger().WithField("user", record).WithError(err).Error("create record error")
		return 0, err
	}

	return record.MerchantId, nil
}

func hashPassword(passwordString string) (string, string) {

	return password.Encode(passwordString, getPasswordOptions())
}

func checkMobileExist(mobile string) bool {

	data := &model.Merchant{}
	err := database.GetDB(constant.DatabaseConfigKey).Where(&model.Merchant{
		Mobile: mobile,
	}).First(data).Error

	if err != gorm.ErrRecordNotFound || data.MerchantId > 0 {
		return true
	}

	return false
}
