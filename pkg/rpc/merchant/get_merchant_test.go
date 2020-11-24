package merchant

import (
	"testing"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/request/requestid"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func TestController_GetMerchant(t *testing.T) {

	tests := []struct {
		initData func()
		context  context.Context
		input    *protos.GetMerchantRequest
		want     *protos.GetMerchantReply
		err      error
	}{
		{
			initData: func() {
				database.GetDB(constant.DatabaseConfigKey).Create(&model.Merchant{
					MerchantId: 30339,
					Mobile:     "18916000000",
				})
			},
			context: context.WithValue(context.Background(), requestid.Key, "test"),
			input:   &protos.GetMerchantRequest{MerchantId: 30339},
			want: &protos.GetMerchantReply{
				Merchant: &protos.Merchant{
					MerchantId:       30339,
					Mobile:           "18916000000",
					Name:             "",
					Area:             "",
					ManagementCentre: "",
					NetworkStation:   "",
					Status:           1,
				},
			},
		},
		{
			context: context.WithValue(context.Background(), requestid.Key, "test2"),
			input:   &protos.GetMerchantRequest{MerchantId: 404},
			want: &protos.GetMerchantReply{
				Err: &protos.Error{
					Code:    constant.ErrorCodeMerchantIdIsNotExist,
					Message: constant.ErrorMessageMerchantIdIsNotExit,
				},
			},
		},
	}

	for _, test := range tests {

		if test.initData != nil {
			test.initData()
		}

		controller := &Controller{}
		reply, err := controller.GetMerchant(test.context, test.input)
		assert.Equal(t, test.err, err, test)
		if err != nil {
			continue
		}

		assert.Equal(t, test.want, reply, test)
	}
}
