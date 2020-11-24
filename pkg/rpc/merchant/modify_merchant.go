package merchant

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/errors"
	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
	"github.com/shomali11/util/xstrings"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

func (controller Controller) ModifyMerchant(ctx context.Context, req *protos.ModifyMerchantRequest) (*protos.ModifyMerchantReply, error) {

	logger := rpclog.WithRequestId(ctx, log.GetLogger()).WithField("requestData", req)

	if req == nil {
		logger.Error("get request data is nil")
		return nil, fmt.Errorf(constant.ErrorMessageGetRequestDataIsNil)
	}

	if err := controller.validateModifyMerchantRequest(req); err != nil {

		logger.WithError(err).Info("validate parameter failed")
		return &protos.ModifyMerchantReply{
			Err: convertError(err, int64(errors.RequestParamError.ErrorCode)),
		}, nil
	}

	merchant, err := getMerchantByMerchantId(req.GetMerchantId())
	if err != nil {

		if !gorm.IsRecordNotFoundError(err) {
			logger.WithError(err).Error("get merchant error")
			return &protos.ModifyMerchantReply{
				Err: convertError(err, constant.ErrorCodeGetMerchantError),
			}, nil
		}

		logger.WithError(err).Info("merchant not exist")

		return &protos.ModifyMerchantReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeMerchantIdIsNotExist,
				Message: constant.ErrorMessageMerchantIdIsNotExit,
			},
		}, nil
	}

	logger = logger.WithField("merchant", merchant)

	record := map[string]interface{}{}
	if xstrings.IsNotBlank(req.GetName()) {
		record["name"] = req.GetName()
	}

	if xstrings.IsNotBlank(req.GetArea()) {
		record["area"] = req.GetArea()
	}

	if xstrings.IsNotBlank(req.GetManagementCentre()) {
		record["managementCentre"] = req.GetManagementCentre()
	}

	if xstrings.IsNotBlank(req.GetNetworkStation()) {
		record["networkStation"] = req.GetNetworkStation()
	}

	if xstrings.IsNotBlank(req.GetMobile()) {
		record["mobile"] = req.Mobile
	}

	if xstrings.IsNotBlank(req.GetPassword()) {
		record["salt"], record["password"] = hashPassword(req.Password)
	}

	record["isRebate"] = req.GetIsRebate()

	if req.GetStatus() >= 0 {
		record["status"] = req.GetStatus()
	}

	err = database.GetDB(constant.DatabaseConfigKey).Table(merchant.TableName()).Where("merchantId = ?", merchant.MerchantId).Updates(record).Error
	if err != nil {

		logger.WithError(err).Error("update merchant error")
		return &protos.ModifyMerchantReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeUpdateMerchantError,
				Message: constant.ErrorMessageUpdateMerchantError,
			},
		}, nil
	}

	return &protos.ModifyMerchantReply{
		Merchant: &protos.Merchant{
			MerchantId: merchant.MerchantId,
		},
	}, nil
}

func (controller Controller) validateModifyMerchantRequest(req *protos.ModifyMerchantRequest) error {

	return validation.ValidateStruct(req,
		validation.Field(&req.MerchantId, MerchantIdRule...),
	)
}
