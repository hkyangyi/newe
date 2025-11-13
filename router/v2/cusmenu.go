package v2

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/model"
	"github.com/hkyangyi/newe/router/app"
)

// RegisterCustomerMenusRoutes 注册菜单相关路由
func RegisterCustomerMenusRoutes(g *gin.RouterGroup) {

	g.POST("add", cusmenuCreate)
	g.PUT("edit", cusmenuUpdate)
	g.DELETE("del", cusmenuDelete)
	g.GET("tree", cusmenuTree)
	g.GET("buttons", cusmenuButtons)
}

func cusmenuCreate(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.CustomerMenus
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	if err := a.Add(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(a)
}

// 不使用统一 Bind 校验版本的更新：仅解析 JSON
func cusmenuUpdate(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.CustomerMenus
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	fmt.Println(a)

	if err := a.Edit(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(a)
}

func cusmenuDelete(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.CustomerMenus
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	if err := a.Del(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS("删除成功")
}

func cusmenuTree(c *gin.Context) {
	var g = app.NewApp(c)
	tree := model.CustomerMenuTree()
	g.SUCCESS(tree)
}

func cusmenuButtons(c *gin.Context) {
	var g = app.NewApp(c)
	var items []model.CustomerMenus
	db.Db.Where("is_deleted = 1 AND status = 1 AND type = 3").Order("order_no asc").Find(&items)
	perms := make([]string, 0, len(items))
	for _, v := range items {
		if v.Permission != "" {
			perms = append(perms, v.Permission)
		}
	}
	g.SUCCESS(perms)
}

// 可选：异步重建全部 full_path 与 level
func RebuildCusMenuFullPath() error {
	var items []model.CustomerMenus
	if err := db.Db.Where("is_deleted=1").Order("pid asc").Find(&items).Error; err != nil {
		return err
	}
	// 简单多轮迭代，避免递归多次查询
	for r := 0; r < 5; r++ { // 5层足够，一般不会很深
		changed := false
		for i := range items {
			if items[i].Pid == "" || items[i].Pid == "0" {
				items[i].FullPath = items[i].Path
				items[i].Level = 0
				continue
			}
			var parent *model.CustomerMenus
			for j := range items {
				if items[j].ID == items[i].Pid {
					parent = &items[j]
					break
				}
			}
			if parent != nil && parent.FullPath != "" {
				newFull := strings.TrimRight(parent.FullPath, "/") + items[i].Path
				if newFull != items[i].FullPath || items[i].Level != parent.Level+1 {
					items[i].FullPath = newFull
					items[i].Level = parent.Level + 1
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	now := time.Now().Unix()
	for _, it := range items {
		db.Db.Model(&model.CustomerMenus{}).Where("id=?", it.ID).Updates(map[string]interface{}{"full_path": it.FullPath, "level": it.Level, "updated_at": now})
	}
	return nil
}
