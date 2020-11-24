package model

import (
	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameShop = "shop"

type Shop struct {
	ShopId             uint64     `gorm:"column:shopId;type:bigint unsigned;primary_key;comment:'店铺 ID'"`
	MerchantId         uint64     `gorm:"column:merchantId;type:bigint unsigned;not null;comment:'商家用户ID'"`
	Name               string     `gorm:"column:name;type:varchar(40);not null;default:'0';comment:'店铺名称'"`
	BusinessHoursStart uint32     `gorm:"column:businessHoursStart;type:int unsigned;not null;default:'0';comment:'开始营业时间，区间左边'"`
	BusinessHoursEnd   uint32     `gorm:"column:businessHoursEnd;type:int unsigned;not null;default:'86400';comment:'结束营业时间，区间右边'"`
	ShopStatus         ShopStatus `gorm:"column:shopStatus;type:tinyint unsigned;not null;default:'1';comment:'店铺营业状态： 1 => 营业中；2 => 停业整顿'"`
	ShopType           ShopType   `gorm:"column:shopType;type:tinyint unsigned;not null;default:'1';comment:'门店类型： 1 => 直营店；2 => 加盟店'"`
	FuYouMerchantId    string     `gorm:"column:fuYouMerchantId;type:char(20);not null;default:'';comment:'门店富有支付号'"`
	IsRebate           bool       `gorm:"column:isRebate;type:tinyint unsigned;not null;default:'0';comment:'是否返利'"`
	databases.Time
}

type ShopStatus uint8
type ShopType uint8

const (
	ShopStatusOpen  ShopStatus = 1
	ShopStatusClose ShopStatus = 2
)

const (
	ShopTypeDirect ShopType = 1
	ShopTypeLeague ShopType = 2
)

func (shop *Shop) TableName() string {

	return TableNameShop
}
