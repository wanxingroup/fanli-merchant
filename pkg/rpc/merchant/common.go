package merchant

import (
	"crypto/sha1"
	"fmt"
	"strconv"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/databases"
	"github.com/anaskhan96/go-password-encoder"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jinzhu/gorm"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

const saltLen = 10
const iterations = 10000
const keyLen = 50
const DefaultExpireTime = 10 // 十分钟

func getPasswordOptions() *password.Options {
	return &password.Options{
		SaltLen:      saltLen,
		Iterations:   iterations,
		KeyLen:       keyLen,
		HashFunction: sha1.New,
	}
}

func getMerchantByMobile(mobile string) (*model.Merchant, error) {

	merchant := &model.Merchant{}
	err := database.GetDB(constant.DatabaseConfigKey).Where(&model.Merchant{
		Mobile: mobile,
	}).First(&merchant).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.GetLogger().WithError(err).Error("find merchant error")
		} else {
			log.GetLogger().Info("not found merchant")
		}

		return nil, err
	}

	return merchant, nil
}

func getMerchantByMerchantId(merchantId uint64) (*model.Merchant, error) {

	merchant := &model.Merchant{}
	err := database.GetDB(constant.DatabaseConfigKey).Where(&model.Merchant{
		MerchantId: merchantId,
	}).First(&merchant).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.GetLogger().WithError(err).Error("find merchant error")
		} else {
			log.GetLogger().Info("not found merchant")
		}

		return nil, err
	}

	return merchant, nil
}

func getMerchantListByConditions(conditions map[string]interface{}, pageData map[string]uint64) ([]*model.Merchant, uint64, error) {
	db := database.GetDB(constant.DatabaseConfigKey).Model(&model.Merchant{})

	if merchantId, has := conditions["merchantId"]; has {
		db = db.Where("merchantId = ?", merchantId)
	}

	if mobileFuzzySearch, has := conditions["mobileFuzzySearch"]; has {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobileFuzzySearch))
	}

	if nameFuzzySearch, has := conditions["nameFuzzySearch"]; has {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", nameFuzzySearch))
	}

	if networkStationFuzzySearch, has := conditions["networkStationFuzzySearch"]; has {
		db = db.Where("networkStation LIKE ?", fmt.Sprintf("%%%s%%", networkStationFuzzySearch))
	}

	if area, has := conditions["area"]; has {
		db = db.Where("area = ?", area)
	}

	if managementCentre, has := conditions["managementCentre"]; has {
		db = db.Where("managementCentre = ?", managementCentre)
	}

	db = db.Order("merchantId DESC")

	results := make([]*model.Merchant, 0, pageData["pageSize"])
	var count uint64
	err := databases.FindPage(db, pageData, &results, &count)
	return results, count, err
}

func convertError(err error, defaultErrorCode int64) *protos.Error {

	if err == nil {
		return nil
	}

	log.GetLogger().Debugf("error type: %T", err)
	if validationError, ok := err.(validation.Error); ok {

		log.GetLogger().WithError(validationError).Debug("is validation.Error")
		return convertValidationErrorToProtobuf(validationError, defaultErrorCode)
	}

	if validationErrors, ok := err.(validation.Errors); ok {

		log.GetLogger().WithError(validationErrors).Debug("is validation.Errors")
		return convertValidationErrorsToProtobuf(validationErrors, defaultErrorCode)
	}

	log.GetLogger().WithError(err).Debug("not a validation.Error")
	return &protos.Error{
		Code:    defaultErrorCode,
		Message: err.Error(),
	}
}

func convertValidationErrorsToProtobuf(validationErrors validation.Errors, defaultErrorCode int64) *protos.Error {
	if len(validationErrors) > 1 {

		return &protos.Error{
			Code:    defaultErrorCode,
			Message: validationErrors.Error(),
		}
	}

	for _, validationError := range validationErrors {

		if validationError, ok := validationError.(validation.Error); ok {

			return convertValidationErrorToProtobuf(validationError, defaultErrorCode)
		}
	}

	return &protos.Error{
		Code:    defaultErrorCode,
		Message: validationErrors.Error(),
	}
}

func convertValidationErrorToProtobuf(validationError validation.Error, defaultErrorCode int64) *protos.Error {

	code, convertError := strconv.ParseInt(validationError.Code(), 10, 64)
	if convertError != nil {
		return &protos.Error{
			Code:    defaultErrorCode,
			Message: validationError.Error(),
		}
	}
	return &protos.Error{
		Code:    code,
		Message: validationError.Error(),
	}
}

func ConvertMerchantStatusToProtobuf(status model.MerchantStatus) protos.MerchantStatus {

	switch status {
	case model.MerchantStatusOpen:
		return protos.MerchantStatus_MerchantStatusOpen
	case model.MerchantStatusClose:
		return protos.MerchantStatus_MerchantStatusClose
	}
	return protos.MerchantStatus_MerchantStatusOpen
}
