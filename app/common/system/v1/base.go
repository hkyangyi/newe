package v1

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/app/common/app"
	"github.com/hkyangyi/newe/app/common/system/moddle"
	"github.com/hkyangyi/newe/common/db"
)

// 获取菜单列表
func BaseGetMenuTree(c *gin.Context) {
	var a = app.NewApp(c)
	merdata, b := a.C.Get("AdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(moddle.AdminAuth)
	//a.SUCCESS(mer)

	var params []interface{}
	var where []string

	if mer.Merdb.Username != "admin" {
		fd, err := moddle.SysDepartGetByCode(mer.Merdb.OrgCode)
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

// @验证字段唯一

// @Tags Base
// @Description 验证字段是否存在
// @Accept  json
// @Produce json
// @Param tablename  path   string    true   "表名字"
// @Param fieldname  path   string    true   "字段名字"
// @Param tablevalue  path   string    true   "字段值"
// @Param tableid  path   int    true  "过滤ID 在编辑时传入ID"
// @Success 200 {string} string	"ok"
// @Router /api/base/verifysole [get]
func Verifysole(c *gin.Context) {
	var a = app.NewApp(c)
	var data verifyform
	if e := a.Bind(&data); e != nil {
		a.Error(errors.New("参数错误"))
		return
	}

	where := make(map[string]interface{})
	where[data.FieldName] = data.Tablevalue

	res := db.VerifyOnly(data.TableName, data.TableId, where)
	if res {
		a.Error(errors.New("已存在"))
		return
	}
	a.SUCCESS(nil)
}
