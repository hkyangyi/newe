package cus

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/cusmodel"
	"github.com/hkyangyi/newe/router/app"
)

func RegisterCusMemberRoutes(g *gin.RouterGroup) {
	g.POST("add", MemberAdd)
	g.PUT("edit", MemberEdit)
	g.DELETE("del", MemberDel)
	g.GET("page", MemberGetPage)
}

func MemberAdd(c *gin.Context) {
	var (
		g = app.NewApp(c)
		a cusmodel.CustomerMember
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	merdata, b := g.C.Get("CusAdminAuthData")
	if !b {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	a.Password = utils.EncodeMD5(a.Password)
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = time.Now().Unix()
	a.Cid = mer.CDB.Id
	fdp, err := cusmodel.CustomerDepartGetBYId(a.RoleId)
	if err != nil {
		g.Error(errors.New("部门参数错误"))
		return
	}
	if len(fdp.ID) == 0 {
		g.Error(errors.New("部门参数错误"))
		return
	}
	a.OrgCode = fdp.Code

	inn := strings.Index(fdp.Code, mer.MerDb.OrgCode)
	if mer.MerDb.Username != "admin" && inn < 0 {
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
		a cusmodel.CustomerMember
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	merdata, b := g.C.Get("CusAdminAuthData")
	if !b {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)

	if len(a.ID) != 32 {
		g.LoginError(errors.New("参数ID错误"))
		return
	}

	if len(a.Password) > 0 {
		a.Password = utils.EncodeMD5(a.Password)
	}
	fdp, err := cusmodel.CustomerDepartGetBYId(a.DepartId)
	if err != nil {
		g.Error(errors.New("部门参数错误"))
		return
	}
	if len(fdp.ID) == 0 {
		g.Error(errors.New("部门参数错误"))
		return
	}
	a.Cid = mer.CDB.Id
	a.OrgCode = fdp.Code

	inn := strings.Index(fdp.Code, mer.MerDb.OrgCode)
	if mer.MerDb.Username != "admin" && inn < 0 {
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
		a cusmodel.CustomerMember
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	merdata, b := g.C.Get("CusAdminAuthData")
	if !b {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)

	a.Cid = mer.CDB.Id

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

func MemberGetPage(c *gin.Context) {

	var (
		g = app.NewApp(c)
		a cusmodel.CustomerMember
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	merdata, b := g.C.Get("CusAdminAuthData")
	if !b {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	a.Cid = mer.CDB.Id

	//测试消息
	//message.AddMessage(mer.MerDb.ID, "系统消息", "欢迎使用newE系统！", "info", "system	", "notify")

	var wheremap []string
	var params []interface{}

	if mer.MerDb.Username != "admin" {
		wheremap = append(wheremap, " org_code like ?")
		params = append(params, mer.MerDb.OrgCode+"%")
	}

	if len(a.Username) > 0 {
		wheremap = append(wheremap, " username = ?")
		params = append(params, a.Username)
	}

	if len(a.DepartId) > 0 {
		wheremap = append(wheremap, " depart_id = ?")
		params = append(params, a.DepartId)
	}

	if len(a.RealName) > 0 {
		wheremap = append(wheremap, " real_name like ? ")
		params = append(params, a.RealName+"%")
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
