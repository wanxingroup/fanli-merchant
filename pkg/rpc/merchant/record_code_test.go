package merchant

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/rpc/protos"
)

func TestController_RecordCode(t *testing.T) {
	tests := []struct {
		input     *protos.RecordCodeRequest
		wantError error
	}{
		{
			input: &protos.RecordCodeRequest{
				Mobile:    "13167185999",
				Code:      "123456",
				ExpiredAt: "2018-12-02 13:14:15",
			},
		},
	}

	for _, test := range tests {

		c := Controller{}
		reply, err := c.RecordCode(context.Background(), test.input)

		assert.Equal(t, test.wantError, err)
		assert.Nil(t, err)
		assert.Nil(t, reply.Err)

	}
}
