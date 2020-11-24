package merchant

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func TestController_ValidatePassword(t *testing.T) {

	registerUser("validatePas", "correct password", "", "", "", "", false, model.MerchantStatusOpen)
	tests := []struct {
		input          *protos.ValidatePasswordRequest
		wantPass       bool
		wantError      error
		wantReplyError *protos.Error
	}{
		{
			wantError: fmt.Errorf("request data is nil"),
		},
		{
			input: &protos.ValidatePasswordRequest{
				Mobile:   "13800138000",
				Password: "",
			},
			wantPass: false,
		},
		{
			input: &protos.ValidatePasswordRequest{
				Mobile:   "",
				Password: "13800138000",
			},
			wantPass: false,
		},
		{
			input: &protos.ValidatePasswordRequest{
				Mobile:   "not found",
				Password: "not found",
			},
			wantPass: false,
		},
		{
			input: &protos.ValidatePasswordRequest{
				Mobile:   "validatePas",
				Password: "wrong password",
			},
			wantPass: false,
		},
		{
			input: &protos.ValidatePasswordRequest{
				Mobile:   "validatePas",
				Password: "correct password",
			},
			wantPass: true,
		},
	}

	for _, test := range tests {

		c := Controller{}
		reply, err := c.ValidatePassword(context.Background(), test.input)

		assert.Equal(t, test.wantError, err)
		if test.wantError != nil {

			assert.Nil(t, reply, test)
			continue
		}

		if reply == nil {
			assert.NotNil(t, reply, test)
			continue
		}

		assert.Equal(t, test.wantPass, reply.Pass, test)
	}
}
