package merchant

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/errors"
	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

func (controller *Controller) GetMerchant(ctx context.Context, req *protos.GetMerchantRequest) (*protos.GetMerchantReply, error) {

	logger := rpclog.WithRequestId(ctx, log.GetLogger()).WithField("req", req)

	err := controller.validateGetMerchantRequestData(req)
	if err != nil {

		logger.WithField("error", err).Info("validate get.go merchant data invalid")
		return &protos.GetMerchantReply{
			Err: convertError(err, errors.CodeRequestParamError),
		}, nil
	}

	merchant, err := getMerchantByMerchantId(req.GetMerchantId())
	if err != nil {

		if !gorm.IsRecordNotFoundError(err) {
			logger.WithField("error", err).Error("get merchant error")
			return &protos.GetMerchantReply{
				Err: convertError(err, constant.ErrorCodeGetMerchantError),
			}, nil
		}

		return &protos.GetMerchantReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeMerchantIdIsNotExist,
				Message: constant.ErrorMessageMerchantIdIsNotExit,
			},
		}, nil
	}

	return &protos.GetMerchantReply{
		Merchant: &protos.Merchant{
			MerchantId:       merchant.MerchantId,
			Mobile:           merchant.Mobile,
			IsRebate:         merchant.IsRebate == true,
			Name:             merchant.Name,
			Area:             merchant.Area,
			ManagementCentre: merchant.ManagementCentre,
			NetworkStation:   merchant.NetworkStation,
			Status:           ConvertMerchantStatusToProtobuf(merchant.Status),
		},
	}, nil
}

func (controller *Controller) validateGetMerchantRequestData(req *protos.GetMerchantRequest) error {

	return validation.ValidateStruct(req,
		validation.Field(&req.MerchantId, MerchantIdRule...),
	)
}
