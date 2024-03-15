package router

import (
	systemrouter "github.com/hkyangyi/newe/app/common/system/route"
	v1 "github.com/hkyangyi/newe/app/common/system/v1"
	"github.com/hkyangyi/newe/app/middle"

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
	admin.POST("upload", v1.UpImage)
	admin.GET("verifysole", v1.Verifysole)
	admin.GET("GetMenuTree", v1.BaseGetMenuTree)
	admin.GET("getUserInfo", v1.GetUserInfo)
	admin.GET("getPermCode", v1.GetPermCode)
	admin.GET("getMenuList", v1.GetMenuList)
	admin.PUT("upUserInfo", v1.UpUserInfo)
	admin.POST("EditPass", v1.EditPass)

	sys := admin.Group("system")
	systemrouter.InitRouter(sys)
}
