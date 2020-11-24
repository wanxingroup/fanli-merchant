package merchant

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/errors"
	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

func (controller Controller) ModifyPassword(ctx context.Context, req *protos.ModifyPasswordRequest) (*protos.ModifyPasswordReply, error) {

	logger := rpclog.WithRequestId(ctx, log.GetLogger()).WithField("requestData", req)

	if req == nil {
		logger.Error("get request data is nil")
		return nil, fmt.Errorf(constant.ErrorMessageGetRequestDataIsNil)
	}

	if err := controller.validateModifyPasswordRequest(req); err != nil {

		logger.WithError(err).Info("validate parameter failed")
		return &protos.ModifyPasswordReply{
			Err: convertError(err, int64(errors.RequestParamError.ErrorCode)),
		}, nil
	}

	merchant, err := getMerchantByMerchantId(req.GetMerchantId())
	if err != nil {

		if !gorm.IsRecordNotFoundError(err) {
			logger.WithError(err).Error("get merchant error")
			return &protos.ModifyPasswordReply{
				Err: convertError(err, constant.ErrorCodeGetMerchantError),
			}, nil
		}

		logger.WithError(err).Info("merchant not exist")

		return &protos.ModifyPasswordReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeMerchantIdIsNotExist,
				Message: constant.ErrorMessageMerchantIdIsNotExit,
			},
		}, nil
	}

	logger = logger.WithField("merchant", merchant)

	merchant.Salt, merchant.Password = hashPassword(req.Password)
	err = database.GetDB(constant.DatabaseConfigKey).Model(&model.Merchant{}).Updates(merchant).Error
	if err != nil {

		logger.WithError(err).Error("update merchant error")
		return &protos.ModifyPasswordReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeUpdateMerchantError,
				Message: constant.ErrorMessageUpdateMerchantError,
			},
		}, nil
	}

	return &protos.ModifyPasswordReply{
		Merchant: &protos.Merchant{
			MerchantId: merchant.MerchantId,
			Mobile:     merchant.Mobile,
		},
	}, nil
}

func (controller Controller) validateModifyPasswordRequest(req *protos.ModifyPasswordRequest) error {

	return validation.ValidateStruct(req,
		validation.Field(&req.MerchantId, MerchantIdRule...),
		validation.Field(&req.Password, PasswordRules...),
	)
}
