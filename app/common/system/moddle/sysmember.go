package moddle

import (
	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

type SysMember struct {
	ID         string `gorm:"primary_key" json:"id"`                                  //
	DepartId   string `json:"departId" form:"departId"  dict:"DepartName_sysDepart" ` //组织结构ID
	DepartName string `json:"departName" gorm:"-"`
	RoleId     string `json:"roleId" form:"roleId"  dict:"RoleName_sysRole" ` //角色ID
	RoleName   string `json:"roleName" gorm:"-"`
	UID        string `json:"uid"`                      //会员ID
	Username   string `json:"username" form:"username"` //登陆账号
	Password   string `json:"password"`                 //密码
	Nickname   string `json:"nickname" form:"nickname"` //昵称
	Realname   string `json:"realname" form:"realname"` //真实姓名
	Headimgurl string `json:"headimgurl"`               //头像
	Mp         string `json:"mp"`                       //手机号
	Idcard     string `json:"idcard"`                   //身份证号码
	Sex        int    `json:"sex"`                      //性别 1男2女
	Status     int    `json:"status"`                   //1正常，2禁用
	OrgCode    string `json:"orgCode"`                  //组织结构编码
	Remark     string `json:"remark"`                   //备注
	CreateTime int64  `json:"createTime"`               //创建时间
	UpdateTime int64  `json:"updateTime"`               //更新时间
	Files      string `json:"files"`                    //附件
	utils.PageList
}

// 添加
func (a *SysMember) Add() error {
	a.ID = utils.GetUUID()
	err := db.Db.Create(a).Error
	return err
}

// 编辑
func (a *SysMember) Edit() error {
	err := db.Db.Model(a).Updates(a).Error
	return err
}

// 获取列表
func (a *SysMember) GetList(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []SysMember
	db.Db.Model(&SysMember{}).Where(where, v...).Count(&page.Total).Order("create_time desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	page.List = items
	return page
}

// 删除
func (a *SysMember) Del() error {
	err := db.Db.Model(a).Delete(a).Error
	return err
}

// 用户名查询
func FindMemberByUsername(username string) (SysMember, error) {
	var data SysMember
	err := db.Db.Model(&data).Where("username = ?", username).First(&data).Error
	return data, err
}

func (a *SysMember) Refresh() {
	db.Db.First(a)
}

// 关键词搜索
func GetSysMemberByKeyword(keyword, orgcode string) []SysMember {
	var items []SysMember
	db.Db.Model(&SysMember{}).Where("realname like ?", "%"+keyword+"%").Where("org_code like ?", orgcode+"%").Limit(10).Find(&items)
	return items
}

// ORG搜索
func GetSysMemberByOrg(orgcode string) []SysMember {
	var items []SysMember
	db.Db.Model(&SysMember{}).Where("org_code like ?", orgcode+"%").Find(&items)
	return items
}

func (a *SysMember) EditPass(pass string) error {
	err := db.Db.Model(a).Update("password", pass).Error
	return err
}
