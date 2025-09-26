package v2

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/router/app"
)

type verifyform struct {
	TableName  string `form:"tablename" valid:"Required; MaxSize(50)"`
	FieldName  string `form:"fieldname" valid:"Required; MaxSize(50)"`
	Tablevalue string `form:"tablevalue" valid:"Required; MaxSize(50)"`
	TableId    string `form:"tableid"`
}

func Verifysole(c *gin.Context) {
	var a = app.NewApp(c)
	var data verifyform
	if e := a.Bind(&data); e != nil {
		a.Error(errors.New("参数错误"))
		return
	}

	where := make(map[string]interface{})
	where[data.FieldName] = data.Tablevalue

	res := db.VerifyOnly(data.TableName, data.TableId, where)
	if res {
		a.Error(errors.New("已存在"))
		return
	}
	a.SUCCESS(nil)
}
