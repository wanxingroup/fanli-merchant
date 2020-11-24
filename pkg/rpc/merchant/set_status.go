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
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

func (controller Controller) SetStatus(ctx context.Context, req *protos.SetStatusRequest) (*protos.SetStatusReply, error) {

	logger := rpclog.WithRequestId(ctx, log.GetLogger()).WithField("requestData", req)
	logger.Info("Merchant SetStatus")

	if req == nil {
		logger.Error("get request data is nil")
		return nil, fmt.Errorf(constant.ErrorMessageGetRequestDataIsNil)
	}

	if err := controller.validateSetStatusRequest(req); err != nil {

		logger.WithError(err).Info("validate parameter failed")
		return &protos.SetStatusReply{
			Err: convertError(err, int64(errors.RequestParamError.ErrorCode)),
		}, nil
	}

	merchant, err := getMerchantByMerchantId(req.GetMerchantId())
	if err != nil {

		if !gorm.IsRecordNotFoundError(err) {
			logger.WithError(err).Error("get merchant error")
			return &protos.SetStatusReply{
				Err: convertError(err, constant.ErrorCodeGetMerchantError),
			}, nil
		}

		logger.WithError(err).Info("merchant not exist")

		return &protos.SetStatusReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeMerchantIdIsNotExist,
				Message: constant.ErrorMessageMerchantIdIsNotExit,
			},
		}, nil
	}

	record := map[string]interface{}{
		"status": req.GetStatus(),
	}

	// 禁用时不返利
	if req.GetStatus() == protos.MerchantStatus_MerchantStatusClose {
		record["isRebate"] = 0
	}

	err = database.GetDB(constant.DatabaseConfigKey).Table(merchant.TableName()).Where("merchantId = ?", merchant.MerchantId).Updates(record).Error

	logger.Info("SetStatus", record)

	if err != nil {
		logger.WithError(err).Error("update merchant error")
		return &protos.SetStatusReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeUpdateMerchantError,
				Message: constant.ErrorMessageUpdateMerchantError,
			},
		}, nil
	}

	return &protos.SetStatusReply{
		IsSuccess: true,
	}, nil
}

func (controller Controller) validateSetStatusRequest(req *protos.SetStatusRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.MerchantId, MerchantIdRule...),
	)
}
