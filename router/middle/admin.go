package middle

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hkyangyi/newe/common/config"
	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/common/worklog"
	"github.com/hkyangyi/newe/model"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
	Success   string      `json:"success"`
	Timestamp int64       `json:"timestamp"`
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// fmt.Println(path)
		// // //index:= strings.Index(path,"/api")
		// if strings.Index(path, "/") == 0 {
		// 	path = path[1:]
		// 	fmt.Println(path)
		// }
		// patharr := strings.Split(path, "/")
		// if len(patharr) > 0 {
		// 	fmt.Println(patharr)
		// 	pathgroup:= strings.Join(path)
		// 	apidata:= model.SysApi{
		// 		Path: path,
		// 		PathGroup: patharr[0],
		// 	}
		// }

		//Authorization=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTgzNTI2MTcsImlhdCI6MTc1ODA5MzQxNywidXVpZCI6IkFVVEhfNzJiZjJjYWQ3ZjZmNGEzZjhhY2Q4ZDQ4MDg4ODUyOTNfYWRtaW4ifQ.Xvor8lZ4iqueCAtMs1pj_18UvTuIdNflnlBShfIvO_c

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
		var data model.AdminAuth
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
		//data.Refresh()
		if config.Conf.SYS_SQLAUTH == 1 && data.Merdb.Username != "admin" {
			// 可在此调用 SqlRuleVerify 检查接口权限，当前跳过
		}
		redis.REDIS.Set(uuid, data, 3600)
		c.Set("AdminAuthData", data)

		// 预采集请求体（仅非GET/HEAD且小于1MB），以免影响后续绑定
		var bodyLog string
		if m := c.Request.Method; m != "GET" && m != "HEAD" {
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
		var log model.SysOperationLog
		log.MemberID = data.Merdb.ID
		log.Username = data.Merdb.Username
		log.Method = c.Request.Method
		log.RoutePattern = c.FullPath()     // 路由模板，如 /api/v1/user/:id
		log.ActualPath = c.Request.URL.Path // 实际请求路径
		// 分组名可用路由前缀或自定义
		log.GroupName = getGroupName(c.FullPath())
		log.IP = c.ClientIP()
		// 请求参数（GET用Query，其他用预采集的Body样本）
		if c.Request.Method == "GET" {
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
		if apiID, err := (&model.SysOperationLog{RoutePattern: log.RoutePattern, Method: log.Method, GroupName: log.GroupName}).EnsureApiAndGetID(apiName); err == nil {
			log.ApiID = apiID
		}
		// 日志入库
		log.CreatedAt = time.Now()
		_ = log.Add()

		// ...existing code...
		// ...existing code...

	}
}

// getGroupName 从路由模板提取分组名（如 /api/v1/user/:id → user）
func getGroupName(path string) string {
	arr := strings.Split(path, "/")
	// 去除所有 /:xxx 路径段
	var clean []string
	for _, v := range arr {
		if v != "" && !strings.HasPrefix(v, ":") {
			clean = append(clean, v)
		}
	}
	// 取倒数第二个为分组名
	if len(clean) >= 2 {
		return clean[len(clean)-2]
	}
	return ""
}

// // 数据权限中间件
// func SqlRuleVerify(c *gin.Context, data model.AdminAuth) error {
// 	//请求路径
// 	method := c.Request.Method
// 	path := c.Request.URL.Path
// 	fmt.Println(path)
// 	fmt.Println(method)
// 	apidata := GetApiData(path)
// 	fmt.Println("apidata:", apidata)
// 	if len(apidata.ID) != 32 {
// 		return nil
// 	}

// 	if apidata.Status != 1 {
// 		return errors.New("该接口已停用")
// 	}

// 	ruledb := GetRoleRules(apidata.ID, data.Merdb.RoleId)
// 	fmt.Println(ruledb)

// 	if len(ruledb.ID) != 32 {
// 		return errors.New("您未添加该接口权限")
// 	}

// 	if ruledb.Status != 1 {
// 		return errors.New("您已停用该接口权限")
// 	}
// 	return nil
// }

//更具PATH获取API数据
