package base

import (
	"fmt"

	"github.com/hkyangyi/newe/common/config"
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
