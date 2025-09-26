package model

import (
	"errors"

	"github.com/hkyangyi/newe/common/db"
)

// SysApiList API接口表，支持分组（父级），自动写入所有注册路由
// 父级分组如模块名，接口自动归类

type SysApiList struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`       // 主键
	ParentID  int64  `gorm:"type:bigint;index" json:"parent_id"`       // 父级分组ID
	GroupName string `gorm:"type:varchar(32);index" json:"group_name"` // 分组名（如 user、order）
	Name      string `gorm:"type:varchar(64)" json:"name"`             // 接口名称（如 用户列表）
	Route     string `gorm:"type:varchar(128);index" json:"route"`     // 路由模板（如 /api/v1/user/:id）
	Method    string `gorm:"type:varchar(10)" json:"method"`           // 请求方法
	Desc      string `gorm:"type:varchar(255)" json:"desc"`            // 接口说明
	CreatedAt int64  `json:"created_at"`                               // 创建时间
}

func (SysApiList) TableName() string { return "sys_api_list" }

// 执行注册数据库
func regSysApiListTable() error {
	if db.Db == nil {
		return errors.New("db not initialized")
	}
	return db.Db.AutoMigrate(&SysApiList{})
}

func (a *SysApiList) Add() error {
	err := db.Db.Create(a).Error
	return err
}

// AddIfNotExists: 若接口(route+method)不存在则插入
func (a *SysApiList) AddIfNotExists() error {
	var exist SysApiList
	err := db.Db.Where("route = ? AND method = ?", a.Route, a.Method).First(&exist).Error
	if err == nil && exist.ID > 0 {
		return nil // 已存在，不重复插入
	}
	return db.Db.Create(a).Error
}
