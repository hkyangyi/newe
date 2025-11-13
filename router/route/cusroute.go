package router

import (
	"github.com/hkyangyi/newe/router/base"
	"github.com/hkyangyi/newe/router/cus"
	"github.com/hkyangyi/newe/router/middle"
	v2 "github.com/hkyangyi/newe/router/v2"

	"github.com/gin-gonic/gin"
)

func NewesysCusRouter(r *gin.RouterGroup) {

	v2.RegisterWsRoutes(r.Group("ws"))
	// api := r.Group("api")
	auth := r.Group("auth")

	auth.POST("login", cus.Login)
	auth.Use(middle.CusAdminAuth())
	auth.POST("logout", cus.LoginOut)
	auth.POST("refresh", cus.RefreshToken)
	auth.GET("codes", cus.GetPermCode)

	user := r.Group("user", middle.CusAdminAuth())
	user.GET("info", cus.GetUserInfo)
	// user.PUT("upinfo", v2.UpUserInfo)
	// user.POST("editpass", v2.EditPass)
	// user.POST("upimage", v2.UpImage)
	// user.POST("upfile", v2.FileUpload)

	menu := r.Group("menu", middle.CusAdminAuth())
	menu.GET("all", cus.GetMenuAll)
	//用户中心
	cus.RegisterUserRoutes(r.Group("sysuser", middle.CusAdminAuth()))
	cus.RegisterCusDepartRoutes(r.Group("depart", middle.CusAdminAuth()))
	cus.RegisterCusMemberRoutes(r.Group("member", middle.CusAdminAuth()))
	//注册菜单路由
	// cus.RegisterMenuRoutes(menu)
	// //注册字典路由
	// v2.RegisterSysDictRoutes(r.Group("dict", middle.CusAdminAuth()))
	// //注册组织结构路由
	// v2.RegisterSysDepartRoutes(r.Group("depart", middle.CusAdminAuth()))
	// //注册角色路由
	// v2.RegisterSysMemberRoutes(r.Group("member", middle.CusAdminAuth()))
	// //注册操作日志路由
	cus.RegisterOplogRoutes(r.Group("oplog", middle.CusAdminAuth()))
	// //v2.RegisterSysApiListRoutes(r.Group("sysapi", middle.AdminAuth()))
	// v2.RegisterCustomerMenusRoutes(r.Group("cusmenu", middle.CusAdminAuth()))

	bs := r.Group("base", middle.CusAdminAuth())
	//检测唯一性
	bs.GET("verifysole", v2.Verifysole)
	//上传
	bs.GET("dictcode", base.GetDictByCode)
	bs.GET("menusTree", cus.MenuTree)
	bs.POST("upload", cus.UploadFile)
	cus.RegisterUploadRoutes(bs)

}
