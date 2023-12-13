package base

import (
	"fmt"
	"newe/common/config"
	"newe/common/ftpmag"
	"newe/common/redis"
	"newe/common/worklog"
	"newe/common/ws"

	"gorm.io/gorm"
)

var (
	Conf    *config.Config
	WorkLog *worklog.WorkLog
	MYDB    *gorm.DB
	REDIS   *redis.NeRedis
	FtpMag  *ftpmag.FtpMag
	WSMAG   *ws.WsManager
)

// 检测唯一
func VerifyOnly(tablename string, id string, where map[string]interface{}) bool {
	var count int64

	if len(id) > 0 {
		MYDB.Table(tablename).Where(where).Where("id <> ?", id).Count(&count)

	} else {
		MYDB.Table(tablename).Where(where).Count(&count)
	}

	fmt.Println("条数：i%", count)
	if count > 0 {
		return true
	}
	return false
}