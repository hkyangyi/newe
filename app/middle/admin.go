package middle

import (
	"fmt"
	"newe/app/common/system/moddle"
	"newe/common/base"
	"newe/common/utils"
	"time"

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
		uuid, err := utils.AuthToken(token)

		if !err {
			base.WorkLog.WERR(err)
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
		fmt.Printf("middle AdminMid uuid：%s uid \n", uuid)
		//检测KEY是否在缓存中
		res := base.REDIS.Exists(uuid)

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
		var data moddle.SysMember
		errs := base.REDIS.Get(uuid, &data)
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
		base.REDIS.Set(uuid, data, 3600)
		c.Set("AdminAuthData", data)

		c.Next()
	}
}
