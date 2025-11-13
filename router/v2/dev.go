package v2

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/model"
	"github.com/hkyangyi/newe/router/app"
)

func RegisterDevRoutes(g *gin.RouterGroup) {
	g.GET("dbtables", GetDbTables)
	g.POST("dbtableadd", DbTableAdd)
	g.PUT("dbtableedit", DbTableEdit)
	g.DELETE("dbtabledel", DbTableDel)
	g.POST("dbtableInbysql", DbTableInbySql)
	g.GET("dbtablecolumns", GetDbTableColumns)
	g.POST("columadd", ColumnAdd)
	g.PUT("columedit", ColumnEdit)
	g.DELETE("columdel", ColumnDel)
	DevPageRegRoute(g.Group("page"))
}

func DbTableAdd(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevDbtable
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}
	if err := req.Add(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(req)
}

func DbTableInbySql(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevDbtable
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}

	if err := req.InbySql(); err != nil {
		g.Error(err)
		return
	}

	g.SUCCESS(req)

}
func DbTableEdit(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevDbtable
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}
	if err := req.Edit(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(req)
}

func DbTableDel(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevDbtable
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}
	if err := req.Del(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(req)
}

func GetDbTables(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevDbtable
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}

	var wheres []string
	var params []interface{}

	if req.DbTableName != "" {
		wheres = append(wheres, "name LIKE ?")
		params = append(params, "%"+req.DbTableName+"%")
	}

	where := strings.Join(wheres, " AND ")

	data := req.GetList(where, params...)
	g.SUCCESS(data)

}

func GetDbTableColumns(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevDbtable
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}
	if req.ID == "" {
		g.Error(errors.New("缺少参数 ID"))
		return
	}
	columns := req.GetColumns()
	g.SUCCESS(columns)
}

func ColumnAdd(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevDbtableColumns
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}
	if err := req.Add(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(req)
}

func ColumnEdit(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevDbtableColumns
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}
	if err := req.Edit(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(req)
}

func ColumnDel(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevDbtableColumns
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}
	if err := req.Del(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(req)
}
