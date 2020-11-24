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

	baseDb "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/databases"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func (_ Controller) GetList(ctx context.Context, req *protos.GetStaffListRequest) (*protos.GetStaffListReply, error) {

	logger := rpcLog.WithRequestId(ctx, log.GetLogger())

	if req == nil {
		logger.Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	var page = uint64(DefaultPage)
	var pageSize = uint64(DefaultPageSize)

	if req.GetPage() > 0 {
		page = req.GetPage()
	}

	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	pageData := map[string]uint64{
		"page":     page,
		"pageSize": pageSize,
	}

	err := validateGetStaffList(req)
	if err != nil {
		return &protos.GetStaffListReply{
			Err: ConvertErrorToProtobuf(err),
		}, nil
	}

	var staffList []*model.Staff

	db := database.GetDB(constant.DatabaseConfigKey).Model(staffList).Where("shopId = ?", req.GetShopId())

	if !validation.IsEmpty(req.GetUserName()) {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", req.GetUserName()))
	}

	results := make([]*model.Staff, 0, pageData["pageSize"])
	var count uint64
	err = baseDb.FindPage(db, pageData, &results, &count)

	if err != nil {
		logger.WithError(err).Error("get staff list error")
		return &protos.GetStaffListReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeGetStaffListFailed,
				Message: constant.ErrorMessageGetStaffListFailed,
				Stack:   nil,
			},
		}, nil
	}

	staffInfoList := make([]*protos.StaffInfo, 0, len(staffList))
	for _, staffModel := range results {
		staffInfoList = append(staffInfoList, ConvertModelToProtobuf(staffModel))
	}

	return &protos.GetStaffListReply{
		StaffList: staffInfoList,
		Count:     count,
	}, nil
}

func validateGetStaffList(req *protos.GetStaffListRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.ShopId, ShopIdRule...),
	)
}
