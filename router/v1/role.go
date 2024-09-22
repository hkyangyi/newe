package v1

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/model"
	"github.com/hkyangyi/newe/router/app"
)

func RoleAdd(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a model.SysRole
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	if err := a.Add(); err != nil {
		g.Error(err)
		return
	}

	g.SUCCESS(nil)
}

func RoleEdit(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a model.SysRole
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	if err := a.Edit(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}

func RoleStatus(c *gin.Context) {

	var (
		g = app.NewApp(c)
		a model.SysRole
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	if err := a.SetRoleStatus(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}

func RoleDel(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a model.SysRole
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	if len(a.ID) != 32 {
		g.LoginError(errors.New("参数ID错误"))
		return
	}

	if err := a.Del(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}

func RoleGetList(c *gin.Context) {

	var (
		g = app.NewApp(c)
		a model.SysRole
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	// merdata, b := g.C.Get("AdminAuthData")
	// if !b {
	// 	g.LoginError(errors.New("登陆超时"))
	// 	return
	// }
	// mer := merdata.(model.SysMember)

	var wheremap []string
	var params []interface{}

	// if mer.Username != "admin" {
	// 	wheremap = append(wheremap, " org_code like ?")
	// 	params = append(params, mer.OrgCode+"%")
	// }

	if len(a.DepartId) > 0 {
		wheremap = append(wheremap, " depart_id = ?")
		params = append(params, a.DepartId)
	}

	where := strings.Join(wheremap, " AND ")
	page := utils.PageList{
		Page:     a.Page,
		PageSize: a.PageSize,
	}
	items := a.GetList(page, where, params...)
	g.SUCCESS(items)
}

// 根据角色和菜单获取api数据权限
func RoleRulesGetList(c *gin.Context) {

	var (
		g = app.NewApp(c)
		a model.SysRoleRulesList
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	if len(a.MenuId) != 32 || len(a.RoleId) != 32 {
		g.SUCCESS(nil)
		return
	}

	items := a.GetRoleReles()
	g.SUCCESS(items)
}

// 修改角色api数据权限
func RoleRulesEdit(c *gin.Context) {

	var (
		g = app.NewApp(c)
		a model.SysRoleRulesList
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	if len(a.MenuId) != 32 || len(a.RoleId) != 32 {
		g.Error(errors.New("缺少参数"))
		return
	}

	err := a.Edit()
	if err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}
