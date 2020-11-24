package merchant

import (
	"fmt"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	context "golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/log"
)

func (_ Controller) ValidatePassword(ctx context.Context, req *protos.ValidatePasswordRequest) (reply *protos.ValidatePasswordReply, err error) {

	if req == nil {
		log.GetLogger().Error("request data is nil")
		return nil, fmt.Errorf("request data is nil")
	}

	if req.Password == "" || req.Mobile == "" {
		return &protos.ValidatePasswordReply{
			Pass: false,
		}, nil
	}

	merchant, err := getMerchantByMobile(req.Mobile)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &protos.ValidatePasswordReply{
				Pass: false,
			}, nil
		}

		log.GetLogger().WithError(err).Error("get merchant data error")

		return &protos.ValidatePasswordReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeValidatePasswordFailed,
				Message: errors.Wrap(err, constant.ErrorMessageValidatePasswordFailed).Error(),
			},
			Pass: false,
		}, nil
	}

	if merchant.Status != model.MerchantStatusOpen {
		log.GetLogger().Error("merchant status data error")

		return &protos.ValidatePasswordReply{
			Err: &protos.Error{
				Code:    constant.ErrorCodeMerchantStatusError,
				Message: constant.ErrorMessageMerchantStatusError,
			},
			Pass: false,
		}, nil
	}

	if !password.Verify(req.Password, merchant.Salt, merchant.Password, getPasswordOptions()) {

		return &protos.ValidatePasswordReply{
			Pass: false,
		}, nil
	}

	return &protos.ValidatePasswordReply{
		Pass:       true,
		MerchantId: merchant.MerchantId,
	}, nil
}
