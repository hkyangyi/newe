package model

import (
	"errors"
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

type DevPage struct {
	Id             string `json:"id" form:"id"  gorm:"primary_key"`
	PageName       string `json:"pageName" form:"pageName" `             // 页面名称
	RouteRoot      string `json:"routeRoot" form:"routeRoot" `           // 主路由路径
	PageType       int    `json:"pageType" form:"pageType" `             // 1 单表，2树形单表，3 1对1管理
	SqlTable       string `json:"sqlTable" form:"sqlTable" `             // 数据表主表
	SqlTableDeputy string `json:"sqlTableDeputy" form:"sqlTableDeputy" ` // 数据表副表
	SqlKey         string `json:"sqlKey" form:"sqlKey" `                 // 关联副表的ID
	CreateTime     int64  `json:"createTime" form:"createTime" `
	utils.PageList
}

func (a *DevPage) Refresh() error {
	if len(a.Id) != 32 {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Model(a).First(a).Error
	return err
}

// 新增
func (a *DevPage) Add() error {
	a.Id = utils.GetUUID()
	a.CreateTime = time.Now().Unix()
	err := db.Db.Create(a).Error
	return err
}

// 编辑
func (a *DevPage) Edit() error {
	if len(a.Id) != 32 {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Updates(a).Error
	return err
}

// 删除
func (a *DevPage) Del() error {
	if len(a.Id) != 32 {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Model(a).Delete(a).Error
	return err
}

// 获取分页
func (a *DevPage) GetPage(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []DevPage
	db.Db.Model(&DevPage{}).Where(where, v...).Count(&page.Total).Order("create_time desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	page.List = items
	return page
}

// 获取列表
func (a *DevPage) GetList() []DevPage {
	var items []DevPage
	db.Db.Model(&DevPage{}).Order("create_time desc").Find(&items)
	return items
}

// 根据数据表查询
func (a *DevPage) GetOneByTable() (DevPage, error) {
	var one DevPage
	err := db.Db.Where("sql_table=? ", a.SqlTable).First(&one).Error
	return one, err
}
