package v1

import (
	"errors"
	"newe/app/common/app"
	"newe/app/common/system/moddle"
	"newe/common/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func MemberAdd(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a moddle.SysMember
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	merdata, b := g.C.Get("AdminAuthData")
	if !b {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(moddle.SysMember)
	a.Password = utils.EncodeMD5(a.Password)
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = time.Now().Unix()
	fdp, err := moddle.SysDepartGetBYId(a.DepartId)
	if err != nil {
		g.Error(errors.New("部门参数错误"))
		return
	}
	if len(fdp.ID) == 0 {
		g.Error(errors.New("部门参数错误"))
		return
	}
	a.OrgCode = fdp.Code

	inn := strings.Index(fdp.Code, mer.OrgCode)
	if mer.Username != "admin" && inn < 0 {
		g.Error(errors.New("权限错误"))
		return
	}

	err = a.Add()
	if err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
	return
}

func MemberEdit(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a moddle.SysMember
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	merdata, b := g.C.Get("AdminAuthData")
	if !b {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(moddle.SysMember)

	if len(a.ID) != 32 {
		g.LoginError(errors.New("参数ID错误"))
		return
	}

	if len(a.Password) > 0 {
		a.Password = utils.EncodeMD5(a.Password)
	}
	fdp, err := moddle.SysDepartGetBYId(a.DepartId)
	if err != nil {
		g.Error(errors.New("部门参数错误"))
		return
	}
	if len(fdp.ID) == 0 {
		g.Error(errors.New("部门参数错误"))
		return
	}
	a.OrgCode = fdp.Code

	inn := strings.Index(fdp.Code, mer.OrgCode)
	if mer.Username != "admin" && inn < 0 {
		g.Error(errors.New("权限错误"))
		return
	}
	a.UpdateTime = time.Now().Unix()
	if err := a.Edit(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
	return
}

func MemberDel(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a moddle.SysMember
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
	return
}

func MemberGetList(c *gin.Context) {

	var (
		g = app.NewApp(c)
		a moddle.SysMember
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	merdata, b := g.C.Get("AdminAuthData")
	if !b {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(moddle.SysMember)

	var wheremap []string
	var params []interface{}

	if mer.Username != "admin" {
		wheremap = append(wheremap, " org_code like ?")
		params = append(params, mer.OrgCode+"%")
	}

	if len(a.Username) > 0 {
		wheremap = append(wheremap, " username = ?")
		params = append(params, a.Username)
	}

	if len(a.DepartId) > 0 {
		wheremap = append(wheremap, " depart_id = ?")
		params = append(params, a.DepartId)
	}

	if len(a.Realname) > 0 {
		wheremap = append(wheremap, " realname like ? ")
		params = append(params, a.Realname+"%")
	}

	where := strings.Join(wheremap, " AND ")
	page := utils.PageList{
		Page:     a.Page,
		PageSize: a.PageSize,
	}
	items := a.GetList(page, where, params...)
	g.SUCCESS(items)
	return
}
