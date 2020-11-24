package shop

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
	idcreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/validate"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func (_ Controller) CreateShop(ctx context.Context, req *protos.CreateShopRequest) (reply *protos.CreateShopReply, err error) {

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	validateDate := &model.Shop{
		MerchantId:         req.GetMerchantId(),
		Name:               req.GetName(),
		BusinessHoursStart: uint32(req.GetBusinessHoursStart()),
		BusinessHoursEnd:   uint32(req.GetBusinessHoursEnd()),
		ShopStatus:         model.ShopStatus(req.GetShopStatus()),
		ShopType:           model.ShopType(req.GetShopType()),
		FuYouMerchantId:    req.GetFuYouMerchantId(),
		IsRebate:           req.IsRebate,
	}

	ok, errMessage := validate.ValidateShopInfo(validateDate, false)

	if !ok {
		return &protos.CreateShopReply{
			Err: &protos.Error{
				Code:    errMessage.GetCode(),
				Message: errMessage.GetMessage(),
				Stack:   errMessage.GetStack(),
			},
		}, nil
	}

	// 校验 核销员
	verifierMobilesList := req.GetVerifierMobiles()
	if len(verifierMobilesList) > 0 {
		verifierMobilesList = RemoveRepeatedElement(verifierMobilesList)
		for _, verifierMobile := range verifierMobilesList {
			if !validate.CheckMobile(verifierMobile) {
				return &protos.CreateShopReply{
					Err: &protos.Error{
						Code:    constant.ErrorCodeMobileFormatInvalid,
						Message: constant.ErrorMessageMobileFormatInvalid,
						Stack:   nil,
					},
				}, nil
			}
		}
	}

	shopId, errCreateShop := createShop(&model.Shop{
		MerchantId:         validateDate.MerchantId,
		Name:               validateDate.Name,
		BusinessHoursStart: validateDate.BusinessHoursStart,
		BusinessHoursEnd:   validateDate.BusinessHoursEnd,
		ShopStatus:         validateDate.ShopStatus,
		ShopType:           validateDate.ShopType,
		FuYouMerchantId:    validateDate.FuYouMerchantId,
		Time:               databases.Time{},
		IsRebate:           validateDate.IsRebate,
	})
	if errCreateShop != nil {

		log.GetLogger().WithError(err).Error("register shop error")
		return &protos.CreateShopReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeCreateShopFailed,
				Message: constant.ErrorMessageCreateShopFailed,
				Stack:   nil,
			},
		}, nil
	}
	if len(verifierMobilesList) > 0 {
		createVerifierErr := CreateVerifier(shopId, verifierMobilesList)
		if createVerifierErr != nil {

			log.GetLogger().WithError(err).Error("create verifier  error")
			return &protos.CreateShopReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodeCreateVerifierErr,
					Message: constant.ErrorMessageCodeCreateVerifierErr,
					Stack:   nil,
				},
			}, nil
		}
	}

	return &protos.CreateShopReply{
		ShopId: shopId,
	}, nil

}

func createShop(shop *model.Shop) (uint64, error) {

	record := &model.Shop{
		ShopId:             idcreator.NextID(),
		MerchantId:         shop.MerchantId,
		Name:               shop.Name,
		BusinessHoursStart: shop.BusinessHoursStart,
		BusinessHoursEnd:   shop.BusinessHoursEnd,
		ShopStatus:         shop.ShopStatus,
		ShopType:           shop.ShopType,
		FuYouMerchantId:    shop.FuYouMerchantId,
		IsRebate:           shop.IsRebate,
	}

	err := database.GetDB(constant.DatabaseConfigKey).Create(record).Error
	if err != nil {

		log.GetLogger().WithField("shop", record).WithError(err).Error("create record error")
		return 0, err
	}

	return record.ShopId, nil
}

// 数组去重
func RemoveRepeatedElement(arr []string) (newArr []string) {

	result := make([]string, 0, len(arr))
	temp := map[string]struct{}{}
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
