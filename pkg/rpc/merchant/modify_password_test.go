package merchant

import (
	"fmt"
	"testing"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func TestController_ModifyPassword(t *testing.T) {

	tests := []struct {
		init      func() error
		input     *protos.ModifyPasswordRequest
		wantReply *protos.ModifyPasswordReply
		wantError error
		check     func() error
	}{
		{
			wantError: fmt.Errorf(constant.ErrorMessageGetRequestDataIsNil),
		},
		{
			input: &protos.ModifyPasswordRequest{
				MerchantId: 10001,
				Password:   "",
			},
			wantReply: &protos.ModifyPasswordReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodePasswordLengthNotEnough,
					Message: constant.ErrorMessagePasswordLengthNotEnough,
				},
			},
		},
		{
			input: &protos.ModifyPasswordRequest{
				MerchantId: 10001,
				Password:   "toooooooooooooooooooooooooooooooooooooooooooooooo long",
			},
			wantReply: &protos.ModifyPasswordReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodePasswordLengthTooLarge,
					Message: constant.ErrorMessagePasswordLengthTooLarge,
				},
			},
		},
		{
			input: &protos.ModifyPasswordRequest{
				MerchantId: 10001,
				Password:   "password",
			},
			wantReply: &protos.ModifyPasswordReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodeMerchantIdIsNotExist,
					Message: constant.ErrorMessageMerchantIdIsNotExit,
				},
			},
		},
		{
			init: func() error {

				merchant := &model.Merchant{
					MerchantId: 10004,
					Mobile:     "13344112200",
					Password:   "emptyPassword",
					Salt:       "emptySalt",
				}
				err := database.GetDB(constant.DatabaseConfigKey).Create(merchant).Error
				return err
			},
			input: &protos.ModifyPasswordRequest{
				MerchantId: 10004,
				Password:   "abcdefghijklmn",
			},
			wantReply: &protos.ModifyPasswordReply{
				Merchant: &protos.Merchant{
					MerchantId: 10004,
					Mobile:     "13344112200",
				},
			},
			check: func() error {
				merchant := &model.Merchant{}
				err := database.GetDB(constant.DatabaseConfigKey).Where(&model.Merchant{
					MerchantId: 10004,
				}).First(merchant).Error
				if err != nil {
					return err
				}

				assert.NotEqual(t, merchant.Password, "emptyPassword", merchant)
				assert.NotEqual(t, merchant.Salt, "emptySalt", merchant)
				return nil
			},
		},
	}

	for _, test := range tests {

		if test.init != nil {
			assert.Nil(t, test.init(), test)
		}

		controller := &Controller{}
		reply, err := controller.ModifyPassword(context.Background(), test.input)
		assert.Equal(t, test.wantError, err, test)
		assert.Equal(t, test.wantReply, reply, test)

		if test.check != nil {
			assert.Nil(t, test.check(), test)
		}
	}
}
