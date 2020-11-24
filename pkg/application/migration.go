package application

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
)

func autoMigration() {

	db := database.GetDB(constant.DatabaseConfigKey)
	db.AutoMigrate(model.Merchant{})
	db.AutoMigrate(model.Shop{})
	db.AutoMigrate(model.ShopVerification{})
	db.AutoMigrate(model.MerchantCode{})
	db.AutoMigrate(model.Staff{})
}
