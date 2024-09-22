package v1

import (
	"errors"

	"github.com/hkyangyi/newe/model"
	"github.com/hkyangyi/newe/router/app"
	"github.com/hkyangyi/newe/upload"

	"github.com/gin-gonic/gin"
)

// 公共上次，可以传图片可以传文件
func FileUpload(c *gin.Context) {
	var g = app.NewApp(c)
	authdb, _ := c.Get("AdminAuthData")
	authdata := authdb.(model.AdminAuth)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		g.Error(errors.New("无法找到文件参数"))
		return
	}
	f := upload.NewFiled(file, header, "admin/"+authdata.Merdb.Username)
	err = f.AcWriteFile()
	if err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(f.Url)

}

func UpImage(c *gin.Context) {
	var g = app.NewApp(c)
	authdb, _ := c.Get("AdminAuthData")
	authdata := authdb.(model.AdminAuth)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		g.Error(errors.New("无法找到文件参数"))
		return
	}
	f := upload.NewFiled(file, header, "admin/"+authdata.Merdb.Username)
	err = f.AcWriteImage()
	if err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(f.Url)
}
