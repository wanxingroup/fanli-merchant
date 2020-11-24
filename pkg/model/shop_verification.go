package model

import (
	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameShopVerification = "shop_verification"

type ShopVerification struct {
	ShopId uint64 `gorm:"column:shopId;type:bigint unsigned;primary_key;comment:'店铺 ID'"`
	Mobile string `gorm:"column:mobile;type:varchar(11);not null;default:'';primary_key;comment:'核销用户的手机号'"`
	databases.Time
}

func (shop *ShopVerification) TableName() string {
	return TableNameShopVerification
}
