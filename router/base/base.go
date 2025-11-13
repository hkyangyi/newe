package base

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/model"
	"github.com/hkyangyi/newe/router/app"
)

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
