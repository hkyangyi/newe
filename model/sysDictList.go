package model

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
	Status     int    `json:"status" form:"status"`             //1正常，2停用，10已删除
	Sort       int    `json:"sort"`                             //排序
	Type       int    `json:"type"`                             //1自定义 2 数据表
	TableKey   string `json:"tableKey"`                         //数据表主键
	TableVal   string `json:"tableVal"`                         //数据包值
	Issys      int    `json:"issys"`                            //是否系统内置，1是，0否
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

func (a *SysDictList) GetList() []SysDictList {
	var items []SysDictList
	db.Db.Model(&SysDictList{}).Where("parent_id = ?", "").Order("sort asc").Find(&items)
	return items
}

// 获取列表
func (a *SysDictList) GetPage(page utils.PageList, where string, v ...interface{}) utils.PageList {
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
	db.Db.Model(&SysDictList{}).Where("parent_id = ? and  issys <> 1", a.ID).Delete(&SysDictList{})
	err := db.Db.Model(a).Delete(a).Error
	return err
}

func GetDictByParent(parent string) []SysDictList {

	var fdb SysDictList
	db.Db.Model(&fdb).Where("id = ?", parent).First(&fdb)
	db.UpdateDict(fdb.Value) //更新字典缓存
	var items []SysDictList
	if fdb.Type == 1 {
		db.Db.Table("sys_dict_list").Where("parent_id = ? ", fdb.ID).Order("sort asc").Find(&items)
		return items
	} else {
		table := utils.Camel2Case(fdb.Value)
		selstr := fmt.Sprintf("%s as value, %s as name", fdb.TableKey, fdb.TableVal)
		db.Db.Table(table).Select(selstr).Scan(&items)
		return items
	}

}

// 更新状态
func (a *SysDictList) UpdateStatus() error {
	if len(a.ID) != 32 {
		return errors.New("缺乏参数ID")
	}
	err := db.Db.Model(a).Update("status", a.Status).Error
	return err
}

type dictlist struct {
	Id     string `json:"id"`
	Name   string `json:"name" gorm:"column:name"`
	Value  string `json:"value" gorm:"column:value"`
	Status int    `json:"status"` //1正常，2停用，10已删除
	Sort   int    `json:"sort"`   //排序
}

func (a *SysDictList) GetByCode(code string) []dictlist {
	var fdb SysDictList
	db.Db.Model(&fdb).Where("value = ?", code).First(&fdb)
	var items []dictlist
	if fdb.Type == 1 {
		db.Db.Table("sys_dict_list").Where("parent_id = ? and status = ?", fdb.ID, 1).Order("sort asc").Find(&items)
		return items
	} else {

		table := utils.Camel2Case(fdb.Value)
		selstr := fmt.Sprintf("%s as value, %s as name", fdb.TableKey, fdb.TableVal)
		db.Db.Table(table).Select(selstr).Scan(&items)
		return items
	}
}
