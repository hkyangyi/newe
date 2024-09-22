package router

import (
	"github.com/hkyangyi/newe/router/middle"
	v1 "github.com/hkyangyi/newe/router/v1"

	"github.com/gin-gonic/gin"
)

//	@Summary	系统设置
//	@Schemes	base
//	@Description	系统管理
//	@Tags			v1
//	@Accept			json
//	@Produce		json
//	@Router			/api/v1/

func RegRouter(r *gin.RouterGroup) {
	// api := r.Group("api")
	admin := r.Group("admin")

	admin.POST("login", v1.Login)
	admin.GET("logout", v1.LoginOut)
	admin.POST("createMenu", v1.CreateMenu)
	admin.GET("dict", v1.GetDictByCode)
	admin.GET("GetRoleByDepart", v1.GetRoleByDepart)
	//中间件验证登录
	admin.Use(middle.AdminAuth())
	admin.GET("verifysole", v1.Verifysole)
	admin.GET("GetMenuTree", v1.BaseGetMenuTree)
	admin.GET("getUserInfo", v1.GetUserInfo)
	admin.GET("getPermCode", v1.GetPermCode)
	admin.GET("getMenuList", v1.GetMenuList)
	admin.PUT("upUserInfo", v1.UpUserInfo)
	admin.POST("EditPass", v1.EditPass)
	admin.POST("upimage", v1.UpImage)
	admin.POST("upfile", v1.FileUpload)

	sys := admin.Group("system")
	SystemRouter(sys)
}

func SystemRouter(r *gin.RouterGroup) {
	//菜单管理
	r.GET("getMenuList", v1.GetSysMenuList)
	r.POST("addMenuList", v1.SysMenuAdd)
	r.PUT("editMenuList", v1.SysMenuEdit)
	r.DELETE("delMenuList", v1.SysMenuDel)
	r.GET("MerGetMenu", v1.MerGetMenu)
	//字典管理
	r.POST("dictadd", v1.DictCreate)
	r.PUT("dictedit", v1.DictUpdate)
	r.GET("dictgetlist", v1.DictGetList)
	r.DELETE("dictdel", v1.DictDel)
	r.GET("dictgetbycode", v1.GetDictByCode)
	//组织结构
	r.POST("deptadd", v1.DepartAdd)
	r.PUT("deptedit", v1.DepartEdit)
	r.DELETE("deptdel", v1.DepartDel)
	r.GET("deptgetlist", v1.DepartGetList)
	r.GET("deptgetrules", v1.DepartRulesGet)
	r.POST("deptrulessave", v1.DepartRulesSave)
	//角色管理
	r.POST("roleadd", v1.RoleAdd)
	r.PUT("roleedit", v1.RoleEdit)
	r.GET("rolegetlist", v1.RoleGetList)
	r.DELETE("roledel", v1.RoleDel)
	r.POST("setRoleStatus", v1.RoleStatus)
	r.GET("roleRulesGetList", v1.RoleRulesGetList)
	r.PUT("roleRulesEdit", v1.RoleRulesEdit)
	//用户管理
	r.POST("memberadd", v1.MemberAdd)
	r.PUT("memberedit", v1.MemberEdit)
	r.GET("membergetlist", v1.MemberGetList)
	r.DELETE("memberdel", v1.MemberDel)
	//api管理
	r.POST("ApiListAdd", v1.ApiListAdd)
	r.PUT("ApiListEdit", v1.ApiListEdit)
	r.GET("ApiListGetList", v1.ApiListGetList)
	r.DELETE("ApiListDel", v1.ApiListDel)
	r.PUT("ApiListSetStatus", v1.ApiListSetStatus)
}
