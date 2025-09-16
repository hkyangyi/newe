package worklog

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/hkyangyi/newe/common/file"

	"github.com/aiwuTech/fileLogger"
)

// 日志类型常量
const (
	LogTypeSQL   = "sql"
	LogTypeError = "error"
	LogTypeTrace = "trace"
)

// 时间格式常量
const (
	DateFormat     = "20060102"
	MonthFormat    = "200601"
	LogFileExt     = "log"
	LogFileSuffix  = "-"
	LogFileMaxSize = 10 // MB
)

// 日志配置常量
const (
	LogFileMaxBackups = 2
	LogFileMaxAge     = 300
	LogFileMaxLine    = 5000
)

// WorkLog 工作日志结构体
type WorkLog struct {
	logTime  string
	fileName string
	path     string
	mu       sync.RWMutex
	sql      *fileLogger.FileLogger
	err      *fileLogger.FileLogger
	trace    *fileLogger.FileLogger
}

// 全局日志实例
var Logio *WorkLog

// Init 初始化工作日志
func Init(path string) *WorkLog {
	if err := file.IsNotExistMkDir(path); err != nil {
		log.Printf("创建日志目录失败: %v", err)
	}

	Logio = &WorkLog{
		path:    path,
		logTime: time.Now().Format(DateFormat),
	}

	// 初始化日志文件
	Logio.rotateLogFiles()

	return Logio
}

// GetLogger 获取全局日志实例
func GetLogger() *WorkLog {
	return Logio
}

// rotateLogFiles 轮转日志文件
func (l *WorkLog) rotateLogFiles() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logTime = time.Now().Format(DateFormat)
	l.fileName = fmt.Sprintf("%s.%s", l.logTime, LogFileExt)
	monthDir := time.Now().Format(MonthFormat)

	// 创建并初始化各类型日志
	l.trace = l.createLogger(LogTypeTrace, monthDir)
	l.err = l.createLogger(LogTypeError, monthDir)
	l.sql = l.createLogger(LogTypeSQL, monthDir)
}

// createLogger 创建指定类型的日志记录器
func (l *WorkLog) createLogger(logType, monthDir string) *fileLogger.FileLogger {
	logPath := filepath.Join(l.path, logType, monthDir)
	
	if err := file.IsNotExistMkDir(logPath); err != nil {
		log.Printf("创建%s日志目录失败: %v", logType, err)
		// 尝试使用基础路径作为备选
		logPath = l.path
	}

	return fileLogger.NewSizeLogger(
		logPath,
		l.fileName,
		LogFileSuffix,
		LogFileMaxSize,
		LogFileMaxBackups,
		fileLogger.MB,
		LogFileMaxAge,
		LogFileMaxLine,
	)
}

// checkRotate 检查是否需要轮转日志文件
func (l *WorkLog) checkRotate() {
	currentDate := time.Now().Format(DateFormat)
	
	l.mu.RLock()
	needRotate := l.logTime != currentDate
	l.mu.RUnlock()
	
	if needRotate {
		l.rotateLogFiles()
	}
}

// SQL 记录SQL日志
func (l *WorkLog) SQL(v ...interface{}) {
	l.checkRotate()
	
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	l.sql.Println(v...)
}

// SQLf 记录格式化的SQL日志
func (l *WorkLog) SQLf(format string, v ...interface{}) {
	l.checkRotate()
	
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	l.sql.Printf(format, v...)
}

// Error 记录错误日志
func (l *WorkLog) Error(v ...interface{}) {
	l.checkRotate()
	
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	l.err.Println(v...)
}

// Errorf 记录格式化的错误日志
func (l *WorkLog) Errorf(format string, v ...interface{}) {
	l.checkRotate()
	
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	l.err.Printf(format, v...)
}

// Trace 记录跟踪日志
func (l *WorkLog) Trace(v ...interface{}) {
	l.checkRotate()
	
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	l.trace.Println(v...)
}

// Tracef 记录格式化的跟踪日志
func (l *WorkLog) Tracef(format string, v ...interface{}) {
	l.checkRotate()
	
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	l.trace.Printf(format, v...)
}

// 为了向后兼容，保留原有的函数名，但内部调用新的实现
func WorkLogInit(path string) *WorkLog {
	return Init(path)
}

func (l *WorkLog) NewFile() {
	l.rotateLogFiles()
}

func (l *WorkLog) WSQL(v ...interface{}) {
	l.SQL(v...)
}

func (l *WorkLog) WERR(v ...interface{}) {
	l.Error(v...)
}

func (l *WorkLog) WTRACE(v ...interface{}) {
	l.Trace(v...)
}

// 全局便捷函数
func SQL(v ...interface{}) {
	if Logio != nil {
		Logio.SQL(v...)
	}
}

func SQLf(format string, v ...interface{}) {
	if Logio != nil {
		Logio.SQLf(format, v...)
	}
}

func Error(v ...interface{}) {
	if Logio != nil {
		Logio.Error(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if Logio != nil {
		Logio.Errorf(format, v...)
	}
}

func Trace(v ...interface{}) {
	if Logio != nil {
		Logio.Trace(v...)
	}
}

func Tracef(format string, v ...interface{}) {
	if Logio != nil {
		Logio.Tracef(format, v...)
	}
}
