package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/model"
	"github.com/hkyangyi/newe/router/app"

	"github.com/gin-gonic/gin"
)

type AuthFrom struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Tags 登录
// @Produce  json
// @Param username query string true "账号" maxlength(100)
// @Param password query string true "密码" maxlength(100)
// @Success 200 {object}  "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/login [post]
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

// @Tags 获取用户登录信息
// @Produce  json
// @Success 200 {object}  "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/getuserinfo [get]
func GetUserInfo(c *gin.Context) {
	var a = app.NewApp(c)

	merdata, b := a.C.Get("AdminAuthData")
	if !b {
		a.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(model.AdminAuth)
	mer.RefreshByMerdb()
	a.SUCCESS(mer)
}

// @Tags 获取按钮权限代码
// @Produce  json
// @Success 200 {object}  "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/getuserinfo [get]
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

//测试添加菜单

type VueRouteMenu struct {
	Id        string                 `json:"id"`
	Pid       string                 `json:"pid"`
	Name      string                 `json:"name"`
	Meta      Mate                   `json:"meta"`
	Component string                 `json:"component"`
	Children  []VueRouteMenu         `json:"children"`
	FullPath  string                 `json:"fullPath"`
	Path      string                 `json:"path"`
	Redirect  string                 `json:"redirect"`
	Props     map[string]interface{} `json:"props"`
}

type Mate struct {
	OrderNo int `json:"orderNo"`
	// title
	Title string `json:"title"`
	// 动态路由器级别.
	DynamicLevel int `json:"dynamicLevel"`
	// 动态路由器真实路由路径
	RealPath string `json:"realPath"`
	// 是否忽略权限
	IgnoreAuth bool `json:"ignoreAuth"`
	// 角色
	Roles []string `json:"roles"`
	// 是否缓存
	IgnoreKeepAlive bool `json:"ignoreKeepAlive"`
	// 是否固定在标签上
	Affix bool `json:"affix"`
	// icon on tab
	Icon     string `json:"icon"`
	FrameSrc string `json:"frameSrc"`
	// current page transition
	TransitionName string `json:"transitionName"`
	// Whether the route has been dynamically added
	HideBreadcrumb bool `json:"hideBreadcrumb"`
	// 隐藏子菜单
	HideChildrenInMenu bool `json:"hideChildrenInMenu"`
	// 承载参数
	CarryParam bool `json:"carryParam"`
	// Used internally to mark single-level menus
	Single bool `json:"single"`
	// Currently active menu
	CurrentActiveMenu string `json:"currentActiveMenu"`
	// 从不显示在选项卡中
	HideTab bool `json:"hideTab"`
	// 不在菜单中显示
	HideMenu bool `json:"hideMenu"`
	IsLink   bool `json:"IsLink"`
	// 仅为菜单生成
	IgnoreRoute bool `json:"ignoreRoute"`
	// 隐藏子项的路径
	HidePathForChildren bool `json:"HidePathForChildren"`
}

func CreateMenu(c *gin.Context) {
	var data []VueRouteMenu
	var a = app.NewApp(c)
	if err := a.Bind(&data); err != nil {
		return
	}
	CreateMenuls(data, "demo2")

	a.SUCCESS(data)
}

func CreateMenuls(data []VueRouteMenu, pid string) {
	for _, v := range data {
		var typei int
		if len(v.Children) > 0 {
			typei = 1
		} else {
			typei = 2
		}

		var props string
		if len(v.Props) > 0 {
			marshal, _ := json.Marshal(v.Props)
			props = string(marshal)
		}

		var a = model.SysMenus{
			ID:         utils.GetUUID(),
			Pid:        pid,
			Name:       v.Meta.Title,
			Component:  v.Component,
			Icon:       v.Meta.Icon,
			IsExt:      utils.BoolToInt(v.Meta.IsLink),          //是否外链
			Keepalive:  utils.BoolToInt(v.Meta.IgnoreKeepAlive), //是否缓存
			Show:       utils.BoolToInt(v.Meta.HideMenu),        //是否显示
			Type:       typei,                                   //类型 1目录2菜单 3按钮
			SortNo:     v.Meta.OrderNo,                          //排序
			RoutePath:  v.Path,                                  //路由地址
			Status:     0,                                       //是否启用 0启用1禁用
			CreateTime: time.Now().Unix(),                       //
			IframeSrc:  v.Meta.FrameSrc,                         //是否嵌套
			RouteName:  v.Name,                                  //路由名称
			Redirect:   v.Redirect,                              //默认页
			Props:      props,
		}
		a.Add()

		if len(v.Children) > 0 {
			CreateMenuls(v.Children, a.ID)
		}
	}
}

type verifyform struct {
	TableName  string `form:"tablename" valid:"Required; MaxSize(50)"`
	FieldName  string `form:"fieldname" valid:"Required; MaxSize(50)"`
	Tablevalue string `form:"tablevalue" valid:"Required; MaxSize(50)"`
	TableId    string `form:"tableid"`
}

// 更新用户信息
func UpUserInfo(c *gin.Context) {
	var g = app.NewApp(c)
	var data model.SysMember
	g.Bind(&data)
	merdata, b := g.C.Get("AdminAuthData")
	if !b {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(model.AdminAuth)
	data.ID = mer.Merdb.ID
	data.Edit()

	g.SUCCESS(nil)
}

type STEditPass struct {
	OldPass string `json:"oldPass" form:"oldPass"`
	NewPass string `json:"newPass" form:"newPass"`
}

// 修改密码
func EditPass(c *gin.Context) {
	var (
		a STEditPass
		g = app.NewApp(c)
	)

	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	merdata, b := g.C.Get("AdminAuthData")
	if !b {
		g.LoginError(errors.New("登陆超时"))
		return
	}
	mer := merdata.(model.AdminAuth)
	mer.RefreshByMerdb()
	oldpass := utils.EncodeMD5(a.OldPass)
	if oldpass != mer.Merdb.Password {
		g.Error(errors.New("密码错误"))
		return
	}
	newpass := utils.EncodeMD5(a.NewPass)

	err := mer.Merdb.EditPass(newpass)
	if err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}
