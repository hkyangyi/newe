package route

import (
	v1 "neweadmin/app/common/system/v1"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.RouterGroup) {
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
	//组织结构
	r.POST("deptadd", v1.DepartAdd)
	r.PUT("deptedit", v1.DepartEdit)
	r.DELETE("deptdel", v1.DepartDel)
	r.GET("deptgetlist", v1.DepartGetList)
	r.GET("deptgetrules", v1.DepartRulesGet)
	r.POST("deptrulessave", v1.DepartRulesSave)
	//用户管理
	r.POST("memberadd", v1.MemberAdd)
	r.PUT("memberedit", v1.MemberEdit)
	r.GET("membergetlist", v1.MemberGetList)
	r.DELETE("memberdel", v1.MemberDel)
}
