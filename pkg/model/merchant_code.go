package model

import (
	"time"

	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const (
	TableNameMerchantCode = "merchant_code"
	CodeTypeFindPassword  = 0
)

type MerchantCode struct {
	Id        uint64    `gorm:"column:id;type:bigint unsigned;primary_key;comment:'主键 Id'"`
	Mobile    string    `gorm:"column:mobile;type:char(11);not null;default:'';index:mobile;comment:'手机号'"`
	Code      string    `gorm:"column:code;type:char(6);not null;default:'';index:mobile;comment:'验证码'"`
	ExpiredAt time.Time `gorm:"column:expiredAt;not null;comment:'过期时间'"`
	CodeType  uint8     `gorm:"column:codeType;type:tinyint unsigned;not null;default:'0';comment:'验证码类型：0 => 找回密码'"`
	databases.Time
}

func (shop *MerchantCode) TableName() string {

	return TableNameMerchantCode
}
