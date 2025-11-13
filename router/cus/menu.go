package cus

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/cusmodel"
	"github.com/hkyangyi/newe/router/app"
)

// 保留简化后的前端菜单结构

// SimpleMenuMeta 对齐用户期望返回结构的 meta 字段
type SimpleMenuMeta struct {
	Title      string   `json:"title"`
	Icon       string   `json:"icon,omitempty"`
	Order      int      `json:"order,omitempty"`
	Authority  []string `json:"authority,omitempty"`
	KeepAlive  bool     `json:"keepAlive,omitempty"`
	HideInMenu bool     `json:"hideInMenu,omitempty"`
	ActivePath string   `json:"activePath,omitempty"`
	AffixTab   bool     `json:"affixTab,omitempty"`
	IframeSrc  string   `json:"iframeSrc,omitempty"`
	Link       string   `json:"link,omitempty"`
	HideInTabe bool     `json:"hideInTab,omitempty"`
}

// SimpleMenu 对齐用户期望返回结构
type SimpleMenu struct {
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Redirect  string         `json:"redirect,omitempty"`
	Component string         `json:"component"`
	Meta      SimpleMenuMeta `json:"meta"`
	Children  []SimpleMenu   `json:"children,omitempty"`
}

func GetMenuAll(c *gin.Context) {
	var a = app.NewApp(c)
	authData, ok := a.C.Get("CusAdminAuthData")
	if !ok {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	admin := false
	// 简单判断是否 admin 用户
	mer, ok2 := authData.(cusmodel.CusAuth)
	if ok2 && mer.Isadmin {
		admin = true
	}
	var tree []cusmodel.CustomerMenus

	if admin {
		// admin 用户直接返回全部菜单
		tree = cusmodel.CustomerMenuTree()
	} else {
		tree = cusmodel.CustomerMenuTreeByDepart(mer.MerDb.RoleId)
	}

	// 取菜单树（全部正常记录）
	if c.Query("debug") == "1" {
		a.SUCCESS(tree)
		return
	}
	// 转换成前端需要结构
	res := convertMenus(tree, admin)
	a.SUCCESS(res)
}

func MenuTree(c *gin.Context) {
	var g = app.NewApp(c)
	merdata, ok := g.C.Get("CusAdminAuthData")
	if !ok {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	if err := mer.RefreshByMerdb(); err != nil {
		g.LoginError(err)
		return
	}
	if mer.Isadmin {
		tree := cusmodel.CustomerMenuTree()
		g.SUCCESS(tree)
	} else {
		tree := cusmodel.CustomerMenusListGetByDepart(mer.MerDb.RoleId)
		g.SUCCESS(tree)
	}
}

// convertMenus 递归转换 SysMenus -> SimpleMenu
func convertMenus(items []cusmodel.CustomerMenus, isAdmin bool) []SimpleMenu {
	result := make([]SimpleMenu, 0)
	for _, m := range items {
		// 过滤已删除或禁用
		if m.IsDeleted != 1 || m.Status != 1 {
			continue
		}
		// 按类型过滤：仅输出目录与菜单(type 1/2)，按钮(3)不在树结构中
		if m.Type == 3 {
			continue
		}
		sm := SimpleMenu{
			Name:      m.Name,
			Path:      m.FullPath,
			Redirect:  m.Redirect,
			Component: resolveComponent(m),
			Meta: SimpleMenuMeta{
				Title:      m.Title,
				Icon:       m.Icon,
				Order:      m.OrderNo,
				KeepAlive:  m.KeepAlive == 1,
				HideInMenu: m.HideInMenu == 1,
				ActivePath: m.ActivePath,
				AffixTab:   m.AffixTab == 1,
				IframeSrc:  m.IframeSrc,
				Link:       m.Link,
				HideInTabe: m.HideInTab == 1,
			},
		}
		if len(m.Children) > 0 {
			if m.HideChildrenInMenu == 1 {
				// 不传递 children
			} else {
				sm.Children = convertMenus(m.Children, isAdmin)
			}
		}
		result = append(result, sm)
	}
	return result
}

func authoritySlice(m cusmodel.CustomerMenus, isAdmin bool) []string {
	if isAdmin {
		return []string{"admin"}
	}
	if len(m.Authority) > 0 {
		return m.Authority
	}
	if m.Permission != "" {
		return []string{m.Permission}
	}
	return nil
}

func resolveComponent(m cusmodel.CustomerMenus) string {
	// 目录可能不需要实际组件，或由前端动态布局；如果有 iframe 则前端用特定组件
	if m.IframeSrc != "" {
		return "components/IFrame"
	}
	if m.Component != "" {
		return m.Component
	}
	// 默认：目录给一个布局标识，菜单没有则空由前端兜底
	if m.Type == 1 {
		return ""
	}
	return m.Component
}
