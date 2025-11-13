package middle

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/common/worklog"
	"github.com/hkyangyi/newe/cusmodel"
)

func CusAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")
		fmt.Println("token:", token)
		if token == "" {
			c.JSON(401, Response{
				Code:      10000,
				Message:   "登录超时1",
				Result:    nil,
				Success:   "fail",
				Timestamp: time.Now().Unix(),
			})
			c.Abort()
			return
		}

		// 兼容 Bearer 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		uuid, err := utils.AuthToken(token)

		if err != nil {
			worklog.Logio.WERR(fmt.Sprintf("JWT验证失败: %s, token: %s", err.Error(), token))
			c.JSON(401, Response{
				Code:      10000,
				Message:   "身份验证失败，请重新登录",
				Result:    nil,
				Success:   "fail",
				Timestamp: time.Now().Unix(),
			})
			c.Abort()
			return
		}
		fmt.Printf("middle AdminMid uuid：%s uid \n", uuid)
		//检测KEY是否在缓存中
		res := redis.REDIS.Exists(uuid)

		if !res {
			c.JSON(401, Response{
				Code:      10000,
				Message:   "登录超时2",
				Result:    nil,
				Success:   "fail",
				Timestamp: time.Now().Unix(),
			})
			c.Abort()
			return
		}
		//从缓存中拿取数据
		var data cusmodel.CusAuth
		errs := redis.REDIS.Get(uuid, &data)
		if errs != nil {
			c.JSON(401, Response{
				Code:      10000,
				Message:   "登录超时3",
				Result:    nil,
				Success:   "fail",
				Timestamp: time.Now().Unix(),
			})
			c.Abort()
			return
		}

		redis.REDIS.Set(uuid, data, 3600)
		c.Set("CusAdminAuthData", data)

		// 预采集请求体（仅非GET/HEAD且小于1MB），上传相关请求直接过滤不采集
		var bodyLog string
		if isCusUploadRequest(c) {
			bodyLog = "[filtered]"
		} else if m := c.Request.Method; m != "GET" && m != "HEAD" {
			cl := c.Request.ContentLength
			if cl >= 0 && cl <= 1<<20 { // 1MB上限
				if b, err := io.ReadAll(c.Request.Body); err == nil {
					bodyLog = string(b)
					c.Request.Body = io.NopCloser(bytes.NewBuffer(b)) // 还原给后续handler使用
				}
			}
		}
		c.Set("LOG_REQ_BODY", bodyLog)

		// 操作日志自动写入
		c.Next() // 先执行后续业务，保证能拿到响应码

		// 日志内容收集
		var log cusmodel.CustomerOperationLog
		log.Cid = data.CDB.Id
		if data.Isadmin {

			log.MemberID = data.CDB.Id
			log.Username = data.CDB.Username
		} else {
			log.MemberID = data.MerDb.ID
			log.Username = data.MerDb.Username
		}

		log.Method = c.Request.Method
		log.RoutePattern = c.FullPath() // 路由模板，如 /api/v1/user/:id

		log.ActualPath = c.Request.URL.Path // 实际请求路径
		// 分组名可用路由前缀或自定义
		log.GroupName = getGroupName(c.FullPath())
		log.IP = c.ClientIP()
		// 请求参数（GET用Query，其他用预采集的Body样本）
		if isCusUploadRequest(c) {
			log.Params = "[filtered]"
		} else if c.Request.Method == "GET" {
			log.Params = c.Request.URL.RawQuery
		} else if v, ok := c.Get("LOG_REQ_BODY"); ok {
			if s, ok2 := v.(string); ok2 {
				log.Params = s
			}
		}
		// 操作结果与响应码
		log.StatusCode = c.Writer.Status()
		if log.StatusCode >= 200 && log.StatusCode < 400 {
			log.Result = "success"
		} else {
			log.Result = "fail"
		}
		// 可补充 message 字段（如错误信息）
		// 先确保接口存在于 sys_api_list，并获取 api_id
		apiName := ""
		if log.Method == "GET" {
			apiName = "GET " + log.RoutePattern
		} else {
			apiName = log.Method + " " + log.RoutePattern
		}
		if apiID, err := (&cusmodel.CustomerOperationLog{RoutePattern: log.RoutePattern, Method: log.Method, GroupName: log.GroupName}).EnsureApiAndGetID(apiName); err == nil {
			log.ApiID = apiID
		}
		// 日志入库
		log.CreatedAt = time.Now()
		_ = log.Add()

		// ...existing code...
		// ...existing code...

	}
}

// isUploadRequest 判断是否上传相关请求，避免记录大二进制/敏感内容
func isCusUploadRequest(c *gin.Context) bool {
	ct := c.GetHeader("Content-Type")
	if strings.Contains(strings.ToLower(ct), "multipart/form-data") {
		return true
	}
	p := strings.ToLower(c.Request.URL.Path)
	fp := strings.ToLower(c.FullPath())
	if strings.Contains(p, "/upload") || strings.Contains(fp, "/upload") {
		return true
	}
	if strings.Contains(p, "uploadavatar") || strings.Contains(fp, "uploadavatar") {
		return true
	}
	return false
}
