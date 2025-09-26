package v2

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/model"
	"github.com/hkyangyi/newe/router/app"
)

// RegisterSysDictRoutes 注册字典相关路由
// 建议在外部 route 初始化时：
//
//	api := r.Group("/newesys/dict")
//	v2.RegisterSysDictRoutes(api)
func RegisterSysDictRoutes(g *gin.RouterGroup) {
	g.POST("add", DictCreate)
	g.PUT("edit", DictUpdate)
	g.DELETE("del", DictDel)
	g.GET("page", DictGetPage)
	g.GET("parent", DictGetByParent)
	g.GET("code", GetDictByCode)
	g.GET("list", DictGetList)
	g.GET("status", UpdateStatus)
}

func DictGetList(c *gin.Context) {
	var a = app.NewApp(c)
	var data model.SysDictList
	items := data.GetList()
	a.SUCCESS(items)
}

// DictGetList 字典项分页列表（可按 parentId / parentName 过滤）
func DictGetPage(c *gin.Context) {
	a := app.NewApp(c)
	var req model.SysDictList
	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}

	var wheres []string
	var params []interface{}

	wheres = append(wheres, "parent_id = ?")
	params = append(params, "")

	if req.ParentName != "" {
		wheres = append(wheres, "parent_name = ?")
		params = append(params, req.ParentName)
	}
	if req.Name != "" {
		wheres = append(wheres, "name LIKE ?")
		params = append(params, "%"+req.Name+"%")
	}

	where := strings.Join(wheres, " AND ")
	page := utils.PageList{Page: req.Page, PageSize: req.PageSize}
	items := req.GetPage(page, where, params...)
	a.SUCCESS(items)
}

func DictGetByParent(c *gin.Context) {
	a := app.NewApp(c)
	var req model.SysDictList
	if err := a.Bind(&req); err != nil {
		a.Error(err)
	}

	if req.ParentId == "" {
		a.Error(errors.New("缺少参数 parentId"))
		return
	}

	items := model.GetDictByParent(req.ParentId)
	a.SUCCESS(items)
}

// DictCreate 新建字典或字典项
func DictCreate(c *gin.Context) {
	a := app.NewApp(c)
	var req model.SysDictList
	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}
	if req.Status == 0 {
		req.Status = 1
	}
	if err := req.Add(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

// DictUpdate 更新字典项
func DictUpdate(c *gin.Context) {
	a := app.NewApp(c)
	var req model.SysDictList
	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}
	if err := req.Edit(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

// DictDel 逻辑删除（标记已删除）或物理删除子项
func DictDel(c *gin.Context) {
	a := app.NewApp(c)
	var req model.SysDictList
	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}
	if err := req.Del(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

// GetDictByCode 根据父级编码获取字典数据（支持自定义或数据表映射）
func GetDictByCode(c *gin.Context) {
	a := app.NewApp(c)
	var req model.SysDictList
	code := c.Query("code")
	if code == "" {
		a.Error(errors.New("缺少参数 code"))
		return
	}

	data := req.GetByCode(code)
	a.SUCCESS(data)
}

// 更新状态
func UpdateStatus(c *gin.Context) {
	a := app.NewApp(c)
	var req model.SysDictList
	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}
	if err := req.UpdateStatus(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
	return
}
