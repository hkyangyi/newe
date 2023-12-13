package upload

import (
	"fmt"
	"math/rand"
	"mime/multipart"
	"newe/app/common/system/moddle"
	"newe/common/base"
	"newe/common/file"
	"newe/common/utils"
	"os"
	"path"
	"strings"
	"time"
)

type Field struct {
	Usdb   moddle.SysMember
	Field  multipart.File
	Header *multipart.FileHeader
	Size   int64
	Path   string
	Url    string
	Name   string
}

// 生成图片路径
func (a *Field) GetPath() string {
	now := time.Now().Format("20060102")
	a.Path = base.Conf.IMG_SavePath + "/admin" + "/" + now + "/"
	return a.Path
}

func (a *Field) GetUrl() string {
	now := time.Now().Format("20060102")
	a.Url = base.Conf.IMG_PrefixUrl + "/admin" + "/" + now + "/" + a.Name
	return a.Url
}

// 获取图片名称
func (a *Field) GetName() string {
	//a.Name = a.Image.Filename
	uuid := utils.GetUUID()
	ext := path.Ext(a.Header.Filename)
	imgname := uuid + ext
	a.Name = imgname
	return imgname
}

// 获取文件格式
func (a *Field) GetType() string {
	ext := path.Ext(a.Header.Filename)
	return ext
}

// 随机生成符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 检查图片后缀
func (a *Field) CheckImageExt() bool {
	ext := file.GetExt(a.Header.Filename)
	for _, allowExt := range strings.Split(base.Conf.IMG_AllowExts, ",") {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

// 检查图片大小
func (a *Field) CheckImageSize() bool {
	a.Size = a.Header.Size
	size := int(a.Size)
	return size <= base.Conf.IMG_MaxSize
}

// 检查图
func (a *Field) CheckImage() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + a.Path)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(a.Path)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", a.Path)
	}
	return nil
}

//保存到数据库

func (a *Field) ImagesSave() {
	var data = moddle.SysImages{
		Path:       a.Path + a.Name,
		URL:        a.Url,             //
		Size:       a.Size,            //
		CreateTime: time.Now().Unix(), //
		CreateName: a.Usdb.Realname,   //
		DepartId:   a.Usdb.DepartId,   //
	}
	moddle.ImagesAdd(data)
}

// 生成图片路径
func (a *Field) GetFieldPath() string {
	now := time.Now().Format("20060102")
	a.Path = base.Conf.IMG_SavePath + "/" + a.Usdb.ID + "/" + now + "/"
	a.Url = base.Conf.IMG_PrefixUrl + "/" + a.Usdb.ID + "/" + now + "/" + a.Name
	return a.Path
}

func (a *Field) GetFieldName() string {
	//a.Name = a.Image.Filename
	uuid := utils.GetUUID()
	fieldname := uuid + "_" + a.Header.Filename
	a.Name = fieldname
	return fieldname
}

func (a *Field) CheckFieldSize() bool {
	a.Size = a.Header.Size
	size := int(a.Size)
	return size <= base.Conf.IMG_MaxSize
}
