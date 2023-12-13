package router

import (
	systemrouter "neweadmin/app/common/system/route"
	v1 "neweadmin/app/common/system/v1"
	"neweadmin/app/middle"

	"github.com/gin-gonic/gin"
)

//	@Summary	系统设置
//	@Schemes	base
//	@Description	系统管理
//	@Tags			v1
//	@Accept			json
//	@Produce		json
//	@Router			/api/v1/

func RegRouter(r *gin.Engine) {
	api := r.Group("api")
	admin := api.Group("admin")

	admin.POST("login", v1.Login)
	admin.GET("loginOut", v1.LoginOut)
	admin.POST("createMenu", v1.CreateMenu)
	admin.GET("dict", v1.GetDictByCode)
	//中间件验证登录
	admin.Use(middle.AdminAuth())
	admin.POST("upload", v1.UpImage)
	admin.GET("verifysole", v1.Verifysole)
	admin.GET("getUserInfo", v1.GetUserInfo)
	admin.GET("getPermCode", v1.GetPermCode)
	admin.GET("getMenuList", v1.GetMenuList)
	admin.PUT("upUserInfo", v1.UpUserInfo)
	admin.POST("EditPass", v1.EditPass)

	sys := admin.Group("system")
	systemrouter.InitRouter(sys)
}
