package serve

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	roootRouter "github.com/hkyangyi/newe/app/router"
	"github.com/hkyangyi/newe/common/base"

	"github.com/gin-gonic/gin"
)

func InitHttpServe() {
	gin.SetMode(base.Conf.HTTP_RunMode)
	routersInit := InitRouter()
	readTimeout := time.Duration(int64(base.Conf.HTTP_ReadTimeout)) * time.Second
	writeTimeout := time.Duration(int64(base.Conf.HTTP_WriteTimeout)) * time.Second
	endPoint := fmt.Sprintf(":%d", base.Conf.HTTP_Port)
	maxHeaderBytes := 1 << 20
	Server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	fmt.Printf("[info] start http server listening %s", endPoint)
	Server.ListenAndServe()
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
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
		base.WorkLog.WTRACE(str)
		return str
	}))
	r.Use(gin.Recovery())
	r.StaticFS("/upload/images", http.Dir(base.Conf.IMG_SavePath))  //文件目录
	r.StaticFS("/upload/fields", http.Dir(base.Conf.FILE_SavePath)) //文件目录
	r.Use(Cors())
	roootRouter.RegRouter(r)
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
