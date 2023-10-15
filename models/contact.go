package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	OwnerId  uint // 谁的对应关系
	TargetId uint // 对应谁
	Type     int  // 对应的类型 0 1 2
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}
