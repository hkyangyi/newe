package moddle

import (
	"errors"
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

type SysRole struct {
	ID         string `gorm:"primary_key" json:"id" form:"id"`                       //角色ID
	DepartId   string `json:"departId" form:"departId" dict:"DepartName_sysDepart" ` //组织结构ID
	DepartName string `gorm:"-" json:"departName"`
	OrgCode    string `json:"orgCode"`    //组织结构代码
	RoleName   string `json:"roleName"`   //角色名称
	Status     int    `json:"status"`     //角色状态(-1禁用，1启用)
	CreateTime int64  `json:"createTime"` //创建时间
	UpdateTime int64  `json:"updateTime"` //修改时间
	Remark     string `json:"remark"`     //备注
	utils.PageList
}

func GetSysRoleByDepart(id string) []SysRole {
	var items []SysRole
	db.Db.Select("id", "role_name").Where("depart_id = ?", id).Find(&items)
	return items
}

// 添加
func (a *SysRole) Add() error {
	a.ID = utils.GetUUID()
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = time.Now().Unix()
	if sysdepart, err := SysDepartGetBYId(a.DepartId); err == nil {
		a.OrgCode = sysdepart.Code
	}

	err := db.Db.Create(a).Error
	return err
}

// 编辑
func (a *SysRole) Edit() error {
	if sysdepart, err := SysDepartGetBYId(a.DepartId); err != nil {
		a.OrgCode = sysdepart.Code
	}
	a.UpdateTime = time.Now().Unix()
	err := db.Db.Model(a).Updates(a).Error
	return err
}

// 启停
func (a *SysRole) SetRoleStatus() error {
	if a.ID == "" {
		return errors.New("缺少参数ID")
	}
	a.UpdateTime = time.Now().Unix()
	err := db.Db.Model(a).Updates(a).Error
	return err
}

// 获取列表
func (a *SysRole) GetList(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []SysRole
	db.Db.Model(&SysRole{}).Where(where, v...).Count(&page.Total).Order("create_time desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	page.List = items
	return page
}

// 删除
func (a *SysRole) Del() error {
	err := db.Db.Model(a).Delete(a).Error
	return err
}

// 根据角色获取菜单列表
func (a *SysRole) GetRoleMenus() []SysMenus {
	//获取该角色可用菜单
	items := SysMenusListGetByDepart(a.DepartId)
	var list []SysMenus
	list = SysMenusDigui(items, "", list)
	//根据菜单获取API列表
	return list
}

type SysRoleRules struct {
	ID         string `gorm:"primary_key" json:"id" form:"id"` //
	RoleId     string `json:"roleId"`                          //
	ApiId      string `json:"apiId"`
	MenuId     string `json:"menuId"`     //
	OrgCode    string `json:"orgCode"`    //部门数据
	WhereCode  string `json:"whereCode"`  //其他条件
	Status     int    `json:"status"`     //权限状态(-1禁用，1启用)
	CreateTime int64  `json:"createTime"` //
	UpdateTime int64  `json:"updateTime"` //修改时间
}

type SysRoleRulesList struct {
	ID         string `gorm:"primary_key" json:"id" form:"id"`
	MenuId     string `json:"menuId" form:"menuId"  dict:"MenuName_sysMenus" `
	MenuName   string `gorm:"-" json:"menuName"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Name       string `json:"name"`
	RulesId    string `json:"rulesId" gorm:"column:rulesId"` //权限ID
	RoleId     string `json:"roleId" form:"roleId"`          //校色ID
	OrgCode    string `json:"orgCode"`                       //部门数据
	WhereCode  string `json:"whereCode"`                     //其他条件
	Status     int    `json:"status"`                        //权限状态(-1禁用，1启用)
	CreateTime int64  `json:"createTime"`                    //
	UpdateTime int64  `json:"updateTime"`                    //修改时间
}

// 根据角色获取api列表
func (a *SysRoleRulesList) GetRoleReles() []SysRoleRulesList {
	var items []SysRoleRulesList

	db.Db.Select("a.id,a.menu_id,a.path,a.method,a.name,b.id as rulesId,b.org_code,b.where_code,b.status,b.create_time,b.update_time").
		Table("sys_api_list as a").
		Where("a.menu_id = ?", a.MenuId).
		Joins("left join sys_role_rules as b on b.api_id = a.id and b.role_id = ?", a.RoleId).
		Find(&items)
	for i := 0; i < len(items); i++ {
		if items[i].RoleId == "" {
			items[i].RoleId = a.RoleId
		}
	}
	return items
}

func (a *SysRoleRulesList) Edit() error {
	var item = SysRoleRules{
		ID:         a.RulesId,
		ApiId:      a.ID,
		RoleId:     a.RoleId,
		MenuId:     a.MenuId,
		OrgCode:    a.OrgCode,
		WhereCode:  a.WhereCode,
		Status:     a.Status,
		CreateTime: a.CreateTime,
		UpdateTime: time.Now().Unix(),
	}
	if item.ID == "" {
		item.ID = utils.GetUUID()
	}
	if item.CreateTime == 0 {
		item.CreateTime = time.Now().Unix()
	}
	err := db.Db.Save(&item).Error
	return err
}

// //角色api启停
// func (a *SysRole)
