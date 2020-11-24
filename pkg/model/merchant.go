package model

import (
	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameMerchant = "merchant"

type Merchant struct {
	MerchantId       uint64         `gorm:"column:merchantId;type:bigint unsigned;primary_key;comment:'商家用户ID'"`
	Mobile           string         `gorm:"column:mobile;type:varchar(11);not null;default:'';comment:'商家手机号'"`
	Password         string         `gorm:"column:password;type:char(100);not null;default:'';comment:'密码'"`
	Salt             string         `gorm:"column:salt;type:char(50);not null;default:'';comment:'密码的盐'"`
	IsRebate         bool           `gorm:"column:isRebate;type:tinyint unsigned;not null;default:'0';comment:'是否返利'"`
	Name             string         `gorm:"column:name;type:varchar(40);not null;default:'';comment:'名称'"`
	Area             string         `gorm:"column:area;type:varchar(40);not null;default:'';comment:'区域'"`
	ManagementCentre string         `gorm:"column:managementCentre;type:varchar(40);not null;default:'';comment:'管理中心'"`
	NetworkStation   string         `gorm:"column:networkStation;type:varchar(40);not null;default:'';comment:'网点'"`
	Status           MerchantStatus `gorm:"column:status;type:tinyint unsigned;not null;default:'1';comment:'状态： 1 => 正常；2 => 禁用'"`

	databases.Time
}

func (merchant *Merchant) TableName() string {
	return TableNameMerchant
}

type MerchantStatus uint8

const (
	MerchantStatusOpen  MerchantStatus = 1 // 正常
	MerchantStatusClose MerchantStatus = 2 // 禁用
)
