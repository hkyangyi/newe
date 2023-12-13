package v1

import (
	"encoding/json"
	"errors"
	"newe/app/common/app"
	"newe/app/common/system/moddle"
	"newe/common/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetMenuList(c *gin.Context) {
	var a = app.NewApp(c)
	merdata, b := a.C.Get("AdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(moddle.SysMember)
	var items []moddle.SysMenus
	if mer.Username == "admin" {
		items = moddle.GetMenusList()
	} else {
		items = moddle.SysMenusListGetByDepart(mer.DepartId)
	}

	data := GetMenuListByAdmin(items)
	a.SUCCESS(data)
}

func GetMenuListByAdmin(menus []moddle.SysMenus) []VueRouteMenu {
	var items []VueRouteMenu
	for _, v := range menus {
		if v.Type == 3 {
			break
		}
		if v.Status == 0 {
			break
		}
		var item VueRouteMenu
		item.Id = v.ID
		item.Pid = v.Pid
		item.Path = v.RoutePath
		item.Name = v.RouteName
		item.Component = v.Component
		item.Meta.Title = v.Name
		item.Meta.Icon = v.Icon
		item.Redirect = v.Redirect
		if v.Keepalive > 1 {
			item.Meta.IgnoreKeepAlive = false
		} else {
			item.Meta.IgnoreKeepAlive = true
		}
		if v.Show > 0 {
			item.Meta.HideMenu = false
		} else {
			item.Meta.HideMenu = true
		}
		item.Meta.OrderNo = v.SortNo
		item.Meta.FrameSrc = v.IframeSrc

		if len(v.Props) > 4 {
			var props = make(map[string]interface{})
			json.Unmarshal([]byte(v.Props), &props)
			item.Props = props
		}

		if v.IsExt > 0 {
			item.Meta.IsLink = true
		} else {
			item.Meta.IsLink = false
		}
		items = append(items, item)
	}
	var bitems []VueRouteMenu
	bitems = menuDigui(items, "", bitems)
	return bitems
}

func menuDigui(items []VueRouteMenu, pid string, list []VueRouteMenu) []VueRouteMenu {
	var item []VueRouteMenu
	for _, v := range items {
		if v.Pid == pid {
			var ls VueRouteMenu
			utils.StAtoB(v, ls, &ls)
			ls.Children = menuDigui(items, v.Id, item)
			list = append(list, ls)
		}
	}
	return list
}

// @Tags 菜单管理
// @Produce  json
// @Param name query string true "菜单名称" maxlength(100)
// @Success 200 {object}  "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/system/getMenuList [get]
func GetSysMenuList(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysMenus
	a.Bind(&data)
	var params []interface{}
	var where []string

	if len(data.Name) > 2 {
		w := "name like ?"
		params = append(params, data.Name)
		where = append(where, w)
	}

	ws := strings.Join(where, " AND ")

	items := data.GetList(ws, params...)
	a.SUCCESS(items)
	return
}

// 添加菜单
func SysMenuAdd(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysMenus
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	if err := data.Add(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

// 编辑菜单
func SysMenuEdit(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysMenus
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	if err := data.Edit(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

// 删除菜单
func SysMenuDel(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysMenus
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	if err := data.Del(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

// 获取菜单列表
func MerGetMenu(c *gin.Context) {
	var a = app.NewApp(c)
	merdata, b := a.C.Get("AdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(moddle.SysMember)
	//a.SUCCESS(mer)

	var params []interface{}
	var where []string

	if mer.Username != "admin" {
		fd, err := moddle.SysDepartGetByCode(mer.OrgCode)
		if err != nil {
			a.Error(err)
			return
		}
		items := moddle.DepartRulesByDeaprtId(fd.ID)
		var meuls []string
		for _, v := range items {
			meuls = append(meuls, v.MenuId)
		}

		w := "id in ?"
		params = append(params, meuls)
		where = append(where, w)

	}

	w := "status = ?"
	params = append(params, 1)
	where = append(where, w)

	var md moddle.SysMenus

	ws := strings.Join(where, " AND ")

	items := md.GetList(ws, params...)
	a.SUCCESS(items)
	return
}
