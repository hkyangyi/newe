package v1

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/app/common/app"
	"github.com/hkyangyi/newe/app/common/system/moddle"
)

// 获取菜单列表
func BaseGetMenuTree(c *gin.Context) {
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
	params = append(params, 3)
	where = append(where, " type < ? ")

	w := "status = ?"
	params = append(params, 1)
	where = append(where, w)

	var md moddle.SysMenus

	ws := strings.Join(where, " AND ")

	items := md.GetList(ws, params...)
	a.SUCCESS(items)
	return
}

// 更具部门ID获取角色
func GetRoleByDepart(c *gin.Context) {
	var (
		g  = app.NewApp(c)
		id = c.Query("id")
	)

	if id == "" {
		g.Error(errors.New("缺少参数"))
		return
	}

	items := moddle.GetSysRoleByDepart(id)
	g.SUCCESS(items)

}
