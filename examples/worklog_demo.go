package main

import (
	"fmt"
	"time"

	"github.com/hkyangyi/newe/common/worklog"
)

func main() {
	// 初始化日志系统，指定日志存储路径
	logPath := "./logs"
	logger := worklog.Init(logPath)
	fmt.Println("日志系统初始化完成，日志存储路径:", logPath)

	// 示例1：使用实例方法记录日志
	logger.SQL("执行SQL查询:", "SELECT * FROM users WHERE id = 1")
	logger.Error("发生错误:", "连接数据库失败")
	logger.Trace("跟踪信息:", "用户登录", "用户ID: 12345")

	// 示例2：使用格式化方法记录日志
	logger.SQLf("执行查询: %s，参数: %v", "SELECT * FROM products WHERE category = ?", "electronics")
	logger.Errorf("错误码: %d, 错误信息: %s", 500, "服务器内部错误")
	logger.Tracef("用户 %s (ID: %d) 执行了 %s 操作", "张三", 10001, "购买商品")

	// 示例3：使用全局便捷函数记录日志
	worklog.SQL("全局SQL日志:", "INSERT INTO logs VALUES ('test')")
	worklog.Error("全局错误日志:", "文件不存在")
	worklog.Trace("全局跟踪日志:", "API调用", "/api/v1/users", "GET")

	// 示例4：使用全局格式化便捷函数
	worklog.SQLf("执行更新操作: %s", "UPDATE users SET status = 'active'")
	worklog.Errorf("权限错误: 用户 %s 没有 %s 权限", "李四", "删除")
	worklog.Tracef("处理请求耗时: %.2f ms", 123.45)

	// 示例5：兼容旧API
	logger.WSQL("使用旧API记录SQL:", "SELECT version()")
	logger.WERR("使用旧API记录错误:", "配置文件解析失败")
	logger.WTRACE("使用旧API记录跟踪:", "系统启动完成")

	// 模拟日志轮转（实际应用中不需要这样做，系统会自动检测日期变化）
	fmt.Println("模拟日志轮转...")
	logger.NewFile()
	
	// 继续记录日志
	worklog.Trace("日志轮转后的记录")
	
	fmt.Println("日志示例执行完成，请查看", logPath, "目录下的日志文件")
	
	// 等待一秒，确保日志写入完成
	time.Sleep(time.Second)
}