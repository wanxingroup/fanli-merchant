package staff

import (
	"fmt"

	rpcLog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func (_ Controller) Modify(ctx context.Context, req *protos.ModifyStaffInfoRequest) (*protos.ModifyStaffInfoReply, error) {

	logger := rpcLog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	err := validateModifyStaff(req)
	if err != nil {
		return &protos.ModifyStaffInfoReply{
			Err: ConvertErrorToProtobuf(err),
		}, nil
	}

	err = modifyStaff(req)
	if err != nil {
		logger.WithError(err).Error("modify staff error")
		return &protos.ModifyStaffInfoReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeModifyStaffFailed,
				Message: constant.ErrorMessageModifyStaffFailed,
				Stack:   nil,
			},
		}, nil
	}

	staff, err := getStaff(req.StaffId)
	if err != nil {
		logger.WithError(err).Error("get staff error")
		return &protos.ModifyStaffInfoReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeGetStaffFailed,
				Message: constant.ErrorMessageGetStaffFailed,
				Stack:   nil,
			},
		}, nil
	}

	return &protos.ModifyStaffInfoReply{
		StaffInfo: ConvertModelToProtobuf(staff),
	}, nil
}

func modifyStaff(req *protos.ModifyStaffInfoRequest) error {

	record := map[string]interface{}{
		"shopId":   req.ShopId,
		"name":     req.Name,
		"mobile":   req.Mobile,
		"isRebate": req.IsRebate,
		"status":   req.Status,
	}

	err := database.GetDB(constant.DatabaseConfigKey).Table(model.TableNameStaff).Where("staffId = ?", req.StaffId).Updates(record).Error

	if err != nil {
		log.GetLogger().WithField("record", record).WithError(err).Error("modify staff record error")
		return err
	}

	log.GetLogger().WithField("record", record).Info("modify staff success")

	return nil
}

func validateModifyStaff(req *protos.ModifyStaffInfoRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.ShopId, ShopIdRule...),
		validation.Field(&req.Name, NameRule...),
		validation.Field(&req.Mobile, MobileRule...),
	)
}
