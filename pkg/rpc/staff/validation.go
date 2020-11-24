package staff

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
)

var ShopIdRule = []validation.Rule{
	validation.Required.ErrorObject(
		validation.NewError(constant.ErrorCodeShopIdEmpty, constant.ErrorMessageShopIdEmpty),
	),
}

var UserIdRule = []validation.Rule{
	validation.Required.ErrorObject(
		validation.NewError(constant.ErrorCodeUserIdEmpty, constant.ErrorMessageUserIdEmpty),
	),
}

var NameRule = []validation.Rule{
	validation.Required.ErrorObject(validation.NewError(constant.ErrorCodeNameEmpty, constant.ErrorMessageNameEmpty)),
	validation.RuneLength(2, 40).ErrorObject(validation.NewError(constant.ErrorCodeNameLengthOutOfRange, constant.ErrorMessageNameLengthOutOfRange)),
}

var MobileRule = []validation.Rule{
	validation.Required.ErrorObject(validation.NewError(constant.ErrorCodeMobileEmpty, constant.ErrorMessageMobileEmpty)),
	validation.RuneLength(11, 11).ErrorObject(validation.NewError(constant.ErrorCodeMobileLengthOutOfRange, constant.ErrorMessageMobileLengthOutOfRange)),
}
