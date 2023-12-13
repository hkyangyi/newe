package v1

import (
	"errors"
	"neweadmin/app/common/app"
	"neweadmin/app/common/system/moddle"
	"neweadmin/common/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func DepartAdd(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}

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
	var data moddle.SysDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}

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
	var data moddle.SysDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}

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
	var data moddle.SysDepart
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	merdata, b := a.C.Get("AdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(moddle.SysMember)
	items := data.GetList(mer)
	a.SUCCESS(items)
	return
}

func DepartRulesGet(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysDepart
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
	var dp = moddle.SysDepart{
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
		b := utils.InArrar(ids[i], farr)
		if !b {
			delarr = append(delarr, ids[i])
		}
	}

	for _, v := range farr {
		if b := utils.InArrar(v, ids); !b {
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
