package newe

import (
	"github.com/hkyangyi/newe/common/config"
	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/common/worklog"
	"github.com/hkyangyi/newe/common/ws"
	"gorm.io/gorm"
)

var (
	Conf    *config.Config
	WorkLog *worklog.WorkLog
	MYDB    *gorm.DB
	REDIS   *redis.NeRedis
	WSMAG   *ws.WsManager
)

type logprint struct {
}

func InitConfig() {
	//读取配置文件
	Conf = config.ReadConfig()

	//开启日志系统
	WorkLog = worklog.WorkLogInit(Conf.HTTP_RuntimeRootPath)
	//数据库连接
	MYDB = db.GormInit(Conf.DB_Host, Conf.DB_User, Conf.DB_Password, Conf.DB_Name)
	REDIS, _ = redis.NewRedis(Conf.REDIS_Host, Conf.REDIS_Password, Conf.REDIS_IdleTimeout, Conf.REDIS_MaxIdle, Conf.REDIS_MaxActive)
	//初始化路由
	RouteInit()
	//启动WEB服务
	//InitHttpServe()
}

func Run() {
	InitHttpServe()
	//启动ws管理器
	WSMAG = ws.NewMag()
	go WSMAG.SetUp()
	//读取字典写入REDIS
	//go DictInit()
	HttpServeRun()
}
