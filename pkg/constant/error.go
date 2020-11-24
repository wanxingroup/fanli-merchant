package constant

const (
	ErrorCodeMobileExist             = 411001
	ErrorCodePasswordLengthNotEnough = 411002
	ErrorCodePasswordLengthTooLarge  = 411003
	ErrorCodeMobileFormatInvalid     = 411004
	ErrorCodeShopNameTooLarge        = 411005
	ErrorCodeBusinessTimeError       = 411006
	ErrorCodeBusinessStatusError     = 411007
	ErrorCodeBusinessTypeError       = 411008
	ErrorCodeMerchantIdIsNotExist    = 411009
	ErrorCodeShopNameIsExist         = 411010
	ErrorCodeGetShopListError        = 411011
	ErrorCodeGetShopInfoError        = 411012
	ErrorCodeShopIsNotExist          = 411013
	ErrorCodeUserIdEmpty             = "411014"
	ErrorCodeCodeEmpty               = "411015"
	ErrorCodeRecordCodeError         = 411016

	ErrorCodeRegisterUserFailed     = 511001
	ErrorCodeValidatePasswordFailed = 511002
	ErrorCodeCreateShopFailed       = 511003
	ErrorCodeUpdateShopFailed       = 511004
	ErrorCodeCreateVerifierErr      = 511005
	ErrorCodeGetMerchantError       = 511006
	ErrorCodeUpdateMerchantError    = 511007
	ErrorCodeGetMerchantListError   = 511008
	ErrorCodeMerchantStatusError    = 511009
)

const (
	ErrorMessageMobileExist             = "相同的手机号已存在"
	ErrorMessagePasswordLengthNotEnough = "密码长度不足"
	ErrorMessagePasswordLengthTooLarge  = "密码长度过大"
	ErrorMessageMobileFormatInvalid     = "手机号格式不正确"
	ErrorMessageRegisterUserFailed      = "注册用户失败，内部服务暂时不可用"
	ErrorMessageValidatePasswordFailed  = "验证密码失败，内部服务暂时不可用"
	ErrorMessageShopNameTooLarge        = "商铺名字太长"
	ErrorMessageBusinessTimeError       = "商铺运营时间有误"
	ErrorMessageBusinessStatusError     = "商铺状态有误"
	ErrorMessageBusinessTypeError       = "门店类型有误"
	ErrorMessageMerchantIdIsNotExit     = "商家不存在"
	ErrorMessageCreateShopFailed        = "创建店铺失败，内部服务暂时不可用"
	ErrorMessageShopNameIsExist         = "相同的店铺名已存在"
	ErrorMessageGetShopListError        = "获取店铺列表错误"
	ErrorMessageGetShopInfoError        = "获取店铺详情错误"
	ErrorMessageShopIsNotExist          = "店铺不存在"
	ErrorMessageUserIdEmpty             = "用户 ID 为必填"
	ErrorMessageCodeEmpty               = "生成的验证码有误"
	ErrorMessageRecordCodeError         = "Code写入数据库有误"
	ErrorMessageUpdateShopFailed        = "更新店铺失败，内部服务暂时不可用"
	ErrorMessageCodeCreateVerifierErr   = "创建店铺成功, 但是创建核销员失败"
	ErrorMessageUpdateMerchantError     = "更新商家信息错误"
	ErrorMessageMerchantStatusError     = "商家已禁用"

	ErrorMessageGetRequestDataIsNil = "get request data is nil"
)

const (
	ErrorCodeShopIdEmpty    = "411017"
	ErrorMessageShopIdEmpty = "店铺 ID 为必填"
)

const (
	ErrorCodeNameEmpty    = "411018"
	ErrorMessageNameEmpty = "名称为必填"
)

const (
	ErrorCodeNameLengthOutOfRange    = "411019"
	ErrorMessageNameLengthOutOfRange = "名称长度超出限制"
)

const (
	ErrorCodeMobileEmpty    = "411020"
	ErrorMessageMobileEmpty = "手机号必填"
)

const (
	ErrorCodeMobileLengthOutOfRange    = "411021"
	ErrorMessageMobileLengthOutOfRange = "手机号长度超出限制"
)

const (
	ErrorCodeCreateStaffFailed    = 511008
	ErrorMessageCreateStaffFailed = "添加员工失败，内部服务暂时不可用"
)

const (
	ErrorCodeModifyStaffFailed    = 511009
	ErrorMessageModifyStaffFailed = "修改员工失败，内部服务暂时不可用"
)

const (
	ErrorCodeGetStaffFailed    = 511010
	ErrorMessageGetStaffFailed = "获取员工失败，内部服务暂时不可用"
)

const (
	ErrorCodeGetStaffListFailed    = 511011
	ErrorMessageGetStaffListFailed = "获取员工列表失败，内部服务暂时不可用"
)
