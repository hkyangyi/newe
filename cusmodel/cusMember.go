package cusmodel

import (
	"errors"
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

type CustomerMember struct {
	ID         string `gorm:"primary_key" json:"id"`                                  //
	DepartId   string `json:"departId" form:"departId"  dict:"DepartName_sysDepart" ` //组织结构ID
	Cid        string `json:"cid" form:"cid"`                                         //客户ID
	DepartName string `json:"departName" gorm:"-"`
	RoleId     string `json:"roleId" form:"roleId"  dict:"RoleName_sysDepart" ` //角色ID
	RoleName   string `json:"roleName" gorm:"-"`
	UID        string `json:"uid"`                      //会员ID
	Username   string `json:"username" form:"username"` //登陆账号
	Password   string `json:"password"`                 //密码
	Nickname   string `json:"nickname" form:"nickname"` //昵称
	RealName   string `json:"realName" form:"realName"` //真实姓名
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
	IsAdmin    int    `json:"isAdmin"`                  //是否管理员 1是0否
	utils.PageList
}

// 添加
func (a *CustomerMember) Add() error {
	a.ID = utils.GetUUID()
	if a.Cid == "" {
		return errors.New("缺少参数Cid")
	}
	err := db.Db.Create(a).Error
	return err
}

// 编辑
func (a *CustomerMember) Edit() error {

	err := db.Db.Model(a).Where("cid = ?", a.Cid).Updates(a).Error
	return err
}

// 获取列表
func (a *CustomerMember) GetList(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []CustomerMember
	db.Db.Model(&CustomerMember{}).Where("cid = ?", a.Cid).Where(where, v...).Count(&page.Total).Order("create_time desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	for i := 0; i < len(items); i++ {
		items[i].Password = ""
	}
	page.List = items
	return page
}

// 删除
func (a *CustomerMember) Del() error {
	a.Refresh()
	if a.IsAdmin == 1 {
		return errors.New("管理员账号无法删除")
	}
	err := db.Db.Model(a).Where("cid = ? and id = ? ", a.Cid, a.ID).Delete(a).Error
	return err
}

// 用户名查询
func FindMemberByUsername(username, password string) (CustomerMember, error) {
	var data CustomerMember
	err := db.Db.Model(&CustomerMember{}).Where("username = ? and password = ?", username, password).First(&data).Error
	return data, err
}

func (a *CustomerMember) Refresh() {
	db.Db.First(a)
}

// 关键词搜索
func GetCustomerMemberByKeyword(keyword, orgcode string) []CustomerMember {
	var items []CustomerMember
	db.Db.Model(&CustomerMember{}).Where("realname like ?", "%"+keyword+"%").Where("org_code like ?", orgcode+"%").Limit(10).Find(&items)
	return items
}

// ORG搜索
func GetCustomerMemberByOrg(orgcode string) []CustomerMember {
	var items []CustomerMember
	db.Db.Model(&CustomerMember{}).Where("org_code like ?", orgcode+"%").Find(&items)
	return items
}

func (a *CustomerMember) EditPass(pass string) error {
	err := db.Db.Model(a).Update("password", pass).Error
	return err
}

// MemberChangePassword 修改用户密码（校验旧密码）
func MemberChangePassword(memberID, oldPassword, newPassword string) error {
	if memberID == "" || oldPassword == "" || newPassword == "" {
		return errors.New("参数不完整")
	}
	var m CustomerMember
	if err := db.Db.Where("id = ? ", memberID).First(&m).Error; err != nil {
		return err
	}

	newpass := utils.EncodeMD5(newPassword)
	if newpass == m.Password {
		return errors.New("新密码不能与原密码相同")
	}

	err := db.Db.Model(&m).Update("password", newpass).Error
	return err
}

// 用户修改信息
func (a *CustomerMember) EditProfile() error {
	if a.ID == "" {
		return errors.New("缺少参数ID")
	}

	err := db.Db.Model(a).Where("id = ? and cid = ?", a.ID, a.Cid).Updates(map[string]interface{}{
		"nickname":    a.Nickname,
		"real_name":   a.RealName,
		"sex":         a.Sex,
		"mp":          a.Mp,
		"update_time": time.Now().Unix(),
	}).Error

	return err
}
