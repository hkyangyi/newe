package cus

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/cusmodel"
	"github.com/hkyangyi/newe/router/app"
)

func RegisterCusDepartRoutes(g *gin.RouterGroup) {
	g.POST("add", DepartAdd)
	g.PUT("edit", DepartEdit)
	g.DELETE("del", DepartDel)
	g.GET("list", DepartGetList)
	g.GET("getrules", DepartRulesGet)
	g.POST("saverules", DepartRulesSave)
	g.GET("tree", departTree)
	g.GET("roles", deptGetRoles)
}

func DepartAdd(c *gin.Context) {
	var a = app.NewApp(c)
	var data cusmodel.CustomerDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	merdata, b := a.C.Get("CusAdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	data.Cid = mer.CDB.Id

	err := data.Add()
	if err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
	return
}

func DepartEdit(c *gin.Context) {
	var a = app.NewApp(c)
	var data cusmodel.CustomerDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}

	merdata, b := a.C.Get("CusAdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	data.Cid = mer.CDB.Id

	err := data.Edit()
	if err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
	return
}

func DepartDel(c *gin.Context) {
	var a = app.NewApp(c)
	var data cusmodel.CustomerDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	merdata, b := a.C.Get("CusAdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	data.Cid = mer.CDB.Id

	err := data.Del()
	if err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
	return
}

func DepartGetList(c *gin.Context) {
	var a = app.NewApp(c)
	var data cusmodel.CustomerDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	merdata, b := a.C.Get("CusAdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	data.Cid = mer.CDB.Id

	if mer.Isadmin {
		items := data.GetListByAdmin()
		a.SUCCESS(items)
		return
	} else {
		items := data.GetListByMember(mer.MerDb)
		a.SUCCESS(items)
		return
	}
}

func DepartRulesGet(c *gin.Context) {
	var a = app.NewApp(c)
	var data cusmodel.CustomerDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}

	items := data.GetRules()
	var ids []string
	for _, v := range items {
		ids = append(ids, v.MenuId)
	}

	a.SUCCESS(ids)
}

type DepartRulesForm struct {
	Id    string `json:"id" form:"id"`
	Menus string `json:"menus" form:"menus"`
}

func DepartRulesSave(c *gin.Context) {
	var a = app.NewApp(c)
	var data DepartRulesForm
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	var dp = cusmodel.CustomerDepart{
		ID: data.Id,
	}

	items := dp.GetRules()
	var ids []string
	for _, v := range items {
		ids = append(ids, v.MenuId)
	}
	farr := strings.Split(data.Menus, ",")

	var delarr, addarr []string
	for i := 0; i < len(ids); i++ {
		b := utils.InArray(ids[i], farr)
		if !b {
			delarr = append(delarr, ids[i])
		}
	}

	for _, v := range farr {
		if b := utils.InArray(v, ids); !b {
			addarr = append(addarr, v)
		}
	}

	if len(addarr) > 0 {
		err := dp.AddRules(addarr)
		if err != nil {
			a.Error(err)
			return
		}
	}

	if len(delarr) > 0 {
		err := dp.DelRules(delarr)
		if err != nil {
			a.Error(err)
			return
		}
	}

	a.SUCCESS(nil)
	return
}

// 获取不含角色的树形结构
func departTree(c *gin.Context) {
	var a = app.NewApp(c)
	var data cusmodel.CustomerDepart
	merdata, b := a.C.Get("CusAdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	items := data.GetDepartTree(mer.MerDb)
	a.SUCCESS(items)
	return
}

// 根据部门ID获取角色
func deptGetRoles(c *gin.Context) {
	var a = app.NewApp(c)
	var data cusmodel.CustomerDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	items := data.GetRoleById()
	a.SUCCESS(items)
	return
}
