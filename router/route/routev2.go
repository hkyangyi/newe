package router

import (
	"github.com/hkyangyi/newe/router/middle"
	v2 "github.com/hkyangyi/newe/router/v2"

	"github.com/gin-gonic/gin"
)

func NewesysRouter(r *gin.RouterGroup) {

	v2.RegisterWsRoutes(r.Group("ws"))
	// api := r.Group("api")
	auth := r.Group("auth")

	auth.POST("login", v2.Login)
	auth.Use(middle.AdminAuth())
	auth.POST("logout", v2.LoginOut)
	auth.POST("refresh", v2.RefreshToken)
	auth.GET("codes", v2.GetPermCode)

	user := r.Group("user", middle.AdminAuth())
	user.GET("info", v2.GetUserInfo)
	// user.PUT("upinfo", v2.UpUserInfo)
	// user.POST("editpass", v2.EditPass)
	// user.POST("upimage", v2.UpImage)
	// user.POST("upfile", v2.FileUpload)

	menu := r.Group("menu", middle.AdminAuth())
	menu.GET("all", v2.GetMenuAll)
	//用户中心
	v2.RegisterUserRoutes(r.Group("sysuser", middle.AdminAuth()))
	//注册菜单路由
	v2.RegisterSysMenuRoutes(menu)
	//注册字典路由
	v2.RegisterSysDictRoutes(r.Group("dict", middle.AdminAuth()))
	//注册组织结构路由
	v2.RegisterSysDepartRoutes(r.Group("depart", middle.AdminAuth()))
	//注册角色路由
	v2.RegisterSysMemberRoutes(r.Group("member", middle.AdminAuth()))
	//注册操作日志路由
	v2.RegisterOplogRoutes(r.Group("oplog", middle.AdminAuth()))
	//v2.RegisterSysApiListRoutes(r.Group("sysapi", middle.AdminAuth()))

	base := r.Group("base", middle.AdminAuth())
	//检测唯一性
	base.GET("verifysole", v2.Verifysole)
	//上传
	base.POST("upload", v2.UploadFile)

}
