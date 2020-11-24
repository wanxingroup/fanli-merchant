package shop

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/validate"
)

func (_ Controller) ModifyShopInfo(ctx context.Context, req *protos.ModifyShopInfoRequest) (reply *protos.ModifyShopReply, err error) {

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	validateDate := &model.Shop{
		ShopId:             req.GetShopId(),
		MerchantId:         req.GetMerchantId(),
		Name:               req.GetName(),
		BusinessHoursStart: uint32(req.GetBusinessHoursStart()),
		BusinessHoursEnd:   uint32(req.GetBusinessHoursEnd()),
		ShopStatus:         model.ShopStatus(req.GetShopStatus()),
		ShopType:           model.ShopType(req.GetShopType()),
		FuYouMerchantId:    req.GetFuYouMerchantId(),
		IsRebate:           req.IsRebate,
	}

	ok, errMessage := validate.ValidateShopInfo(validateDate, true)

	if !ok {
		return &protos.ModifyShopReply{
			Err: &protos.Error{
				Code:    errMessage.GetCode(),
				Message: errMessage.GetMessage(),
				Stack:   errMessage.GetStack(),
			},
		}, nil
	}

	errEditShop := modifyShop(validateDate)
	if errEditShop != nil {

		log.GetLogger().WithError(err).Error("update shopInfo error")
		return &protos.ModifyShopReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeUpdateShopFailed,
				Message: constant.ErrorMessageUpdateShopFailed,
				Stack:   nil,
			},
		}, nil
	}

	verifierMobilesList := req.GetVerifierMobiles()
	if len(verifierMobilesList) > 0 {
		verifierMobilesList = RemoveRepeatedElement(verifierMobilesList)
		for _, verifierMobile := range verifierMobilesList {
			if !validate.CheckMobile(verifierMobile) {
				return &protos.ModifyShopReply{
					Err: &protos.Error{
						Code:    constant.ErrorCodeMobileFormatInvalid,
						Message: constant.ErrorMessageMobileFormatInvalid,
						Stack:   nil,
					},
				}, nil
			}
		}

		ClearVerifier(req.GetShopId())

		createVerifierErr := CreateVerifier(req.GetShopId(), verifierMobilesList)
		if createVerifierErr != nil {

			log.GetLogger().WithError(err).Error("create verifier  error")
			return &protos.ModifyShopReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodeCreateVerifierErr,
					Message: constant.ErrorMessageCodeCreateVerifierErr,
					Stack:   nil,
				},
			}, nil
		}
	}

	return &protos.ModifyShopReply{
		ShopId: validateDate.ShopId,
	}, nil
}

func modifyShop(shopInfo *model.Shop) error {

	record := map[string]interface{}{
		"name":               shopInfo.Name,
		"businessHoursStart": shopInfo.BusinessHoursStart,
		"businessHoursEnd":   shopInfo.BusinessHoursEnd,
		"shopStatus":         shopInfo.ShopStatus,
		"shopType":           shopInfo.ShopType,
		"fuYouMerchantId":    shopInfo.FuYouMerchantId,
		"isRebate":           shopInfo.IsRebate,
	}

	err := database.GetDB(constant.DatabaseConfigKey).Table(shopInfo.TableName()).Where("shopId = ?", shopInfo.ShopId).Updates(record).Error
	if err != nil {
		log.GetLogger().WithField("shop", record).WithError(err).Error("update record error")
		return err
	}

	return nil
}
