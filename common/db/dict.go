package db

import (
	"fmt"

	"github.com/hkyangyi/newe/common/redis"
	"github.com/hkyangyi/newe/common/utils"
)

type SysDictList struct {
	ID         string `gorm:"primary_key" json:"id"  form:"id"` //
	ParentId   string `json:"parentId" form:"parentId"`         //
	ParentName string `json:"parentName" form:"code"`           //字典名称
	Name       string `json:"name" gorm:"column:name"`
	Value      string `json:"value" gorm:"column:value"`
	Status     int    `json:"status"`   //1正常，2停用，10已删除
	Sort       int    `json:"sort"`     //排序
	Type       int    `json:"type"`     //1自定义 2 数据表
	TableKey   string `json:"tableKey"` //数据表主键
	TableVal   string `json:"tableVal"` //数据包值
}

func DictInit() {
	var items []SysDictList
	Db.Model(&SysDictList{}).Where("parent_id = ?", "").Find(&items)
	for _, v := range items {
		rediskey := "DICT_" + v.ParentName
		var its []SysDictList
		if v.Type == 1 {
			Db.Table("sys_dict_list").Where("parent_id = ?", v.ID).Order("sort asc").Find(&its)
		} else {
			table := utils.Camel2Case(v.ParentName)
			selstr := fmt.Sprintf("%s as value, %s as name", v.TableKey, v.TableVal)
			Db.Table(table).Select(selstr).Scan(&its)
		}

		var dictdata = make(map[string]string, len(its))
		for i := 0; i < len(its); i++ {
			ks := its[i].Value
			dictdata[ks] = its[i].Name
		}

		SetDict(rediskey, dictdata)
	}

	// 自动创建表
	// err := MYDB.AutoMigrate(&BaseHikGateList{})
	// if err != nil {
	// 	panic("failed to migrate BaseHikGateList")
	// }

	// err = MYDB.AutoMigrate(&BaseHikGateAvList{})
	// if err != nil {
	// 	panic("failed to migrate BaseHikGateAvList")
	// }
	//手动将表放入字典

}

func SetDict(key string, data map[string]string) {
	redis.REDIS.SetLong(key, data)
}

// 根据字典编码获取字典
func GetDict(key string) map[string]string {
	rediskey := "DICT_" + key
	var data = make(map[string]string)
	//查询REDIS是否存在
	if b := redis.REDIS.Exists(key); !b {

		var (
			fdb SysDictList
			its []SysDictList
		)
		Db.Model(&SysDictList{}).Where("parent_name = ?", key).Find(&fdb)
		if fdb.Type == 1 {
			Db.Table("sys_dict_list").Where("parent_id = ?", fdb.ID).Order("sort asc").Find(&its)
		} else {

			table := utils.Camel2Case(fdb.ParentName)
			selstr := fmt.Sprintf("%s as value, %s as name", fdb.TableKey, fdb.TableVal)
			Db.Table(table).Select(selstr).Scan(&its)

		}
		for i := 0; i < len(its); i++ {
			ks := its[i].Value
			data[ks] = its[i].Name
		}
		redis.REDIS.SetLong(rediskey, data)
		return data
	}

	redis.REDIS.Get(rediskey, &data)
	return data
}

// 更新字典
func UpdateDict(key string) {
	rediskey := "DICT_" + key

	var data = make(map[string]string)
	var (
		fdb SysDictList
		its []SysDictList
	)
	Db.Model(&SysDictList{}).Where("parent_name = ?", key).Find(&fdb)

	if fdb.Type == 1 {
		Db.Table("sys_dict_list").Where("parent_id = ?", fdb.ID).Order("sort asc").Find(&its)
	} else {

		table := utils.Camel2Case(fdb.ParentName)
		selstr := fmt.Sprintf("%s as value, %s as name", fdb.TableKey, fdb.TableVal)
		Db.Table(table).Select(selstr).Scan(&its)

	}

	for i := 0; i < len(its); i++ {
		ks := its[i].Value
		data[ks] = its[i].Name
	}
	redis.REDIS.SetLong(rediskey, data)
}
