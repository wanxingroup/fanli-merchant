package merchant

import (
	"strings"

	rpcLog "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/rpc/utils/log"
	"github.com/shomali11/util/xstrings"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

func (controller *Controller) GetMerchantList(ctx context.Context, req *protos.GetMerchantListRequest) (*protos.GetMerchantListReply, error) {

	logger := rpcLog.WithRequestId(ctx, log.GetLogger()).WithField("req", req)
	var page = uint64(defaultPageNum)
	var pageSize = uint64(defaultPageSize)

	if req.GetPage() > 0 {
		page = req.GetPage()
	}

	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}

	conditions := map[string]interface{}{}

	if xstrings.IsNotBlank(req.GetMobileFuzzySearch()) {
		conditions["mobileFuzzySearch"] = strings.TrimSpace(req.GetMobileFuzzySearch())
	}

	if xstrings.IsNotBlank(req.GetNameFuzzySearch()) {
		conditions["nameFuzzySearch"] = strings.TrimSpace(req.GetNameFuzzySearch())
	}

	if xstrings.IsNotBlank(req.GetArea()) {
		conditions["area"] = strings.TrimSpace(req.GetArea())
	}

	if xstrings.IsNotBlank(req.GetManagementCentre()) {
		conditions["managementCentre"] = strings.TrimSpace(req.GetManagementCentre())
	}

	if xstrings.IsNotBlank(req.GetNetworkStationFuzzySearch()) {
		conditions["networkStationFuzzySearch"] = strings.TrimSpace(req.GetNetworkStationFuzzySearch())
	}

	pageData := map[string]uint64{
		"page":     page,
		"pageSize": pageSize,
	}

	list, count, err := getMerchantListByConditions(conditions, pageData)
	if err != nil {
		logger.WithError(err).Error("get merchant list failed")
		return &protos.GetMerchantListReply{
			Err: convertError(err, constant.ErrorCodeGetMerchantListError),
		}, nil
	}

	protoList := make([]*protos.Merchant, len(list))
	for key, item := range list {
		protoList[key] = &protos.Merchant{
			Mobile:           item.Mobile,
			MerchantId:       item.MerchantId,
			Name:             item.Name,
			Area:             item.Area,
			ManagementCentre: item.ManagementCentre,
			NetworkStation:   item.NetworkStation,
			Status:           ConvertMerchantStatusToProtobuf(item.Status),
			IsRebate:         item.IsRebate,
		}
	}

	return &protos.GetMerchantListReply{
		MerchantList: protoList,
		Count:        count,
	}, nil
}
