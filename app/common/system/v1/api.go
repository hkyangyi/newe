package v1

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/app/common/app"
	"github.com/hkyangyi/newe/app/common/system/moddle"
	"github.com/hkyangyi/newe/common/utils"
)

func ApiListAdd(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a moddle.SysApiList
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

func ApiListEdit(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a moddle.SysApiList
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

func ApiListDel(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a moddle.SysApiList
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

func ApiListGetList(c *gin.Context) {

	var (
		g = app.NewApp(c)
		a moddle.SysApiList
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
	// mer := merdata.(moddle.SysMember)

	var wheremap []string
	var params []interface{}

	// if mer.Username != "admin" {
	// 	wheremap = append(wheremap, " org_code like ?")
	// 	params = append(params, mer.OrgCode+"%")
	// }

	if len(a.MenuId) > 0 {
		wheremap = append(wheremap, " menu_id = ?")
		params = append(params, a.MenuId)
	}

	where := strings.Join(wheremap, " AND ")
	page := utils.PageList{
		Page:     a.Page,
		PageSize: a.PageSize,
	}
	items := a.GetList(page, where, params...)
	g.SUCCESS(items)

}

func ApiListSetStatus(c *gin.Context) {

	var (
		g = app.NewApp(c)
		a moddle.SysApiList
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	if err := a.SetStatus(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)

}
