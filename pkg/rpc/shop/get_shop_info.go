package shop

import (
	rpclog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/jinzhu/gorm"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"

	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func (_ Controller) GetShopInfo(ctx context.Context, req *protos.GetShopInfoRequest) (reply *protos.GetShopInfoReply, err error) {

	var shop model.Shop
	logger := rpclog.WithRequestId(ctx, log.GetLogger())

	err = database.GetDB(constant.DatabaseConfigKey).Where(map[string]interface{}{"shopId": req.GetShopId()}).First(&shop).Error

	if err != nil {

		if !gorm.IsRecordNotFoundError(err) {

			logger.WithError(err).Error("get shopInfo record error")
			return &protos.GetShopInfoReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodeGetShopInfoError,
					Message: constant.ErrorMessageGetShopInfoError,
				},
			}, nil
		}

		logger.Info("shop not exist")
		return &protos.GetShopInfoReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeShopIsNotExist,
				Message: constant.ErrorMessageShopIsNotExist,
			},
		}, nil
	}

	// 查询核销人员信息
	var verifierMobiles []model.ShopVerification

	err = database.GetDB(constant.DatabaseConfigKey).Where(map[string]interface{}{"shopId": req.GetShopId()}).Find(&verifierMobiles).Error
	if err != nil {

		logger.WithError(err).Error("get verifier mobiles error")
		return &protos.GetShopInfoReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeGetShopInfoError,
				Message: constant.ErrorMessageGetShopInfoError,
			},
		}, nil
	}

	shopInfo := &protos.GetShopInfoReply{
		ShopInfo: &protos.ShopInfo{
			ShopId:             shop.ShopId,
			Name:               shop.Name,
			BusinessHoursStart: int64(shop.BusinessHoursStart),
			BusinessHoursEnd:   int64(shop.BusinessHoursEnd),
			ShopStatus:         ConvertShopStatusToProtobuf(shop.ShopStatus),
			ShopType:           int64(shop.ShopType),
			MerchantId:         shop.MerchantId,
			CreatedAt:          shop.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:          shop.UpdatedAt.Format("2006-01-02 15:04:05"),
			VerifierMobiles:    make([]string, 0, len(verifierMobiles)),
			FuYouMerchantId:    shop.FuYouMerchantId,
			IsRebate:           shop.IsRebate,
		},
	}

	for _, mobile := range verifierMobiles {
		shopInfo.ShopInfo.VerifierMobiles = append(shopInfo.ShopInfo.VerifierMobiles, mobile.Mobile)
	}

	return shopInfo, nil
}

func ConvertShopStatusToProtobuf(status model.ShopStatus) protos.ShopStatus {

	switch status {
	case model.ShopStatusOpen:
		return protos.ShopStatus_ShopStatusOpen
	case model.ShopStatusClose:
		return protos.ShopStatus_ShopStatusClose
	}
	return protos.ShopStatus_ShopStatusOpen
}
