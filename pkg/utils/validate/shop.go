package validate

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/jinzhu/gorm"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

const shopNameMaxLength = 20

func ValidateShopInfo(data *model.Shop, isUpdate bool) (bool, *protos.Error) {

	// 校验店铺名称长度
	if len([]rune(data.Name)) > shopNameMaxLength {
		return false, &protos.Error{
			Code:    constant.ErrorCodeShopNameTooLarge,
			Message: constant.ErrorMessageShopNameTooLarge,
			Stack:   nil,
		}
	}

	// 校验营业时间是否准确
	if data.BusinessHoursStart < 0 || data.BusinessHoursEnd < 0 || data.BusinessHoursStart > 86400 || data.BusinessHoursEnd > 172800 {
		return false, &protos.Error{
			Code:    constant.ErrorCodeBusinessTimeError,
			Message: constant.ErrorMessageBusinessTimeError,
			Stack:   nil,
		}
	}

	// 校验门店状态
	if data.ShopStatus != model.ShopStatusOpen && data.ShopStatus != model.ShopStatusClose {
		return false, &protos.Error{
			Code:    constant.ErrorCodeBusinessStatusError,
			Message: constant.ErrorMessageBusinessStatusError,
			Stack:   nil,
		}
	}

	// 校验门店类型
	if data.ShopType != model.ShopTypeDirect && data.ShopType != model.ShopTypeLeague {
		return false, &protos.Error{
			Code:    constant.ErrorCodeBusinessTypeError,
			Message: constant.ErrorMessageBusinessTypeError,
			Stack:   nil,
		}
	}

	// 校验商家是否存在
	if !checkMerchantIdExist(data.MerchantId) {
		return false, &protos.Error{
			Code:    constant.ErrorCodeMerchantIdIsNotExist,
			Message: constant.ErrorMessageMerchantIdIsNotExit,
			Stack:   nil,
		}
	}

	// 校验店铺名是否重复
	if checkShopNameExist(data.Name, data.MerchantId, data.ShopId) {
		return false, &protos.Error{
			Code:    constant.ErrorCodeShopNameIsExist,
			Message: constant.ErrorMessageShopNameIsExist,
			Stack:   nil,
		}
	}

	return true, nil
}

func checkMerchantIdExist(merchantId uint64) bool {

	data := &model.Merchant{}
	err := database.GetDB(constant.DatabaseConfigKey).Where(&model.Merchant{
		MerchantId: merchantId,
	}).First(data).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}

		log.GetLogger().WithError(err).Error("get merchant record error")
		return true
	}

	if data.MerchantId > 0 {
		return true
	}

	return false
}

func checkShopNameExist(name string, merchantId, shopId uint64) bool {
	data := &model.Shop{}

	db := database.GetDB(constant.DatabaseConfigKey)

	err := db.Where(&model.Shop{
		Name:       name,
		MerchantId: merchantId,
	}).First(data).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}

		log.GetLogger().WithError(err).Error("get shop record error")
		return true
	}

	return data.ShopId != shopId
}
