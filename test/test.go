package test

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe"
)

// @title			NEWE API
// @version		1.0
// @description	后台系统框架
func main() {
	newe.InitConfig()
	//newe.NewRoute()
	RegRoute()

	newe.Run()
	fmt.Println(newe.Conf)
}

func RegRoute() {
	//r := newe.NewRoute("www", middle.AdminAuth())
	r := newe.NewRoute("www")
	r.Any("home", test)
}

func test(c *gin.Context) {
	path := c.Request.URL.Path
	fmt.Println(path)

	c.String(200, "aaaa")
}
