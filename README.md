go get https://github.com/hkyangyi/newe


func main() {
	newe.InitConfig()
	//newe.NewRoute()
	RegRoute()

	newe.Run()
	fmt.Println(newe.Conf)
}

func RegRoute() {
	r := newe.NewRoute("www", middle.AdminAuth())
	r.Any("home", test)
}

func test(c *gin.Context) {
	c.String(200, "aaaa")
}
