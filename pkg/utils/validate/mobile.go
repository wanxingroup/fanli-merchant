package validate

import (
	"regexp"
)

// 公式来源： https://github.com/VincentSit/ChinaMobilePhoneNumberRegex
const mobileRegularExpress = `^(?:\+?86)?1(?:3\d{3}|5[^4\D]\d{2}|8\d{3}|7(?:[235-8]\d{2}|4(?:0\d|1[0-2]|9\d))|9[0-35-9]\d{2}|66\d{2})\d{6}$`
const mvnoRegularExpress = `^(?:\+?86)?1(?:7[01]|6[257])\d{8}$`

func CheckMobile(mobile string) bool {

	var reg *regexp.Regexp
	reg = regexp.MustCompile(mobileRegularExpress)
	if reg.MatchString(mobile) {
		return true
	}

	reg = regexp.MustCompile(mvnoRegularExpress)
	if reg.MatchString(mobile) {
		return true
	}

	return false
}
