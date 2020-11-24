package merchant

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"dev-gitlab.wanxingrowth.com/fanli/merchant/pkg/constant"
)

var MerchantIdRule = []validation.Rule{
	validation.Required.ErrorObject(validation.NewError(constant.ErrorCodeUserIdEmpty, constant.ErrorMessageUserIdEmpty)),
}

var CodeRules = []validation.Rule{
	validation.Required.ErrorObject(validation.NewError(constant.ErrorCodeCodeEmpty, constant.ErrorMessageCodeEmpty)),
}

var PasswordRules = []validation.Rule{
	validation.Required.ErrorObject(validation.NewError(strconv.FormatUint(constant.ErrorCodePasswordLengthNotEnough, 10), constant.ErrorMessagePasswordLengthNotEnough)),
	validation.Length(passwordMinLength, 0).ErrorObject(validation.NewError(strconv.FormatUint(constant.ErrorCodePasswordLengthNotEnough, 10), constant.ErrorMessagePasswordLengthNotEnough)),
	validation.Length(0, passwordMaxLength).ErrorObject(validation.NewError(strconv.FormatUint(constant.ErrorCodePasswordLengthTooLarge, 10), constant.ErrorMessagePasswordLengthTooLarge)),
}
