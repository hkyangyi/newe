package v1

import (
	"newe/app/common/app"
	"newe/app/common/system/moddle"
	"newe/common/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func DictGetList(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysDictList
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}

	var wheremap []string
	var params []interface{}

	if len(data.ParentId) > 0 {
		wheremap = append(wheremap, "parent_id = ?")
		params = append(params, data.ParentId)
	}

	where := strings.Join(wheremap, " AND ")
	page := utils.PageList{
		Page:     data.Page,
		PageSize: data.PageSize,
	}
	items := data.GetList(page, where, params...)
	a.SUCCESS(items)
}

func DictCreate(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysDictList
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	if err := data.Add(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

func DictUpdate(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysDictList
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	if err := data.Edit(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

func DictDel(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysDictList
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	if err := data.Del(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

// 获取字典列表
func GetDictByCode(c *gin.Context) {
	var a = app.NewApp(c)
	var data moddle.SysDictList
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	items := data.GetByCode()
	a.SUCCESS(items)
}
