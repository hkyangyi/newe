package serve

import (
	"neweadmin/common/base"
	"neweadmin/common/config"
	"neweadmin/common/db"
	"neweadmin/common/ftpmag"
	"neweadmin/common/redis"
	"neweadmin/common/worklog"
	"neweadmin/common/ws"
)

type logprint struct {
}

func SetUp() {
	//读取配置文件
	base.Conf = config.GetConfig()
	//开启日志系统
	base.WorkLog = worklog.WorkLogInit(base.Conf.HTTP_RuntimeRootPath)
	//数据库连接
	base.MYDB = db.GormInit(base.Conf.DB_Host, base.Conf.DB_User, base.Conf.DB_Password, base.Conf.DB_Name)
	base.REDIS, _ = redis.NewRedis(base.Conf.REDIS_Host, base.Conf.REDIS_Password, base.Conf.REDIS_IdleTimeout, base.Conf.REDIS_MaxIdle, base.Conf.REDIS_MaxActive)
	//启动ws管理器
	base.WSMAG = ws.NewMag()
	go base.WSMAG.SetUp()
	//读取字典写入REDIS
	go base.DictInit()
	base.FtpMag = ftpmag.NewFtpNag()
	//启动WEB服务
	InitHttpServe()
}
