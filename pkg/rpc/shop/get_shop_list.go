package shop

import (
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	baseDb "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/utils/databases"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"

	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

const (
	DefaultPage     = 1  // 默认页码
	DefaultPageSize = 20 // 默认条数
)

func (_ Controller) GetShopList(ctx context.Context, req *protos.GetShopListRequest) (reply *protos.GetShopListReply, err error) {
	var page = uint64(DefaultPage)
	var pageSize = uint64(DefaultPageSize)

	if req.GetPage() > 0 {
		page = req.GetPage()
	}

	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}

	pageData := map[string]uint64{
		"page":     page,
		"pageSize": pageSize,
	}

	var shopList []model.Shop

	db := database.GetDB(constant.DatabaseConfigKey).Model(shopList).Where(map[string]interface{}{"merchantId": req.GetMerchantId()})

	if !validation.IsEmpty(req.GetShopName()) {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", req.GetShopName()))
	}

	results := make([]*model.Shop, 0, pageData["pageSize"])
	var count uint64
	err = baseDb.FindPage(db, pageData, &results, &count)

	if err != nil {

		log.GetLogger().WithError(err).Error("get shopList record error")

		return &protos.GetShopListReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeGetShopListError,
				Message: constant.ErrorMessageGetShopListError,
				Stack:   nil,
			},
		}, nil
	}

	var shopInfoList []*protos.ShopInfo
	for _, shop := range results {
		shopInfoList = append(shopInfoList, &protos.ShopInfo{
			ShopId:             shop.ShopId,
			Name:               shop.Name,
			BusinessHoursStart: int64(shop.BusinessHoursStart),
			BusinessHoursEnd:   int64(shop.BusinessHoursEnd),
			ShopStatus:         ConvertShopStatusToProtobuf(shop.ShopStatus),
			ShopType:           int64(shop.ShopType),
			MerchantId:         req.MerchantId,
			IsRebate:           shop.IsRebate,
			CreatedAt:          shop.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:          shop.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &protos.GetShopListReply{
		ShopList: shopInfoList,
		Count:    count,
	}, nil
}
