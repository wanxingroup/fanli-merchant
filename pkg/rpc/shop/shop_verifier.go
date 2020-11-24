package shop

import (
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/jinzhu/gorm"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

func CreateVerifier(shopId uint64, data []string) error {

	tx := database.GetDB(constant.DatabaseConfigKey).Begin()

	for _, datum := range data {
		record := &model.ShopVerification{
			ShopId: shopId,
			Mobile: datum,
		}

		if err := tx.Create(&record).Error; err != nil {
			log.GetLogger().WithField("createVerifier", record).WithError(err).Error("create createVerifier error")
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func ClearVerifier(shopId uint64) {
	database.GetDB(constant.DatabaseConfigKey).Where("shopId = ?", shopId).Unscoped().Delete(&model.ShopVerification{})
}

func CheckVerifier(mobile string, shopId uint64) bool {
	var result *model.ShopVerification
	err := database.GetDB(constant.DatabaseConfigKey).Where(&model.ShopVerification{
		ShopId: shopId,
		Mobile: mobile,
	}).First(result).Error

	if err == gorm.ErrRecordNotFound {
		log.GetLogger().WithField("mobile", mobile).Info("This mobile is not verifier")

		return false
	}

	if err != nil {
		log.GetLogger().WithError(err).Error("get Mobile verifier error ")
		return false
	}

	return true
}
