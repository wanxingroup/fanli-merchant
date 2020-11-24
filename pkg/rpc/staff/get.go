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

func (_ Controller) Get(ctx context.Context, req *protos.GetStaffInfoRequest) (*protos.GetStaffInfoReply, error) {

	logger := rpcLog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	err := validateGetStaff(req)
	if err != nil {
		return &protos.GetStaffInfoReply{
			Err: ConvertErrorToProtobuf(err),
		}, nil
	}

	staff, err := getStaff(req.StaffId)
	if err != nil {
		logger.WithError(err).Error("get staff error")
		return &protos.GetStaffInfoReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeGetStaffFailed,
				Message: constant.ErrorMessageGetStaffFailed,
				Stack:   nil,
			},
		}, nil
	}

	return &protos.GetStaffInfoReply{
		StaffInfo: ConvertModelToProtobuf(staff),
	}, nil
}

func validateGetStaff(req *protos.GetStaffInfoRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.StaffId, UserIdRule...),
	)
}

func getStaff(staffId uint64) (*model.Staff, error) {
	staff := model.Staff{}
	err := database.GetDB(constant.DatabaseConfigKey).First(&staff, staffId).Error
	if err != nil {
		log.GetLogger().WithField("staff", staff).WithError(err).Error("get staff record error")
		return nil, err
	}

	return &staff, nil
}
