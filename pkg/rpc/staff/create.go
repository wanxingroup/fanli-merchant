package staff

import (
	"fmt"

	rpcLog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	idCreator "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/idcreator/snowflake"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func (_ Controller) Create(ctx context.Context, req *protos.CreateStaffInfoRequest) (*protos.CreateStaffInfoReply, error) {

	logger := rpcLog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	err := validateCreateStaff(req)
	if err != nil {
		return &protos.CreateStaffInfoReply{
			Err: ConvertErrorToProtobuf(err),
		}, nil
	}

	staffInfo, err := createStaff(req)
	if err != nil {
		logger.WithError(err).Error("create staff error")
		return &protos.CreateStaffInfoReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeCreateStaffFailed,
				Message: constant.ErrorMessageCreateStaffFailed,
				Stack:   nil,
			},
		}, nil
	}

	return &protos.CreateStaffInfoReply{
		StaffInfo: staffInfo,
	}, nil
}

func createStaff(req *protos.CreateStaffInfoRequest) (*protos.StaffInfo, error) {

	record := &model.Staff{
		StaffId:  idCreator.NextID(),
		ShopId:   req.ShopId,
		Name:     req.Name,
		Mobile:   req.Mobile,
		IsRebate: req.IsRebate,
		Status:   model.StaffStatus(req.Status),
	}

	err := database.GetDB(constant.DatabaseConfigKey).Create(record).Error
	if err != nil {

		log.GetLogger().WithField("staff", record).WithError(err).Error("create staff record error")
		return nil, err
	}

	return ConvertModelToProtobuf(record), nil
}

func validateCreateStaff(req *protos.CreateStaffInfoRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.ShopId, ShopIdRule...),
		validation.Field(&req.Name, NameRule...),
		validation.Field(&req.Mobile, MobileRule...),
	)
}
