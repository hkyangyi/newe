package newe

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	roootRouter "github.com/hkyangyi/newe/router/route"
)

var Route *gin.Engine

func RouteInit() {
	Route = gin.New()
	Route.Use(gin.Logger())
	Route.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 你的自定义格式
		str := fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
		WorkLog.WTRACE(str)
		return str
	}))
	Route.Use(gin.Recovery())
	Route.StaticFS(Conf.IMG_PrefixUrl, http.Dir(Conf.IMG_SavePath))   //文件目录
	Route.StaticFS(Conf.FILE_PrefixUrl, http.Dir(Conf.FILE_SavePath)) //文件目录
	Route.Use(Cors())
	//初始化系统路由
	r := NewRoute("newesys")
	fmt.Println(r)
	roootRouter.RegRouter(r)
}

func NewRoute(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	r := Route.Group(path, handlers...)
	return r
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {

		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		headerKeys = append(headerKeys, "Token")
		headerKeys = append(headerKeys, "AccSn")
		headerKeys = append(headerKeys, "Content-Type")
		headerKeys = append(headerKeys, "Authorization")
		headerKeys = append(headerKeys, "ignoreCancelToken")
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin,x-requested-with, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {

			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", headerStr)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Set("content-type", "application/json")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
