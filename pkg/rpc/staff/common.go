package staff

import (
	"strconv"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

const (
	DefaultPage     = 1   // 默认页码
	DefaultPageSize = 20  // 默认条数
	MaxPageSize     = 100 // 最大条数
)

func ConvertModelToProtobuf(staff *model.Staff) *protos.StaffInfo {
	return &protos.StaffInfo{
		StaffId:  staff.StaffId,
		ShopId:   staff.ShopId,
		Name:     staff.Name,
		Mobile:   staff.Mobile,
		IsRebate: staff.IsRebate,
		Status:   ConvertStaffStatusToProtobuf(staff.Status),
	}
}

func ConvertStaffStatusToProtobuf(status model.StaffStatus) protos.StaffStatus {

	switch status {
	case model.StaffStatusOpen:
		return protos.StaffStatus_StaffStatusOpen
	case model.StaffStatusClose:
		return protos.StaffStatus_StaffStatusClose
	}
	return protos.StaffStatus_StaffStatusOpen
}

func ConvertErrorToProtobuf(err error) *protos.Error {

	if validationError, ok := err.(validation.Error); ok {
		errorCode, convertError := strconv.Atoi(validationError.Code())
		if convertError != nil {
			errorCode = errors.CodeServerInternalError
		}
		return &protos.Error{
			Code:    int64(errorCode),
			Message: validationError.Error(),
		}
	}

	return &protos.Error{
		Code:    errors.CodeServerInternalError,
		Message: err.Error(),
	}
}
