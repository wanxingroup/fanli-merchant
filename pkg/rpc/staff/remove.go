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

func (_ Controller) Remove(ctx context.Context, req *protos.RemoveStaffRequest) (*protos.RemoveStaffReply, error) {

	logger := rpcLog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	err := validateRemoveStaff(req)
	if err != nil {
		return &protos.RemoveStaffReply{
			Err: ConvertErrorToProtobuf(err),
		}, nil
	}

	result, err := removeStaff(req)
	if err != nil {
		logger.WithError(err).Error("remove staff error")
		return &protos.RemoveStaffReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeCreateStaffFailed,
				Message: constant.ErrorMessageCreateStaffFailed,
				Stack:   nil,
			},
		}, nil
	}

	return &protos.RemoveStaffReply{
		Result: result,
	}, nil
}

func removeStaff(req *protos.RemoveStaffRequest) (bool, error) {

	record := &model.Staff{
		StaffId: req.StaffId,
	}

	err := database.GetDB(constant.DatabaseConfigKey).Delete(record).Error
	if err != nil {

		log.GetLogger().WithField("staff", record).WithError(err).Error("remove staff record error")
		return false, err
	}

	return true, nil
}

func validateRemoveStaff(req *protos.RemoveStaffRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.StaffId, UserIdRule...),
	)
}
