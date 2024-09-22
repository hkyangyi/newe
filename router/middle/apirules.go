package middle

import (
	"fmt"

	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/model"
)

// 更具PATH获取API数据
func GetApiData(path string) model.SysApiList {
	rediskey := "APIDATA_" + path
	var data model.SysApiList
	if res := redis.REDIS.Exists(rediskey); !res {
		goto SqlFind
	}

	if errs := redis.REDIS.Get(rediskey, &data); errs != nil {
		goto SqlFind
	}

	return data

SqlFind:
	data.Path = path
	data.GetByPath()
	redis.REDIS.Set(rediskey, data, 60)
	return data

}

// 更具角色获取该api的权限
func GetRoleRules(apiid, roleid string) model.SysRoleRules {
	rediskey := fmt.Sprintf("ROLEAPIRULES_%s_%s", apiid, roleid)
	var data model.SysRoleRules
	if res := redis.REDIS.Exists(rediskey); !res {
		goto SqlFind
	}

	if errs := redis.REDIS.Get(rediskey, &data); errs != nil {
		goto SqlFind
	}

	return data

SqlFind:
	data.ApiId = apiid
	data.RoleId = roleid
	data.GetByApiVsRole()
	redis.REDIS.Set(rediskey, data, 60)
	return data
}
