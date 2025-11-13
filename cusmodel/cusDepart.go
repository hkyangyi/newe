package cusmodel

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
	"gorm.io/gorm"
)

type CustomerDepart struct {
	ID         string           `gorm:"primary_key" json:"id" form:"id"` //uuid
	Cid        string           `json:"cid" form:"cid"`                  //客户ID
	Pid        string           `json:"pid" form:"pid"`                  //父级ID
	Name       string           `json:"name" form:"name"`                //分组名称（机构名称）
	Code       string           `json:"code"`                            //分组编码
	Type       string           `json:"type" dict:"TypeStr_DepartType"`  //类型（1集团，2公司，3部门，4服务门店）
	TypeStr    string           `json:"typeStr" gorm:"-"`
	Telephone  string           `json:"telephone"`  //联系电话
	Phone      string           `json:"phone"`      //联系手机
	Address    string           `json:"address"`    //地址
	SortNo     int              `json:"sortNo"`     //排序
	CreateTime int64            `json:"createTime"` //创建时间
	UpdateTime int64            `json:"updateTime"` //更新时间
	Disabled   bool             `gorm:"-" json:"disabled"`
	List       []CustomerDepart `gorm:"-" json:"children"`
}

// 刷新
func (a *CustomerDepart) Refresh() error {
	err := db.Db.Model(a).First(a).Error
	return err
}

// 添加
func (a *CustomerDepart) Add() error {
	a.ID = utils.GetUUID()
	a.CreateTime = time.Now().Unix()
	a.UpdateTime = time.Now().Unix()
	if len(a.Pid) == 0 {
		a.Code = a.GetCode()
	} else {
		var fdb CustomerDepart
		db.Db.Model(&CustomerDepart{}).Where("id = ?", a.Pid).First(&fdb)
		code := a.GetCode()
		a.Code = fdb.Code + "-" + code
	}
	err := db.Db.Create(a).Error
	return err
}

