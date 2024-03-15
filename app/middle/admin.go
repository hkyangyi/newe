package middle

import (
	"errors"
	"fmt"
	"time"

	"github.com/hkyangyi/newe/app/common/system/moddle"
	"github.com/hkyangyi/newe/common/config"
	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/common/worklog"

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
		token := c.GetHeader("Authorization")
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
		uuid, err := utils.AuthToken(token)

		if err != nil {
			worklog.Logio.WERR(err)
			c.JSON(401, Response{
				Code:      10000,
				Message:   err.Error(),
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
				Message:   "登录超时",
				Result:    nil,
				Success:   "fail",
				Timestamp: time.Now().Unix(),
			})
			c.Abort()
			return
		}
		//从缓存中拿取数据
		var data moddle.AdminAuth
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
			err = SqlRuleVerify(c, data)
			if err != nil {
				c.JSON(200, Response{
					Code:      400,
					Message:   err.Error(),
					Result:    nil,
					Success:   "fail",
					Timestamp: time.Now().Unix(),
				})
				c.Abort()
				return
			}
		}
		redis.REDIS.Set(uuid, data, 3600)
		c.Set("AdminAuthData", data)
		c.Next()
	}
}

// 数据权限中间件
func SqlRuleVerify(c *gin.Context, data moddle.AdminAuth) error {
	//请求路径
	method := c.Request.Method
	path := c.Request.URL.Path
	fmt.Println(path)
	fmt.Println(method)
	apidata := GetApiData(path)
	fmt.Println("apidata:", apidata)
	if len(apidata.ID) != 32 {
		return nil
	}

	if apidata.Status != 1 {
		return errors.New("该接口已停用")
	}

	ruledb := GetRoleRules(apidata.ID, data.Merdb.RoleId)
	fmt.Println(ruledb)

	if len(ruledb.ID) != 32 {
		return errors.New("您未添加该接口权限")
	}

	if ruledb.Status != 1 {
		return errors.New("您已停用该接口权限")
	}
	return nil
}

//更具PATH获取API数据
