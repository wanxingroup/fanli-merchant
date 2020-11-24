package shop_test

import (
	"os"
	"testing"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/utils/test"
)

func TestMain(m *testing.M) {

	test.Init()

	code := m.Run()

	if code == 0 {
		test.Release()
	}

	os.Exit(code)
}
