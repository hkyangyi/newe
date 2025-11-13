package cus

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/config"
	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/cusmodel"
	"github.com/hkyangyi/newe/router/app"
)

func RegisterUserRoutes(g *gin.RouterGroup) {
	g.PUT("editProfile", EditUserProfile)
	g.PUT("changePassword", ChangePassword)
	g.GET("messages", GetUserMessages)
	g.POST("uploadAvatar", UploadAvatar)
	g.PUT("readMessage", MessageReaded)
}

func GetUserInfo(c *gin.Context) {
	var a = app.NewApp(c)

	merdata, b := a.C.Get("CusAdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}

	fmt.Println(merdata)

	mer := merdata.(cusmodel.CusAuth)
	err := mer.RefreshByMerdb()
	if err != nil {
		a.LoginError(err)
		return
	}
	a.SUCCESS(mer.MerDb)
}

// 修改密码
func ChangePassword(c *gin.Context) {
	var a = app.NewApp(c)

	// 获取当前用户
	merdata, ok := a.C.Get("CusAdminAuthData")
	if !ok {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	if err := mer.RefreshByMerdb(); err != nil {
		a.LoginError(err)
		return
	}

	// 绑定入参
	var req struct {
		OldPassword string `json:"oldPassword" form:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" form:"newPassword" binding:"required"`
	}
	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}

	if err := cusmodel.MemberChangePassword(mer.MerDb.ID, req.OldPassword, req.NewPassword); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

// 编辑信息（头像、昵称）
func EditUserProfile(c *gin.Context) {
	var a = app.NewApp(c)
	var req cusmodel.CustomerMember

	merdata, ok := a.C.Get("CusAdminAuthData")
	if !ok {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	if err := mer.RefreshByMerdb(); err != nil {
		a.LoginError(err)
		return
	}

	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}

	req.ID = mer.MerDb.ID
	if err := req.EditProfile(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}

// 获取消息列表（分页）
func GetUserMessages(c *gin.Context) {
	var a = app.NewApp(c)
	var req cusmodel.CustomerMessage

	merdata, ok := a.C.Get("CusAdminAuthData")
	if !ok {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	if err := mer.RefreshByMerdb(); err != nil {
		a.LoginError(err)
		return
	}

	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}

	fmt.Println(req)

	var wheremap []string
	var params []interface{}
	wheremap = append(wheremap, " uid = ? ")
	params = append(params, mer.MerDb.ID)
	if req.MsgType != "" {
		wheremap = append(wheremap, " msg_type = ? ")
		params = append(params, req.MsgType)
	}

	where := strings.Join(wheremap, " AND ")
	page := utils.PageList{Page: req.Page, PageSize: req.PageSize}
	items, err := req.GetPage(page, where, params...)
	if err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(items)
}

// 头像上传
func UploadAvatar(c *gin.Context) {
	a := app.NewApp(c)
	// 鉴权
	merdata, ok := a.C.Get("CusAdminAuthData")
	if !ok {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	if err := mer.RefreshByMerdb(); err != nil {
		a.LoginError(err)
		return
	}

	// 绑定参数
	var req struct {
		ImageBase64 string `json:"imageBase64" form:"imageBase64" binding:"required"`
	}
	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}
	data := strings.TrimSpace(req.ImageBase64)
	if data == "" {
		a.Error(errors.New("缺少图片数据"))
		return
	}

	// 解析 dataURL 前缀（可选）
	var ext string
	if strings.HasPrefix(data, "data:image/") {
		// 形如 data:image/png;base64,xxxxx
		comma := strings.Index(data, ",")
		if comma > 0 {
			header := data[:comma]
			data = data[comma+1:]
			if strings.Contains(header, "image/png") {
				ext = "png"
			} else if strings.Contains(header, "image/jpeg") || strings.Contains(header, "image/jpg") {
				ext = "jpg"
			} else if strings.Contains(header, "image/gif") {
				// 为避免动图处理复杂，这里不支持 gif 作为头像
				a.Error(errors.New("暂不支持GIF头像，请使用 JPG/PNG"))
				return
			}
		}
	}

	// base64 解码
	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		// 可能是URL安全或缺少填充的 Base64
		raw, err = base64.RawStdEncoding.DecodeString(data)
		if err != nil {
			a.Error(errors.New("图片Base64格式不正确"))
			return
		}
	}

	// 尺寸限制（按解码后字节）
	conf := config.Conf
	maxSize := int64(conf.IMG_MaxSize)
	if maxSize > 0 && int64(len(raw)) > maxSize {
		a.Error(errors.New("图片超过大小限制"))
		return
	}

	// 解码为 image
	img, format, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		a.Error(errors.New("无法解析图片"))
		return
	}
	if ext == "" {
		if format == "png" {
			ext = "png"
		} else {
			ext = "jpg" // 默认使用 jpg
		}
	}

	// 编码压缩
	outBuf := &bytes.Buffer{}
	switch ext {
	case "jpg", "jpeg":
		if err := jpeg.Encode(outBuf, img, &jpeg.Options{Quality: 80}); err != nil {
			a.Error(err)
			return
		}
	case "png":
		enc := png.Encoder{CompressionLevel: png.BestCompression}
		if err := enc.Encode(outBuf, img); err != nil {
			a.Error(err)
			return
		}
	default:
		a.Error(errors.New("不支持的图片格式"))
		return
	}

	// 路径与文件名
	datePath := time.Now().Format("2006/01/02")
	baseSave := conf.IMG_SavePath
	saveDir := filepath.Join(baseSave, "avatar", datePath)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		a.Error(err)
		return
	}
	// 检查扩展名白名单
	if !isInList(ext, conf.IMG_AllowExts) {
		a.Error(errors.New("图片格式不允许"))
		return
	}
	saveName := genFileName(ext)
	savePath := filepath.Join(saveDir, saveName)
	if err := os.WriteFile(savePath, outBuf.Bytes(), 0644); err != nil {
		a.Error(err)
		return
	}

	// 更新头像URL
	urlPrefix := strings.TrimRight(conf.IMG_PrefixUrl, "/")
	houst := strings.TrimRight(conf.HTTP_ServeUrl, "/")
	url := houst + urlPrefix + "/" + filepath.ToSlash(filepath.Join("avatar", datePath, saveName))
	if err := db.Db.Model(&cusmodel.CustomerMember{}).Where("id = ?", mer.MerDb.ID).Update("headimgurl", url).Error; err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(gin.H{"url": url})
}

// 消息已读
func MessageReaded(c *gin.Context) {
	a := app.NewApp(c)
	var req cusmodel.CustomerMessage
	// 鉴权
	merdata, ok := a.C.Get("CusAdminAuthData")
	if !ok {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	if err := mer.RefreshByMerdb(); err != nil {
		a.LoginError(err)
		return
	}

	if err := a.Bind(&req); err != nil {
		a.Error(err)
		return
	}

	if len(req.ID) == 0 {
		a.Error(errors.New("缺少消息ID"))
		return
	}

	if err := req.SetRead(); err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(nil)
}