// 编辑
func (a *CustomerDepart) Edit() error {
	if len(a.ID) != 32 {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Transaction(func(tx *gorm.DB) error {
		//查询原始数据
		var odb = CustomerDepart{
			ID: a.ID,
		}
		tx.First(&odb)
		//查询是否修改父级
		if odb.Pid != a.Pid {
			var fdb = CustomerDepart{
				ID: a.Pid,
			}
			tx.Where("cid = ?", a.Cid).First(&fdb)
			code := GetCodeByDepartId(a.ID)
			a.Code = fdb.Code + "-" + code
			var items []CustomerDepart
			db.Db.Model(&CustomerDepart{}).Where("code like ?", odb.Code+"%").Or("org_code like ?", odb.Code+"%").Find(&items)
			for _, v := range items {
				v.Code = strings.Replace(v.Code, odb.Code, a.Code, -1)
				if err := tx.Model(&v).Updates(&v).Error; err != nil {
					return err
				}
				if err := tx.Model(&CustomerMember{}).Where("org_code = ?", odb.Code).Update("org_code", v.Code).Error; err != nil {
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
func (a *CustomerDepart) Del() error {
	var count int64
	db.Db.Table("customer_depart").Where("pid = ?", a.ID).Count(&count)
	if count > 0 {
		return errors.New("请先删除子分类")
	}
	db.Db.Table("customer_member").Where("depart_id = ?", a.ID).Count(&count)
	if count > 0 {
		return errors.New("请先删除该分组下的用户")
	}
	//查找是否有
	err := db.Db.Table("customer_depart").Where("id = ? or pid = ?", a.ID, a.ID).Delete(&CustomerDepart{}).Error
	return err
}

// 获取列表
func (a *CustomerDepart) GetListByAdmin() []CustomerDepart {
	var items []CustomerDepart
	var list []CustomerDepart
	db.Db.Table("customer_depart").Where("cid = ?", a.Cid).Order("sort_no asc").Find(&items)
	list = CustomerDepartDigui(items, "", list)
	return list
}

// 获取列表
func (a *CustomerDepart) GetListByMember(usdata CustomerMember) []CustomerDepart {
	var items []CustomerDepart
	var list []CustomerDepart
	var pdb CustomerDepart
	db.Db.Table("customer_depart").Where("cid = ?", a.Cid).Where("id = ?", usdata.DepartId).Order("sort_no asc").First(&pdb)
	db.Db.Table("customer_depart").Where("cid = ?", a.Cid).Where("code like ?", pdb.Code+"%").Order("sort_no asc").Find(&items)
	//items = append(items, pdb)
	list = CustomerDepartDigui(items, pdb.Pid, list)
	return list
}

// 获取不带角色的
func (a *CustomerDepart) GetDepartTree(usdata CustomerMember) []CustomerDepart {
	var items []CustomerDepart
	var list []CustomerDepart
	if usdata.Username == "admin" {
		db.Db.Table("customer_depart").Where("type <> ?", "R").Order("sort_no asc").Find(&items)
		list = CustomerDepartDigui(items, "", list)
	} else {
		var pdb CustomerDepart
		db.Db.Table("customer_depart").Where("code", usdata.OrgCode).Order("sort_no asc").First(&pdb)
		db.Db.Table("customer_depart").Where("code like ? and type <> ?", pdb.Code+"%", "R").Order("sort_no asc").Find(&items)
		//items = append(items, pdb)
		list = CustomerDepartDigui(items, pdb.Pid, list)
	}
	return list
}

func (a *CustomerDepart) GetRoleById() []CustomerDepart {
	var items []CustomerDepart
	db.Db.Table("customer_depart").Where("pid = ? and type = ?", a.ID, "R").Order("sort_no asc").Find(&items)
	return items
}

func CustomerDepartDigui(items []CustomerDepart, pid string, list []CustomerDepart) []CustomerDepart {
	var item []CustomerDepart
	for _, v := range items {
		if v.Pid == pid {
			var ls CustomerDepart

			utils.StAtoB(v, ls, &ls)
			ls.List = CustomerDepartDigui(items, v.ID, item)
			list = append(list, ls)
		}

	}
	return list
}

func CustomerDepartGetBYId(id string) (CustomerDepart, error) {
	var data = CustomerDepart{
		ID: id,
	}
	err := db.Db.Model(&data).First(&data).Error
	return data, err
}

func CustomerDepartGetByCode(code string) (CustomerDepart, error) {
	var data CustomerDepart
	err := db.Db.Model(&data).Where("code = ?", code).First(&data).Error
	return data, err
}

type CustomerDepartDict struct {
	ID        string `gorm:"primary_key" json:"id"` //
	DepartId  string `json:"departId"`              //结构ID
	Lv        int    `json:"lv"`                    //等级
	Code      int64  `json:"code"`                  //编号
	ServeCode string `json:"serveCode"`
}

func (a *CustomerDepart) GetCode() string {
	var data = CustomerDepartDict{
		ID:        utils.GetUUID(),
		DepartId:  a.ID,
		ServeCode: a.Type,
	}
	var count int64
	db.Db.Model(&CustomerDepartDict{}).Where("serve_code = ?", data.ServeCode).Count(&count)
	data.Code = count
	db.Db.Create(&data)
	return data.ServeCode + strconv.Itoa(int(data.Code))

}

func GetCodeByDepartId(departId string) string {
	var data CustomerDepartDict
	db.Db.Model(&CustomerDepartDict{}).Where("depart_id = ?", departId).First(&data)
	return data.ServeCode + strconv.Itoa(int(data.Code))
}

type CustomerDepartRules struct {
	ID         string `gorm:"primary_key" json:"id"` //
	DepartId   string `json:"departId"`              //组织结构ID
	OrgCode    string `json:"orgCode"`               //组织结构编码
	MenuId     string `json:"menuId"`                //菜单ID
	CreateTime int64  `json:"createTime"`            //
}

// 根据部门获取已有权限
func (a *CustomerDepart) GetRules() []CustomerDepartRules {
	var items []CustomerDepartRules
	db.Db.Model(&items).Where("depart_id = ?", a.ID).Find(&items)
	return items
}

func (a *CustomerDepart) DelRules(ids []string) error {
	err := db.Db.Table("customer_depart_rules").Where("menu_id IN ?", ids).Delete(&CustomerDepartRules{}).Error
	return err
}

func (a *CustomerDepart) AddRules(ids []string) error {
	var items []CustomerDepartRules
	for _, v := range ids {
		item := CustomerDepartRules{
			ID:         utils.GetUUID(),
			OrgCode:    a.Code,
			DepartId:   a.ID,
			MenuId:     v,
			CreateTime: time.Now().Unix(),
		}
		items = append(items, item)
	}
	err := db.Db.Create(&items).Error
	return err
}

func DepartRulesByDeaprtId(id string) []CustomerDepartRules {
	var items []CustomerDepartRules
	db.Db.Table("customer_depart_rules").Where("depart_id  =  ?", id).Find(&items)
	return items
}
