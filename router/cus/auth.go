package cus

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/cusmodel"
	"github.com/hkyangyi/newe/router/app"
)

type AuthFrom struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var data AuthFrom
	var a = app.NewApp(c)
	if err := a.Bind(&data); err != nil {
		a.Error(err)
		return
	}
	res, err := data.MerLogin()
	if err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(res)
}

func (a *AuthFrom) MerLogin() (cusmodel.CusAuth, error) {
	var data cusmodel.CusAuth
	data.Isadmin = false
	if len(a.Username) < 3 {
		return data, errors.New("请输入正确的账号")
	}

	rediskey := "CUS_AUTH_ERROR_COUNT_" + a.Username
	var LoginErrCount int
	if redis.REDIS.Exists(rediskey) {
		redis.REDIS.Get(rediskey, &LoginErrCount)
	} else {
		LoginErrCount = 0
	}
	if LoginErrCount >= 5 {
		return data, errors.New("您今日输入错误次数过多，请24小时后再试")
	}

	//查询用户信息
	merdb, err := cusmodel.FindMemberByUsername(a.Username, utils.EncodeMD5(a.Password))
	if err != nil {
		cdb, err := cusmodel.FindCusUserByUsername(a.Username, utils.EncodeMD5(a.Password))
		if err != nil {
			LoginErrCount++
			redis.REDIS.Set(rediskey, LoginErrCount, 60*60*24)
			return data, errors.New("账号或密码错误,您还有" + strconv.Itoa(5-LoginErrCount) + "次机会")
		}
		data.Isadmin = true
		data.CDB = cdb
	}

	//登陆成功
	data.MerDb = merdb
	if err := data.RefreshByMerdb(); err != nil {
		return data, err
	}
	var authkey string

	if data.Isadmin {
		authkey = fmt.Sprintf("CUS_AUTH_%s_%s", data.CDB.Id, data.CDB.Username)
	} else {
		authkey = fmt.Sprintf("CUS_AUTH_%s_%s", data.MerDb.ID, data.MerDb.Username)
	}

	token, _ := utils.SetToken(authkey)
	data.Token = token

	redis.REDIS.Set(authkey, data, 60*60)

	return data, nil
}

//

/**
 * 刷新accessToken
 */

func RefreshToken(c *gin.Context) {
	var a = app.NewApp(c)
	merdata, b := a.C.Get("CusAdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)

	if err := mer.RefreshByMerdb(); err != nil {
		a.LoginError(err)
		return
	}
	var authkey string
	if mer.Isadmin {
		authkey = fmt.Sprintf("CUS_AUTH_%s_%s", mer.CDB.Id, mer.CDB.Username)
	} else {
		authkey = fmt.Sprintf("CUS_AUTH_%s_%s", mer.MerDb.ID, mer.MerDb.Username)
	}
	token, _ := utils.SetToken(authkey)
	mer.Token = token
	redis.REDIS.Set(authkey, mer, 60*60)
	a.SUCCESS(mer)
}

/**
 * 退出登录
 */

func LoginOut(c *gin.Context) {
	var g = app.NewApp(c)
	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	uuid, err := utils.AuthToken(token)
	if err != nil {
		g.LoginError(nil)
		return
	}
	res := redis.REDIS.Exists(uuid)
	if !res {
		g.LoginError(nil)
		return
	}
	redis.REDIS.Delete(uuid)
	g.SUCCESS(nil)
}

/**
 * 获取用户权限码
 */

func GetPermCode(c *gin.Context) {
	var a = app.NewApp(c)
	merdata, b := a.C.Get("CusAdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(cusmodel.CusAuth)

	var items []string
	if mer.Isadmin {
		items = cusmodel.CustomerButtonGetList()
	} else {
		items = cusmodel.CustomerButtonGetByDepart(mer.MerDb.DepartId)
	}
	a.SUCCESS(items)
}
