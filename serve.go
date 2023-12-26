package newe

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var Serve *http.Server

func InitHttpServe() {
	Serve = &http.Server{}
	gin.SetMode(Conf.HTTP_RunMode)
	readTimeout := time.Duration(int64(Conf.HTTP_ReadTimeout)) * time.Second
	writeTimeout := time.Duration(int64(Conf.HTTP_WriteTimeout)) * time.Second
	endPoint := fmt.Sprintf(":%d", Conf.HTTP_Port)
	maxHeaderBytes := 1 << 20
	Serve.Addr = endPoint
	Serve.Handler = Route
	Serve.ReadTimeout = readTimeout
	Serve.WriteTimeout = writeTimeout
	Serve.MaxHeaderBytes = maxHeaderBytes
}

func HttpServeRun() {
	err := Serve.ListenAndServe()
	WorkLog.WERR(err)
}

func HttpServeStop() {

}
