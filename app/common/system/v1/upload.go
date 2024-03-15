package v1

import (
	"errors"

	"github.com/hkyangyi/newe/app/common/app"
	"github.com/hkyangyi/newe/app/common/system/moddle"
	"github.com/hkyangyi/newe/app/common/upload"

	"github.com/gin-gonic/gin"
)

func UpImage(c *gin.Context) {

	var g = app.NewApp(c)

	file, image, err := c.Request.FormFile("file")
	if err != nil {
		g.Error(errors.New("无法找到文件参数"))
		return
	}

	if image == nil {
		g.Error(errors.New("文件获取失败"))
		return
	}
	authdb, _ := c.Get("AdminAuthData")
	authdata := authdb.(moddle.AdminAuth)
	var f = upload.Field{
		Usdb:   authdata.Merdb,
		Field:  file,
		Header: image,
	}
	path := f.GetPath()
	name := f.GetName()
	src := path + name
	if !f.CheckImageSize() {
		g.Error(errors.New("格式不正确"))
		return
	}
	if !f.CheckImageSize() {
		g.Error(errors.New("文件大小超出限制"))
		return
	}
	err = f.CheckImage()
	if err != nil {
		g.Error(err)
		return
	}

	fullurl := f.GetUrl()

	if err := c.SaveUploadedFile(image, src); err != nil {
		g.Error(errors.New("图片保存失败"))
		return
	}
	f.ImagesSave()
	data := make(map[string]interface{})
	data["url"] = fullurl
	data["path"] = src
	g.SUCCESS(data)
}
