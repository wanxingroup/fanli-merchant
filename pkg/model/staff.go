package model

import (
	databases "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database/models"
)

const TableNameStaff = "staff"

type Staff struct {
	StaffId  uint64      `gorm:"column:staffId;type:bigint unsigned;primary_key;comment:'员工ID'"`
	ShopId   uint64      `gorm:"column:shopId;type:bigint unsigned;index:shopId_idx;comment:'店铺ID'"`
	Name     string      `gorm:"column:name;type:char(20);not null;default:'';comment:'姓名'"`
	Mobile   string      `gorm:"column:mobile;type:varchar(11);unique_index:idx_mobile;not null;default:'';comment:'手机号'"`
	IsRebate bool        `gorm:"column:isRebate;type:tinyint unsigned;not null;default:'0';comment:'是否返利'"`
	Status   StaffStatus `gorm:"column:status;type:tinyint unsigned;not null;default:'1';comment:'用户状态： 1 => 正常；2 => 禁用'"`
	databases.BasicTimeFields
}

func (Staff *Staff) TableName() string {
	return TableNameStaff
}

type StaffStatus uint8

const (
	StaffStatusOpen  StaffStatus = 1 //正常
	StaffStatusClose StaffStatus = 2 //禁用
)
