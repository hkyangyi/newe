package moddle

import (
	"errors"
	"github.com/hkyangyi/newe/common/base"
	"github.com/hkyangyi/newe/common/utils"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SysDepart struct {
	ID         string      `gorm:"primary_key" json:"id" form:"id"` //uuid
	Pid        string      `json:"pid" form:"pid"`                  //父级ID
	Name       string      `json:"name" form:"name"`                //分组名称（机构名称）
	Code       string      `json:"code"`                            //分组编码
	Type       int         `json:"type" dict:"TypeStr_DepartType"`  //类型（1集团，2公司，3部门，4服务门店）
	TypeStr    string      `json:"typeStr" gorm:"-"`
	Telephone  string      `json:"telephone"`  //联系电话
	Phone      string      `json:"phone"`      //联系手机
	Address    string      `json:"address"`    //地址
	SortNo     int         `json:"sortNo"`     //排序
	CreateTime int64       `json:"createTime"` //创建时间
	UpdateTime int64       `json:"updateTime"` //更新时间
	Disabled   bool        `gorm:"-" json:"disabled"`
	List       []SysDepart `gorm:"-" json:"children"`
}

// 添加
func (a *SysDepart) Add() error {
	a.ID = utils.GetUUID()
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = time.Now().Unix()
	if len(a.Pid) == 0 {
		a.Code = a.GetCode()
	} else {
		var fdb SysDepart
		base.MYDB.Model(&SysDepart{}).Where("id = ?", a.Pid).First(&fdb)
		code := a.GetCode()
		a.Code = fdb.Code + "-" + code
	}
	err := base.MYDB.Create(a).Error
	return err
}

// 编辑
func (a *SysDepart) Edit() error {
	if len(a.ID) != 32 {
		return errors.New("缺少参数ID")
	}
	err := base.MYDB.Transaction(func(tx *gorm.DB) error {
		//查询原始数据
		var odb = SysDepart{
			ID: a.ID,
		}
		tx.First(&odb)
		//查询是否修改父级
		if odb.Pid != a.Pid {
			var fdb = SysDepart{
				ID: a.Pid,
			}
			tx.First(&fdb)
			code := GetCodeByDepartId(a.ID)
			a.Code = fdb.Code + "-" + code
			var items []SysDepart
			base.MYDB.Model(&SysDepart{}).Where("code like ?", odb.Code+"%").Or("org_code like ?", odb.Code+"%").Find(&items)
			for _, v := range items {
				v.Code = strings.Replace(v.Code, odb.Code, a.Code, -1)
				if err := tx.Model(&v).Updates(&v).Error; err != nil {
					return err
				}
				if err := tx.Model(&SysMember{}).Where("org_code = ?", odb.Code).Update("org_code", v.Code).Error; err != nil {
					return err
				}
			}
		}
		if err := tx.Model(a).Updates(a).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// 删除
func (a *SysDepart) Del() error {
	var count int64
	base.MYDB.Table("sys_depart").Where("pid = ?", a.ID).Count(&count)
	if count > 0 {
		return errors.New("请先删除子分类")
	}
	base.MYDB.Table("sys_member").Where("depart_id = ?", a.ID).Count(&count)
	if count > 0 {
		return errors.New("请先删除该分组下的用户")
	}
	//查找是否有
	err := base.MYDB.Table("sys_depart").Where("id = ? or pid = ?", a.ID, a.ID).Delete(&SysDepart{}).Error
	return err
}

// 获取列表
func (a *SysDepart) GetList(usdata SysMember) []SysDepart {
	var items []SysDepart
	var list []SysDepart
	if usdata.Username == "admin" {
		base.MYDB.Table("sys_depart").Order("sort_no asc").Find(&items)
		list = SysDepartDigui(items, "", list)
	} else {
		var pdb SysDepart
		base.MYDB.Table("sys_depart").Where("code", usdata.OrgCode).Order("sort_no asc").First(&pdb)
		base.MYDB.Table("sys_depart").Where("code like ?", pdb.Code+"%").Order("sort_no asc").Find(&items)
		//items = append(items, pdb)
		list = SysDepartDigui(items, pdb.Pid, list)
	}
	return list
}

func SysDepartDigui(items []SysDepart, pid string, list []SysDepart) []SysDepart {
	var item []SysDepart
	for _, v := range items {
		if v.Pid == pid {
			var ls SysDepart

			utils.StAtoB(v, ls, &ls)
			ls.List = SysDepartDigui(items, v.ID, item)
			list = append(list, ls)
		}

	}
	return list
}

func SysDepartGetBYId(id string) (SysDepart, error) {
	var data = SysDepart{
		ID: id,
	}
	err := base.MYDB.Model(&data).First(&data).Error
	return data, err
}

func SysDepartGetByCode(code string) (SysDepart, error) {
	var data SysDepart
	err := base.MYDB.Model(&data).Where("code = ?", code).First(&data).Error
	return data, err
}

type SysDepartDict struct {
	ID        string `gorm:"primary_key" json:"id"` //
	DepartId  string `json:"departId"`              //结构ID
	Lv        int    `json:"lv"`                    //等级
	Code      int64  `json:"code"`                  //编号
	ServeCode string `json:"serveCode"`
}

func (a *SysDepart) GetCode() string {
	var data = SysDepartDict{
		ID:        utils.GetUUID(),
		DepartId:  a.ID,
		ServeCode: base.Conf.HTTP_ServeCode,
	}
	var count int64
	base.MYDB.Model(&SysDepartDict{}).Where("serve_code = ?", data.ServeCode).Count(&count)
	data.Code = count
	base.MYDB.Create(&data)
	return data.ServeCode + strconv.Itoa(int(data.Code))

}

func GetCodeByDepartId(departId string) string {
	var data SysDepartDict
	base.MYDB.Model(&SysDepartDict{}).Where("depart_id = ?", departId).First(&data)
	return data.ServeCode + strconv.Itoa(int(data.Code))
}

type SysDepartRules struct {
	ID         string `gorm:"primary_key" json:"id"` //
	DepartId   string `json:"departId"`              //组织结构ID
	OrgCode    string `json:"orgCode"`               //组织结构编码
	MenuId     string `json:"menuId"`                //菜单ID
	CreateTime int64  `json:"createTime"`            //
}

// 根据部门获取已有权限
func (a *SysDepart) GetRules() []SysDepartRules {
	var items []SysDepartRules
	base.MYDB.Model(&items).Where("depart_id = ?", a.ID).Find(&items)
	return items
}

func (a *SysDepart) DelRules(ids []string) error {
	err := base.MYDB.Table("sys_depart_rules").Where("menu_id IN ?", ids).Delete(&SysDepartRules{}).Error
	return err
}

func (a *SysDepart) AddRules(ids []string) error {
	var items []SysDepartRules
	for _, v := range ids {
		item := SysDepartRules{
			ID:         utils.GetUUID(),
			OrgCode:    a.Code,
			DepartId:   a.ID,
			MenuId:     v,
			CreateTime: time.Now().Unix(),
		}
		items = append(items, item)
	}
	err := base.MYDB.Create(&items).Error
	return err
}

func DepartRulesByDeaprtId(id string) []SysDepartRules {
	var items []SysDepartRules
	base.MYDB.Table("sys_depart_rules").Where("depart_id  =  ?", id).Find(&items)
	return items
}
