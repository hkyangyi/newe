package v2

const DevPageModelTpl = `package model

import (
	"errors"
	"time"

	"github.com/hkyangyi/newe"
	"github.com/hkyangyi/newe/common/utils"
)
// {{.ModelName}} 结构体对应 {{.TableName}} 表

{{.StStr}}

func (a *{{.ModelName}}) RefResh() error {
	if a.ID == "" {
		return errors.New("缺少参数ID")
	}
	err := newe.MYDB.Model(a).First(a).Error
	return err
}

//新增
func (a *{{.ModelName}}) Add() error {
	a.Id = utils.GetUUID()
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = time.Now().Unix()
	err := newe.MYDB.Create(a).Error
	return err
}

//编辑
func (a *{{.ModelName}}) Edit() error {
	if len(a.Id) != 32 {
		return errors.New("缺少参数ID")
	}
	a.UpdateTime = time.Now().Unix()
	err := newe.MYDB.Updates(a).Error
	return err
}

//删除
func (a *{{.ModelName}}) Del() error {
	if len(a.Id) != 32 {
		return errors.New("缺少参数ID")
	}
	err := newe.MYDB.Model(a).Delete(a).Error
	return err
}

//获取分页
func (a *{{.ModelName}}) GetPage(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []{{.ModelName}}
	newe.MYDB.Model(&{{.ModelName}}{}).Where(where, v...).Count(&page.Total).Order("create_time desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	page.List = items
	return page
}

//获取列表
func (a *{{.ModelName}}) GetList() []{{.ModelName}} {
	var items []{{.ModelName}}
	newe.MYDB.Model(&{{.ModelName}}{}).Order("create_time desc").Find(&items)
	return items
}
`
const DevPageCorsTpl = `package admin

import (
	"newe-sjcj-serve/common/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/router/app"
)

func {{.ModelName}}RegRoute(r *gin.RouterGroup){
	r.POST("add",{{.ModelName}}Add)
	r.PUT("edit",{{.ModelName}}Edit)
	r.DELETE("del",{{.ModelName}}Del)
	r.GET("page",{{.ModelName}}GetPage)
	r.GET("list",{{.ModelName}}GetList)
	r.GET("tree",{{.ModelName}}GetTree)
}


// 添加
func {{.ModelName}}Add(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.{{.ModelName}}
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	if err := a.Add(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}

// 编辑
func {{.ModelName}}Edit(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.{{.ModelName}}
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	if len(a.Id) != 32 {
		g.Error(errors.New("缺少参数ID"))
		return
	}

	if err := a.Edit(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}

// 删除
func {{.ModelName}}Del(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.{{.ModelName}}
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	if err := a.Del(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}

// 获取分页列表
func {{.ModelName}}GetPage(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.{{.ModelName}}
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	var params []interface{}
	var where []string

	{{.WhereStr}}

	page := utils.PageList{
		Page:     a.Page,
		PageSize: a.PageSize,
	}
	ws := strings.Join(where, " AND ")
	res := a.GetPage(page, ws, params...)
	g.SUCCESS(res)
}

// 获取所有列表
func {{.ModelName}}GetList(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.{{.ModelName}}
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	res := a.GetList()
	g.SUCCESS(res)
}

//获取树形
func {{.ModelName}}GetTree(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.{{.ModelName}}
	res := a.GetList()
	tree := utils.ListToTree(res, "id", "pid", "children")
	g.SUCCESS(tree)
}
`
const DevPageVueApiTpl = `import type { {{.ModelName}} } from './data';

import { confirm } from '@vben/common-ui';

import { requestClient } from '#/api/request';

// 新增
export const Add = (params: Partial<{{.ModelName}}>) =>
  requestClient.post('{{.RootPath}}/add', params);

// 编辑
export const Edit = (params: Partial<{{.ModelName}}>) =>
  requestClient.put('{{.RootPath}}/edit', params);

// 获取分页
export const GetPage = (params: any) =>
  requestClient.get('{{.RootPath}}/page', { params });

// 获取列表
export const GetList = (params: any) =>
  requestClient.get('{{.RootPath}}/list', { params });

// 获取树形结构
export const GetTree = (params: any) =>
  requestClient.get('{{.RootPath}}/tree', { params });


// 删除
export const Del = (params: { id?: string; ids?: string[] }) =>
  confirm({
    content: '是否确认删除所选数据？',
    icon: 'warning',
  })
    .then(() => {
      return requestClient.delete('{{.RootPath}}/del', { data: params });
    })
    .catch(() => {});
`
