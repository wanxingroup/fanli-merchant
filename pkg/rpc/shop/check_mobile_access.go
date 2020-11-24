package shop

import (
	"fmt"

	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

func (_ Controller) CheckMobileAccess(ctx context.Context, req *protos.CheckMobileAccessRequest) (reply *protos.CheckMobileAccessReply, err error) {

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	isVerifier := CheckVerifier(req.GetMobile(), req.GetShopId())

	return &protos.CheckMobileAccessReply{
		IsVerifier: isVerifier,
	}, nil
}
