package moddle

import (
	"errors"
	"fmt"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

type SysDictList struct {
	ID         string `gorm:"primary_key" json:"id"  form:"id"` //
	ParentId   string `json:"parentId" form:"parentId"`         //
	ParentName string `json:"parentName" form:"parentName"`     //字典名称
	Value      string `json:"value"`                            //值
	Name       string `json:"name"`                             //名称
	Status     int    `json:"status"`                           //1正常，2停用，10已删除
	Sort       int    `json:"sort"`                             //排序
	Type       int    `json:"type"`                             //1自定义 2 数据表
	TableKey   string `json:"tableKey"`                         //数据表主键
	TableVal   string `json:"tableVal"`                         //数据包值
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
	Id     string `json:"id"`
	Name   string `json:"name" gorm:"column:name"`
	Value  string `json:"value" gorm:"column:value"`
	Status int    `json:"status"` //1正常，2停用，10已删除
	Sort   int    `json:"sort"`   //排序
}

func (a *SysDictList) GetByCode() []dictlist {
	var fdb SysDictList
	db.Db.Model(&fdb).Where("parent_name = ?", a.ParentName).First(&fdb)
	var items []dictlist
	if fdb.Type == 1 {
		db.Db.Table("sys_dict_list").Where("parent_id = ?", fdb.ID).Order("sort asc").Find(&items)
		return items
	} else {

		table := utils.Camel2Case(fdb.ParentName)
		selstr := fmt.Sprintf("%s as value, %s as name", fdb.TableKey, fdb.TableVal)
		db.Db.Table(table).Select(selstr).Scan(&items)
		return items
	}
}
