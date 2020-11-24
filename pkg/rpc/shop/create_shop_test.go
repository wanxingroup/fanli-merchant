package shop

import (
	"context"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/model"
	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const MerchantId_Test = 10086
const ShopId_Test = 100

func TestController_CreateShop(t *testing.T) {

	addTestData()

	tests := []struct {
		input          *protos.CreateShopRequest
		wantReplyError *protos.Error
		wantError      error
	}{
		{
			wantError: fmt.Errorf("request data is nil"),
		},
		{
			input: &protos.CreateShopRequest{
				Name: "我的名字超长我的名字超长我的名字超长我的名字超长",
			},
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodeShopNameTooLarge,
				Message: constant.ErrorMessageShopNameTooLarge,
				Stack:   nil,
			},
		},
		{
			input: &protos.CreateShopRequest{
				Name:               "我的店铺名字正正好好一共就二十个字符整数",
				BusinessHoursStart: -5,
				BusinessHoursEnd:   999,
			},
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodeBusinessTimeError,
				Message: constant.ErrorMessageBusinessTimeError,
				Stack:   nil,
			},
		},
		{
			input: &protos.CreateShopRequest{
				Name:               "我的店铺名字正正好好一共就二十个字符整数",
				BusinessHoursStart: 998,
				BusinessHoursEnd:   172801,
			},
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodeBusinessTimeError,
				Message: constant.ErrorMessageBusinessTimeError,
				Stack:   nil,
			},
		},
		{
			input: &protos.CreateShopRequest{
				Name:               "我的店铺名字正正好好一共就二十个字符整数",
				BusinessHoursStart: 998,
				BusinessHoursEnd:   172800,
				ShopStatus:         5,
			},
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodeBusinessStatusError,
				Message: constant.ErrorMessageBusinessStatusError,
				Stack:   nil,
			},
		},
		{
			input: &protos.CreateShopRequest{
				Name:               "我的店铺名字正正好好一共就二十个字符整数",
				BusinessHoursStart: 998,
				BusinessHoursEnd:   172800,
				ShopStatus:         1,
				ShopType:           9,
			},
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodeBusinessTypeError,
				Message: constant.ErrorMessageBusinessTypeError,
				Stack:   nil,
			},
		},
		{
			input: &protos.CreateShopRequest{
				Name:               "我的店铺名字正正好好一共就二十个字符整数",
				BusinessHoursStart: 998,
				BusinessHoursEnd:   172800,
				ShopStatus:         1,
				ShopType:           1,
				MerchantId:         336,
			},
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodeMerchantIdIsNotExist,
				Message: constant.ErrorMessageMerchantIdIsNotExit,
				Stack:   nil,
			},
		},
		{
			input: &protos.CreateShopRequest{
				Name:               "我叫张三的店铺",
				BusinessHoursStart: 998,
				BusinessHoursEnd:   172800,
				ShopStatus:         1,
				ShopType:           1,
				MerchantId:         MerchantId_Test,
			},
			wantReplyError: &protos.Error{
				Code:    constant.ErrorCodeShopNameIsExist,
				Message: constant.ErrorMessageShopNameIsExist,
				Stack:   nil,
			},
		},
		{
			input: &protos.CreateShopRequest{
				Name:               "我叫张三的店2",
				BusinessHoursStart: 998,
				BusinessHoursEnd:   172800,
				ShopStatus:         1,
				ShopType:           1,
				MerchantId:         MerchantId_Test,
			},
		},
	}

	for _, test := range tests {

		c := Controller{}
		reply, err := c.CreateShop(context.Background(), test.input)

		assert.Equal(t, test.wantError, err)
		if test.wantError != nil {

			assert.Nil(t, reply, test)
			continue
		}

		if reply == nil {
			assert.NotNil(t, reply, test)
			continue
		}

		assert.Equal(t, test.wantReplyError, reply.Err, test)
	}
}

func addTestData() {
	// 写入 测试数据
	record := &model.Shop{
		ShopId:             ShopId_Test,
		MerchantId:         MerchantId_Test,
		Name:               "我叫张三的店铺",
		BusinessHoursStart: 0,
		BusinessHoursEnd:   86400,
		ShopStatus:         1,
		ShopType:           1,
	}
	database.GetDB(constant.DatabaseConfigKey).Create(record)

	record2 := &model.Merchant{
		MerchantId: MerchantId_Test,
		Mobile:     "15214356383",
		Password:   "abcdefg",
		Salt:       "sst",
	}

	database.GetDB(constant.DatabaseConfigKey).Create(record2)

}
