package cus

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
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

// RegisterUploadRoutes 注册上传路由
// POST /file  表单字段 name=file
func RegisterUploadRoutes(g *gin.RouterGroup) {
	//g.POST("upload", UploadFile)
	g.GET("filespage", getFilePage)
}

// UploadFile 依据配置保存文件（图片/文件），记录到 sys_upload
func UploadFile(c *gin.Context) {
	a := app.NewApp(c)

	// 登录用户信息
	auth, ok := a.C.Get("CusAdminAuthData")
	if !ok {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := auth.(cusmodel.CusAuth)
	if err := mer.RefreshByMerdb(); err != nil {
		a.LoginError(err)
		return
	}
	memberID := mer.MerDb.ID

	// 读取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		a.Error(err)
		return
	}
	defer file.Close()

	conf := config.Conf
	origName := header.Filename
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(origName)), ".")
	if ext == "" {
		ext = "bin"
	}

	// 判断是否图片（按扩展名集合）
	isImage := isInList(ext, conf.IMG_AllowExts)

	// 校验扩展名
	if isImage {
		if !isInList(ext, conf.IMG_AllowExts) {
			a.Error(errors.New("图片格式不允许"))
			return
		}
	} else {
		if conf.FILE_AllowExts != "" && !isInList(ext, conf.FILE_AllowExts) {
			a.Error(errors.New("文件格式不允许"))
			return
		}
	}

	// 保存目录与URL前缀
	datePath := time.Now().Format("2006/01/02")
	baseSave := conf.FILE_SavePath
	urlPrefix := conf.FILE_PrefixUrl
	if isImage {
		baseSave = conf.IMG_SavePath
		urlPrefix = conf.IMG_PrefixUrl
	}
	urlPrefix = strings.TrimRight(conf.HTTP_ServeUrl, "/") + strings.TrimRight(urlPrefix, "/")
	// 生成文件名
	saveName := genFileName(ext)
	saveDir := filepath.Join(baseSave, mer.CDB.Id, datePath)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		a.Error(err)
		return
	}
	savePath := filepath.Join(saveDir, saveName)

	// 文件大小限制（优先用 header.Size 预校验）
	maxSize := int64(conf.FILE_MaxSize)
	if isImage {
		maxSize = int64(conf.IMG_MaxSize)
	}
	if maxSize > 0 && header.Size > 0 && header.Size > maxSize {
		a.Error(errors.New("文件超过大小限制"))
		return
	}

	// 写入文件并计算哈希；图片走压缩流程（jpg/jpeg/png），其他按原样保存
	var (
		size    int64
		md5Str  string
		sha1Str string
	)
	// 判断是否压缩类型
	compressible := isImage && (ext == "jpg" || ext == "jpeg" || ext == "png")
	if compressible {
		// 读取全部内容到内存（受大小限制）
		buf, err := io.ReadAll(file)
		if err != nil {
			a.Error(err)
			return
		}
		if maxSize > 0 && int64(len(buf)) > maxSize {
			a.Error(errors.New("文件超过大小限制"))
			return
		}
		// 解码
		img, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			a.Error(err)
			return
		}
		// 压缩编码
		outBuf := &bytes.Buffer{}
		switch ext {
		case "jpg", "jpeg":
			// 质量 80，可按需调整
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
		}
		// 写入文件
		out, err := os.Create(savePath)
		if err != nil {
			a.Error(err)
			return
		}
		if _, err := out.Write(outBuf.Bytes()); err != nil {
			_ = out.Close()
			a.Error(err)
			return
		}
		_ = out.Close()
		// 计算哈希、大小
		hMd5 := md5.Sum(outBuf.Bytes())
		hSha1 := sha1.Sum(outBuf.Bytes())
		md5Str = hex.EncodeToString(hMd5[:])
		sha1Str = hex.EncodeToString(hSha1[:])
		size = int64(outBuf.Len())
		// 覆盖 header.Size 用于后续逻辑
		header.Size = size
	} else {
		// 非压缩类型：直接落盘并计算哈希
		out, err := os.Create(savePath)
		if err != nil {
			a.Error(err)
			return
		}
		defer out.Close()
		hMd5 := md5.New()
		hSha1 := sha1.New()
		mw := io.MultiWriter(out, hMd5, hSha1)
		if _, err := io.Copy(mw, file); err != nil {
			a.Error(err)
			return
		}
		md5Str = hex.EncodeToString(hMd5.Sum(nil))
		sha1Str = hex.EncodeToString(hSha1.Sum(nil))
		if fi, err := out.Stat(); err == nil {
			size = fi.Size()
		} else {
			size = header.Size
		}
	}
	// 再次大小校验（基于最终文件）
	if maxSize > 0 && size > maxSize {
		_ = os.Remove(savePath)
		a.Error(errors.New("文件超过大小限制"))
		return
	}

	// MIME 检测
	mimeType := "application/octet-stream"
	if f2, err := os.Open(savePath); err == nil {
		defer f2.Close()
		buf := make([]byte, 512)
		if n, _ := f2.Read(buf); n > 0 {
			mimeType = http.DetectContentType(buf[:n])
		}
	}

	// 若是图片，获取宽高
	width, height := 0, 0
	if strings.HasPrefix(mimeType, "image/") || isImage {
		if f3, err := os.Open(savePath); err == nil {
			defer f3.Close()
			if cfg, _, err := image.DecodeConfig(f3); err == nil {
				width, height = cfg.Width, cfg.Height
				isImage = true
			}
		}
	}

	// 统一 URL，避免重复斜杠
	url := strings.TrimRight(urlPrefix, "/") + "/" + filepath.ToSlash(filepath.Join(mer.CDB.Id, datePath, saveName))

	// 记录入库
	rec := cusmodel.CustomerUpload{
		ID:           utils.GetUUID(),
		Cid:          mer.MerDb.Cid,
		MemberID:     memberID,
		OriginalName: origName,
		FileName:     saveName,
		Ext:          ext,
		MimeType:     mimeType,
		Size:         size,
		Md5:          md5Str,
		Sha1:         sha1Str,
		Path:         filepath.ToSlash(filepath.Join(baseSave, datePath, saveName)),
		URL:          url,
		IsImage:      map[bool]int{true: 1, false: -1}[isImage],
		Width:        width,
		Height:       height,
		Status:       1,
		CreatedAt:    time.Now().Unix(),
	}
	if err := db.Db.Create(&rec).Error; err != nil {
		a.Error(err)
		return
	}

	a.SUCCESS(gin.H{
		"id":           rec.ID,
		"memberId":     rec.MemberID,
		"originalName": rec.OriginalName,
		"fileName":     rec.FileName,
		"ext":          rec.Ext,
		"mimeType":     rec.MimeType,
		"size":         rec.Size,
		"md5":          rec.Md5,
		"sha1":         rec.Sha1,
		"path":         rec.Path,
		"url":          rec.URL,
		"isImage":      rec.IsImage == 1,
		"width":        rec.Width,
		"height":       rec.Height,
		"createdAt":    rec.CreatedAt,
	})
}

