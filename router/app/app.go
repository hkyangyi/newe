package app

import (
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type App struct {
	C *gin.Context
	M interface{}
}

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
	Success   bool        `json:"success"`
	Timestamp int64       `json:"timestamp"`
}

func NewApp(c *gin.Context) *App {
	var app = App{
		C: c,
	}
	return &app
}

func (a *App) Bind(v interface{}) error {
	a.M = v
	err := a.C.Bind(v)
	if err != nil {
		return err
	}
	a.C.Set("ResData", v)
	valid := validation.Validation{}
	check, err := valid.Valid(v)
	if err != nil {
		return err
	}
	if !check {
		return valid.Errors[0]
	}
	return nil
}

// 返回成功
func (a *App) SUCCESS(data interface{}) {
	a.C.JSON(200, Response{
		Code:      0,
		Message:   "ok",
		Result:    data,
		Success:   true,
		Timestamp: time.Now().Unix(),
	})
	return
}

// 返回错误
func (a *App) Error(err error) {
	a.C.JSON(200, Response{
		Code:      400,
		Message:   err.Error(),
		Result:    nil,
		Success:   false,
		Timestamp: time.Now().Unix(),
	})
	return
}

// 登陆失败
func (a *App) LoginError(err error) {
	a.C.JSON(401, Response{
		Code:      10000,
		Message:   err.Error(),
		Result:    nil,
		Success:   false,
		Timestamp: time.Now().Unix(),
	})
	return
}
