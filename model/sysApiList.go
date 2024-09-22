package model

import (
	"errors"
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

type SysApiList struct {
	ID         string `gorm:"primary_key" json:"id" form:"id"`
	MenuId     string `json:"menuId" form:"menuId"  dict:"MenuName_sysMenus" `
	MenuName   string `gorm:"-" json:"menuName"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Name       string `json:"name"`
	CreateTime int64  `json:"createTime"` //创建时间
	UpdateTime int64  `json:"updateTime"` //修改时间
	Status     int    `json:"status"`     //角色状态(-1禁用，1启用)
	utils.PageList
}

// 添加
func (a *SysApiList) Add() error {
	a.ID = utils.GetUUID()
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = time.Now().Unix()

	err := db.Db.Create(a).Error
	return err
}

// 删除
func (a *SysApiList) Del() error {
	err := db.Db.Model(a).Delete(a).Error
	return err
}

// 编辑
func (a *SysApiList) Edit() error {
	if len(a.ID) != 32 {
		return errors.New("缺少参数ID")
	}

	a.UpdateTime = time.Now().Unix()
	err := db.Db.Model(a).Updates(a).Error
	return err
}

// 启停
func (a *SysApiList) SetStatus() error {

	a.UpdateTime = time.Now().Unix()
	err := db.Db.Model(a).Update("status", a.Status).Error
	return err
}

// 获取列表
func (a *SysApiList) GetList(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []SysApiList
	db.Db.Model(&SysApiList{}).Where(where, v...).Count(&page.Total).Order("create_time desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	page.List = items
	return page
}

// 根据path获取数据
func (a *SysApiList) GetByPath() {
	db.Db.Table("sys_api_list").Where("path = ?", a.Path).First(a)
}
