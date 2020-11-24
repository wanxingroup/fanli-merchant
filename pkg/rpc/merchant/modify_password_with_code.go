package merchant

import (
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func (controller *Controller) ModifyPasswordWithCode(ctx context.Context, req *protos.ModifyPasswordWithCodeRequest) (*protos.ModifyPasswordWithCodeReply, error) {
	// todo 修改密码
	return &protos.ModifyPasswordWithCodeReply{}, nil
}

func (controller *Controller) validateModifyPasswordWithCodeRequestData(req *protos.ModifyPasswordWithCodeRequest) error {
	return nil
}