// 工具：检查扩展名是否在逗号分隔列表中（不区分大小写）
func isInList(ext string, csv string) bool {
	if csv == "" {
		return true
	}
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	parts := strings.Split(csv, ",")
	for _, p := range parts {
		if strings.TrimSpace(strings.ToLower(p)) == ext {
			return true
		}
	}
	return false
}

// 生成文件名：时间戳+随机串
func genFileName(ext string) string {
	ts := time.Now().UnixNano()
	rnd := sha1.Sum([]byte(strings.Join([]string{time.Now().Format("150405"), "-", hex.EncodeToString([]byte{byte(ts), byte(ts >> 8), byte(ts >> 16)})}, "")))
	name := hex.EncodeToString(rnd[:])[:20]
	if ext != "" {
		return name + "." + ext
	}
	return name
}

type FileQuery struct {
	FileType  int    `form:"fileType,default=0"` // 0 全部 1 图片 2 文件
	StartTime string `form:"startTime"`          // 上传开始时间（Unix时间戳）
	EndTime   string `form:"endTime"`            // 上传结束时间（Unix时间戳）
	Keyword   string `form:"keyword"`            // 原始文件名（含扩展名）

	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=10"`
}

// 获取分页
func getFilePage(c *gin.Context) {
	var g = app.NewApp(c)
	var req FileQuery
	// 鉴权
	merdata, ok := g.C.Get("CusAdminAuthData")
	if !ok {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)
	if err := mer.RefreshByMerdb(); err != nil {
		g.LoginError(err)
		return
	}

	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}

	var where []string
	var params []interface{}
	where = append(where, " member_id = ? and status = 1 and cid = ? ")
	params = append(params, mer.MerDb.ID, mer.MerDb.Cid)
	if req.FileType == 1 {
		where = append(where, " is_image = 1 ")
	}
	if req.FileType == 2 {
		where = append(where, " is_image = -1 ")
	}

	rangetime := utils.RemoveSliceEmpty([]string{req.StartTime, req.EndTime})
	if len(rangetime) == 2 {
		startTime, err1 := utils.ParseDateStringToTimestamp(rangetime[0], "2006-01-02")
		endTime, err2 := utils.ParseDateStringToTimestamp(rangetime[1], "2006-01-02")
		if err1 == nil && err2 == nil {
			where = append(where, " created_at >= ? ")
			params = append(params, startTime)
			where = append(where, " created_at <= ? ")
			params = append(params, endTime+86399)
		}
	}

	if req.Keyword != "" {
		where = append(where, " original_name LIKE ? ")
		params = append(params, "%"+req.Keyword+"%")
	}
	whereStr := strings.Join(where, " AND ")
	page := utils.PageList{Page: req.Page, PageSize: req.PageSize}
	var sup cusmodel.CustomerUpload
	items, _ := sup.GetPage(page, whereStr, params...)

	g.SUCCESS(items)

}
