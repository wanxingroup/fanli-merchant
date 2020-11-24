package merchant

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func TestController_RegisterUser(t *testing.T) {

	tests := []struct {
		input          *protos.RegisterUserRequest
		wantUserId     bool
		wantReplyError *protos.Error
		wantError      error
	}{
		{
			wantError:  fmt.Errorf("request data is nil"),
			wantUserId: false,
		},
		{
			input: &protos.RegisterUserRequest{
				Mobile:   "invalid",
				Password: "12345678",
			},
			wantUserId: false,
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodeMobileFormatInvalid,
				Message: constant.ErrorMessageMobileFormatInvalid,
			},
		},
		{
			input: &protos.RegisterUserRequest{
				Mobile:   "18916039393",
				Password: "12345678",
			},
			wantUserId: true,
		},
		{
			input: &protos.RegisterUserRequest{
				Mobile:   "18916039393",
				Password: "12345678",
			},
			wantUserId: false,
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodeMobileExist,
				Message: constant.ErrorMessageMobileExist,
			},
		},
		{
			input: &protos.RegisterUserRequest{
				Mobile:   "13800138000",
				Password: "short",
			},
			wantUserId: false,
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodePasswordLengthNotEnough,
				Message: constant.ErrorMessagePasswordLengthNotEnough,
			},
		},
		{
			input: &protos.RegisterUserRequest{
				Mobile:   "13800138000",
				Password: "the password is toooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo long",
			},
			wantUserId: false,
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodePasswordLengthTooLarge,
				Message: constant.ErrorMessagePasswordLengthTooLarge,
			},
		},
	}

	for _, test := range tests {

		c := Controller{}
		reply, err := c.RegisterUser(context.Background(), test.input)

		assert.Equal(t, test.wantError, err)
		if test.wantError != nil {

			assert.Nil(t, reply, test)
			continue
		}

		if reply == nil {
			assert.NotNil(t, reply, test)
			continue
		}

		if test.wantUserId {
			assert.Greater(t, reply.MerchantId, uint64(0), test)
		} else {
			assert.Equal(t, reply.MerchantId, uint64(0), "must be no return userId", test)
		}

		assert.Equal(t, test.wantReplyError, reply.Err, test)
	}
}
