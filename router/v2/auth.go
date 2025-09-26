package v2

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/model"
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
	res, err := data.Login()
	if err != nil {
		a.Error(err)
		return
	}
	a.SUCCESS(res)
}

func (a *AuthFrom) Login() (model.AdminAuth, error) {
	var data model.AdminAuth
	if len(a.Username) < 3 {
		return data, errors.New("请输入正确的账号")
	}

	rediskey := "AUTH_ERROR_COUNT_" + a.Username
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
	merdb, err := model.FindMemberByUsername(a.Username)
	if err != nil {
		return data, err
	}

	if len(merdb.ID) != 32 {
		return data, errors.New("账号或密码错误")
	}
	if merdb.Status != 1 {
		return data, errors.New("您的账号已停用")
	}

	md5ps := utils.EncodeMD5(a.Password)
	if md5ps != merdb.Password {
		LoginErrCount++
		redis.REDIS.Set(rediskey, LoginErrCount, 60*60*24)
		return data, errors.New("密码错误,您还有" + strconv.Itoa(5-LoginErrCount) + "次机会")
	}

	merdb.Password = ""
	//登陆成功
	data.Merdb = merdb
	err = data.RefreshByMerdb()
	if err != nil {
		return data, err
	}
	authkey := fmt.Sprintf("AUTH_%s_%s", merdb.ID, merdb.Username)
	token, _ := utils.SetToken(authkey)
	data.Token = token

	redis.REDIS.Set(authkey, data, 60*60)

	return data, nil
}

/**
 * 刷新accessToken
 */

func RefreshToken(c *gin.Context) {
	var a = app.NewApp(c)

	a.SUCCESS(nil)
}

/**
 * 退出登录
 */

func LoginOut(c *gin.Context) {
	var g = app.NewApp(c)
	token := c.GetHeader("Authorization")
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
	merdata, b := a.C.Get("AdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(model.AdminAuth)

	var items []string
	if mer.Merdb.Username == "admin" {
		items = model.SysButtonGetList()
	} else {
		items = model.SysButtonGetByDepart(mer.Merdb.DepartId)
	}
	a.SUCCESS(items)
}
