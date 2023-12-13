package worklog

import (
	"fmt"
	"neweadmin/common/file"
	"sync"
	"time"

	"github.com/aiwuTech/fileLogger"
)

type WorkLog struct {
	LOGTIME  string
	FILENAME string
	Path     string
	WG       *sync.WaitGroup
	SQL      *fileLogger.FileLogger
	ERR      *fileLogger.FileLogger
	TRACE    *fileLogger.FileLogger
}

func WorkLogInit(path string) *WorkLog {
	file.IsNotExistMkDir(path)
	var wl = WorkLog{
		Path:    path,
		LOGTIME: time.Now().Format("20060102"),
		WG:      new(sync.WaitGroup),
	}
	wl.WG.Add(1)
	wl.NewFile()

	return &wl
}

func (a *WorkLog) NewFile() {
	a.LOGTIME = time.Now().Format("20060102")
	FILENAME := fmt.Sprintf("%s.%s", time.Now().Format("20060102"), "log")
	tpath := time.Now().Format("200601")
	tracepath := fmt.Sprintf("%s/trace/%s", a.Path, tpath)
	if err := file.IsNotExistMkDir(tracepath); err != nil {
		fmt.Println("日志创建失败")
	}
	a.TRACE = fileLogger.NewSizeLogger(tracepath, FILENAME, "-", 10, 2, fileLogger.MB, 300, 5000)

	errpath := fmt.Sprintf("%s/error/%s", a.Path, tpath)
	if err := file.IsNotExistMkDir(errpath); err != nil {
		fmt.Println("日志创建失败")
	}
	a.ERR = fileLogger.NewSizeLogger(errpath, FILENAME, "-", 10, 2, fileLogger.MB, 300, 5000)

	sqlpath := fmt.Sprintf("%s/sql/%s", a.Path, tpath)
	if err := file.IsNotExistMkDir(sqlpath); err != nil {
		fmt.Println("日志创建失败")
	}
	a.SQL = fileLogger.NewSizeLogger(sqlpath, FILENAME, "-", 10, 2, fileLogger.MB, 300, 5000)
	a.WG.Done()
}

func (a *WorkLog) WSQL(v ...interface{}) {
	a.WG.Wait()
	if a.LOGTIME != time.Now().Format("20060102") {
		a.WG.Add(1)
		a.NewFile()
	}
	a.SQL.Println(v...)
	return
}

func (a *WorkLog) WERR(v ...interface{}) {
	a.WG.Wait()
	if a.LOGTIME != time.Now().Format("20060102") {
		a.WG.Add(1)
		a.NewFile()
	}
	a.ERR.Println(v...)
	return
}

func (a *WorkLog) WTRACE(v ...interface{}) {
	a.WG.Wait()
	if a.LOGTIME != time.Now().Format("20060102") {
		a.WG.Add(1)
		a.NewFile()
	}
	a.TRACE.Println(v...)
	return
}
