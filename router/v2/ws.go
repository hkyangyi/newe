package v2

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/common/ws"
	"github.com/hkyangyi/newe/model"
)

func RegisterWsRoutes(r *gin.RouterGroup) {
	//注册ws路由管理器

	// 单一路由示例：路径/查询/Header 三处取 token；校验后将连接加入指定分组
	r.GET(":token", func(c *gin.Context) {
		// 1) 提取 token：优先路径参数，其次 query，再次 Authorization: Bearer
		token := c.Param("token")
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			auth := c.GetHeader("Authorization")
			if len(auth) > 7 && strings.EqualFold(auth[:6], "Bearer") {
				token = strings.TrimSpace(auth[7:])
			}
		}
		if token == "" {
			c.Status(400)
			return
		}

		// 2) 校验 JWT，解析 uuid
		uuid, err := utils.AuthToken(token)
		if err != nil || uuid == "" {
			c.Status(401)
			return
		}

		//根据uuid获取用户信息
		//检测KEY是否在缓存中
		res := redis.REDIS.Exists(uuid)

		if !res {
			c.Status(401)
			return
		}
		//从缓存中拿取数据
		var data model.AdminAuth
		errs := redis.REDIS.Get(uuid, &data)
		if errs != nil {
			c.Status(401)
			return
		}
		newesysws := ws.MainHub.Get("NeweSysAdmin")
		// 3) 生成连接 key，并接入分组。此处固定使用 "MAIN" 分组，可按业务扩展
		key := "Ws_NeweSysAdmin_" + data.Merdb.ID
		if err := newesysws.Serve(c.Writer, c.Request, key, func(cli *ws.Client, msg []byte) {
			ReadMessage(cli, msg)
		}); err != nil {
			c.Status(500)
			return
		}
		go InitMessage(key)
	})
}

// 回调消息处理
func ReadMessage(li *ws.Client, msg []byte) {
	fmt.Println("ws message recv:", string(msg))
}

// 初始化消息
func InitMessage(key string) {
	time.Sleep(2 * time.Second)
	newesysws := ws.MainHub.Get("NeweSysAdmin")
	uid := strings.TrimPrefix(key, "Ws_NeweSysAdmin_")
	//查看前10条未读消息
	var req model.SysMessage
	var wheremap []string
	var params []interface{}
	wheremap = append(wheremap, " uid = ? ")
	params = append(params, uid)
	wheremap = append(wheremap, " read_at = 0 ")

	where := strings.Join(wheremap, " AND ")
	page := utils.PageList{Page: 1, PageSize: 5}
	items, err := req.GetPage(page, where, params...)

	if err != nil {
		fmt.Println("InitMessage err:", err)
		return
	}
	var sedmsg ws.Frame
	sedmsg.Type = "initMessage"
	sedmsg.Data = items.List

	//发送未读消息
	newesysws.SendJSON(key, sedmsg)

}
