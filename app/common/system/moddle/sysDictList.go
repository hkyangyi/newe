package moddle

import (
	"errors"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

type SysDictList struct {
	ID         string `gorm:"primary_key" json:"id"  form:"id"` //
	ParentId   string `json:"parentId" form:"parentId"`         //
	ParentName string `json:"parentName" form:"code"`           //字典名称
	Value      string `json:"value"`                            //值
	Name       string `json:"name"`                             //名称
	Status     int    `json:"status"`                           //1正常，2停用，10已删除
	Sort       int    `json:"sort"`                             //排序
	utils.PageList
}

// 添加
func (a *SysDictList) Add() error {
	a.ID = utils.GetUUID()
	err := db.Db.Create(a).Error
	return err
}

// 编辑
func (a *SysDictList) Edit() error {
	if len(a.ID) != 32 {
		return errors.New("缺乏参数ID")
	}
	err := db.Db.Model(a).Updates(a).Error
	return err
}

// 获取列表
func (a *SysDictList) GetList(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []SysDictList
	db.Db.Model(&SysDictList{}).Where(where, v...).Count(&page.Total).Order("sort asc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	page.List = items
	return page
}

// 删除
func (a *SysDictList) Del() error {
	if len(a.ID) != 32 {
		return errors.New("缺乏参数ID")
	}
	db.Db.Model(&SysDictList{}).Where("parent_id = ?", a.ID).Delete(&SysDictList{})
	err := db.Db.Model(a).Delete(a).Error
	return err
}

func GetDictByParent(parent string) []SysDictList {
	var items []SysDictList
	db.Db.Model(&SysDictList{}).Where("parent_id = ?", parent).Find(&items)
	return items
}

type dictlist struct {
	Id    string `json:"id"`
	Name  string `json:"label"`
	Value int    `json:"value"`
}

func (a *SysDictList) GetByCode() []dictlist {
	var fdb SysDictList
	db.Db.Model(&fdb).Where("parent_name = ?", a.ParentName).First(&fdb)
	var items []dictlist
	db.Db.Table("sys_dict_list").Select("id,name ,value").Where("parent_id = ?", fdb.ID).Find(&items)
	return items
}
