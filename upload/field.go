package upload

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/hkyangyi/newe/common/config"
	"github.com/hkyangyi/newe/common/utils"
)

// 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

// 获取文件类型
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// 查文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// 检查文件是否有权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// 如果目录不存在，则创建一个目录
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// 创建目录
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// 创建文件
func MkAll(src, filename string) error {
	fmt.Println(src, filename)

	err := os.MkdirAll(src, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	err = os.WriteFile(src+"/"+filename, []byte(""), os.ModePerm)
	if err != nil {
		return err
	}
	return nil

}

// 根据特定模式打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// 打开文件
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to openfile:%v", err)
	}
	return f, nil
}

type Field struct {
	Field      multipart.File
	Header     *multipart.FileHeader
	Size       int64
	Path       string
	PathPrefix string
	Url        string
	Name       string
	Ext        string
}

func NewFiled(f multipart.File, h *multipart.FileHeader, prefix string) *Field {
	var file = Field{
		Field:      f,
		Header:     h,
		PathPrefix: prefix,
	}
	return &file
}

func (f *Field) AcWriteImage() error {
	now := time.Now().Format("20060102")
	mdsrc := "/" + f.PathPrefix + "/" + now + "/"
	f.Path = config.Conf.IMG_SavePath + mdsrc

	uuid := utils.GetUUID()
	f.Ext = path.Ext(f.Header.Filename)
	if f.Ext == "" {
		f.Ext = getFileExtension(f.Header.Header.Get("Content-Type"))
	}

	imgname := uuid + f.Ext
	f.Name = imgname
	f.Url = config.Conf.HTTP_ServeUrl + config.Conf.IMG_PrefixUrl + mdsrc + f.Name
	if !f.CheckImageExt() {
		return fmt.Errorf("不允许的图片格式：%v", f.Ext)
	}
	if !f.CheckImageSize() {
		return fmt.Errorf("图片大小超出限制：%v", f.Size)
	}

	src, err := f.Header.Open()
	if err != nil {
		return fmt.Errorf("读取缓存出错：%v", err)
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(f.Path+f.Name), 0750); err != nil {
		return fmt.Errorf("创建文件路径出错：%v", err)
	}

	out, err := os.Create(f.Path + f.Name)
	if err != nil {
		return fmt.Errorf("创建文件出错：%v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return fmt.Errorf("文件写入出错%v", err)
	}

	return nil

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
	if config.Conf.IMG_AllowExts == "" {
		return true
	}
	for _, allowExt := range strings.Split(config.Conf.IMG_AllowExts, ",") {
		if strings.ToUpper(allowExt) == strings.ToUpper(a.Ext) {
			return true
		}
	}
	return false
}

// 检查图片大小
func (a *Field) CheckImageSize() bool {
	a.Size = a.Header.Size
	size := int(a.Size)
	return size <= config.Conf.IMG_MaxSize
}

// 检查图
func (a *Field) CheckImage() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = IsNotExistMkDir(dir + "/" + a.Path)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := CheckPermission(a.Path)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", a.Path)
	}
	return nil
}

func (f *Field) AcWriteFile() error {
	now := time.Now().Format("20060102")
	mdsrc := "/" + f.PathPrefix + "/" + now + "/"
	f.Path = config.Conf.FILE_SavePath + mdsrc
	uuid := utils.GetUUID()
	f.Ext = path.Ext(f.Header.Filename)
	if f.Ext == "" {
		f.Ext = getFileExtension(f.Header.Header.Get("Content-Type"))
	}

	filename := uuid + f.Ext
	f.Name = filename
	f.Url = config.Conf.HTTP_ServeUrl + config.Conf.FILE_PrefixUrl + mdsrc + f.Name
	if !f.CheckFileExt() {
		return fmt.Errorf("不允许的文件格式：%v", f.Ext)
	}
	if !f.CheckFieldSize() {
		return fmt.Errorf("文件大小超出限制：%v", f.Size)
	}

	src, err := f.Header.Open()
	if err != nil {
		return fmt.Errorf("读取缓存出错：%v", err)
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(f.Path+f.Name), 0750); err != nil {
		return fmt.Errorf("创建文件路径出错：%v", err)
	}

	out, err := os.Create(f.Path + f.Name)
	if err != nil {
		return fmt.Errorf("创建文件出错：%v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return fmt.Errorf("文件写入出错%v", err)
	}

	return nil

}

func (a *Field) CheckFileExt() bool {
	if config.Conf.FILE_AllowExts == "" {
		return true
	}
	for _, allowExt := range strings.Split(config.Conf.FILE_AllowExts, ",") {
		if strings.ToUpper(allowExt) == strings.ToUpper(a.Ext) {
			return true
		}
	}
	return false
}

func (a *Field) CheckFieldSize() bool {
	a.Size = a.Header.Size
	size := int(a.Size)
	return size <= config.Conf.FILE_MaxSize
}

// 根据 Content-Type 获取文件扩展名
func getFileExtension(contentType string) string {
<<<<<<< HEAD
=======
	fmt.Println(contentType)
>>>>>>> 542fd6c8f893697fb398c9524cb99d9dc41f8971
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
<<<<<<< HEAD
=======
	case "application/vnd.ms-excel":
		return ".xls"
>>>>>>> 542fd6c8f893697fb398c9524cb99d9dc41f8971
	// 其他文件类型的处理逻辑
	default:
		return ""
	}
}
