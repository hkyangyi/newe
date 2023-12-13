package moddle

import (
	"errors"
	"github.com/hkyangyi/newe/common/base"
	"github.com/hkyangyi/newe/common/utils"
	"strings"
	"time"
)

type SysMenus struct {
	ID         string     `gorm:"primary_key" json:"id" form:"name"` //
	Pid        string     `json:"pid"`                               //
	Name       string     `json:"name" form:"name"`                  //
	Component  string     `json:"component"`                         //组件地址
	Icon       string     `json:"icon"`                              //图片
	IsExt      int        `json:"isExt"`                             //是否外链
	Keepalive  int        `json:"keepalive"`                         //是否缓存
	Show       int        `json:"show"`                              //是否显示
	Type       int        `json:"type"`                              //类型 1目录2菜单 3按钮
	SortNo     int        `json:"sortNo"`                            //排序
	RoutePath  string     `json:"routePath"`                         //路由地址
	Permission string     `json:"permission"`                        //权限编码
	Status     int        `json:"status"`                            //是否启用 0启用1禁用
	CreateTime int64      `json:"createTime"`                        //
	IframeSrc  string     `json:"iframeSrc"`                         //是否嵌套
	RouteName  string     `json:"routeName"`                         //路由名称
	Redirect   string     `json:"redirect"`                          //默认页
	Props      string     `json:"props"`
	List       []SysMenus `gorm:"-" json:"children"`
}

func (a *SysMenus) Add() error {
	a.ID = utils.GetUUID()
	a.CreateTime = time.Now().Unix()
	err := base.MYDB.Create(a).Error
	return err
}

func (a *SysMenus) Edit() error {
	if len(a.ID) != 32 {
		return errors.New("缺少参数ID")
	}
	err := base.MYDB.Updates(a).Error
	return err
}

func (a *SysMenus) Del() error {
	if len(a.ID) != 32 {
		return errors.New("缺少参数ID")
	}
	err := base.MYDB.Model(a).Delete(a).Error
	return err
}

func GetMenusList() []SysMenus {
	var items []SysMenus
	base.MYDB.Where("type <> ?", 3).Order("sort_no asc").Find(&items)
	return items
}

// 根据部门ID获取菜单
func SysMenusListGetByDepart(departId string) []SysMenus {
	var items []SysDepartRules
	base.MYDB.Model(&items).Where("depart_id = ?", departId).Find(&items)
	var meuls []string
	for _, v := range items {
		meuls = append(meuls, v.MenuId)
	}

	var res []SysMenus
	base.MYDB.Model(&SysMenus{}).Where("type <> ?", 3).Where("id in ?", meuls).Where("status = ?", 1).Order("sort_no asc").Find(&res)
	return res
}

func SysButtonGetList() []string {
	var items []SysMenus
	base.MYDB.Where("type = ?", 3).Order("sort_no asc").Find(&items)

	var btns []string
	for _, v := range items {
		btns = append(btns, v.Permission)
	}

	return btns
}

// 获取按钮
func SysButtonGetByDepart(departId string) []string {
	var items []SysDepartRules
	base.MYDB.Model(&items).Where("depart_id = ?", departId).Find(&items)
	var meuls []string
	for _, v := range items {
		meuls = append(meuls, v.MenuId)
	}
	var res []SysMenus
	base.MYDB.Where("type = ?", 3).Where("id in ?", meuls).Where("status = ?", 1).Order("sort_no asc").Find(&res)

	var btns []string
	for _, v := range res {
		btns = append(btns, v.Permission)
	}

	return btns
}

func (a *SysMenus) GetList(where string, v ...interface{}) []SysMenus {
	var items []SysMenus
	base.MYDB.Table("sys_menus").Where(where, v...).Order("sort_no asc").Find(&items)
	t := strings.Index(where, "name")
	if t >= 0 {
		return items
	}

	var list []SysMenus
	list = SysMenusDigui(items, "", list)
	return list
}

func SysMenusDigui(items []SysMenus, pid string, list []SysMenus) []SysMenus {
	var item []SysMenus
	for _, v := range items {
		if v.Pid == pid {
			var ls SysMenus
			utils.StAtoB(v, ls, &ls)
			ls.List = SysMenusDigui(items, v.ID, item)
			list = append(list, ls)
		}

	}
	return list
}
