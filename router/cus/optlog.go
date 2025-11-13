package cus

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/cusmodel"
	"github.com/hkyangyi/newe/router/app"
)

func RegisterOplogRoutes(g *gin.RouterGroup) {
	g.GET("page", OplogGetPage)
}

func OplogGetPage(c *gin.Context) {
	var a = app.NewApp(c)
	var data cusmodel.CustomerOperationLog
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}

	merdata, b := a.C.Get("CusAdminAuthData")
	if !b {
		a.LoginError(nil)
		return
	}
	mer := merdata.(cusmodel.CusAuth)

	var wheres []string
	var params []interface{}

	wheres = append(wheres, "cid = ?")
	params = append(params, mer.CDB.Id)

	//用户查询根据用户名
	if data.Username != "" {
		wheres = append(wheres, "username LIKE ?")
		params = append(params, "%"+data.Username+"%")
	}

	if data.Method != "" {
		wheres = append(wheres, "method = ?")
		params = append(params, data.Method)
	}

	if data.ApiID > 0 {
		wheres = append(wheres, "api_id = ?")
		params = append(params, data.ApiID)
	}
	if data.StatusCode != 0 {
		wheres = append(wheres, "status_code = ?")
		params = append(params, data.StatusCode)
	}

	if data.Result != "" {
		wheres = append(wheres, "result = ?")
		params = append(params, data.Result)
	}

	if !data.StartTime.IsZero() {
		wheres = append(wheres, "created_at >= ?")
		params = append(params, data.StartTime)
	}
	if !data.EndTime.IsZero() {
		wheres = append(wheres, "created_at <= ?")
		params = append(params, data.EndTime)
	}

	where := strings.Join(wheres, " AND ")
	page := utils.PageList{Page: data.Page, PageSize: data.PageSize}
	items := data.GetPage(page, where, params...)
	a.SUCCESS(items)
}
