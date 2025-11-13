package model

import (
	"errors"
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
	"gorm.io/gorm"
)

func (m *DevDbtable) TableName() string {
	return "dev_dbtable"
}

type DevDbtable struct {
	ID          string `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`                     // 主键
	DbTableName string `gorm:"column:table_name;type:varchar(64);uniqueIndex;" json:"tableName"` // 数据库表名
	TableRemark string `gorm:"type:varchar(255)" json:"tableRemark"`                             // 表备注
	Engine      string `gorm:"type:varchar(64)" json:"engine"`                                   // 存储引擎
	Charset     string `gorm:"type:varchar(64)" json:"charset"`                                  // 字符集
	CreateTime  int64  `json:"createTime"`                                                       // 创建时间
	Onlyread    int    `json:"onlyread"`                                                         // 只读
	utils.PageList
}

func (a *DevDbtable) RefResh() error {
	if a.ID == "" {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Model(a).First(a).Error
	return err
}

func (a *DevDbtable) Add() error {
	a.ID = utils.GetUUID()
	err := db.Db.Create(a).Error
	return err
}

func (a *DevDbtable) Edit() error {
	if a.ID == "" {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Model(a).Updates(a).Error
	return err
}

func (a *DevDbtable) Del() error {
	if a.ID == "" {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Model(a).Delete(a).Error
	return err
}

func (a *DevDbtable) GetPage(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []DevDbtable
	db.Db.Model(&DevDbtable{}).Where(where, v...).Count(&page.Total).Order("create_time desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	page.List = items
	return page
}

func (a *DevDbtable) GetList(where string, v ...interface{}) []DevDbtable {
	var items []DevDbtable
	db.Db.Model(&DevDbtable{}).Where(where, v...).Order("create_time desc").Find(&items)
	return items
}

func (a *DevDbtable) InbySql() error {
	if a.DbTableName == "" {
		return errors.New("缺少参数tableName")
	}
	// 获取 CREATE TABLE
	var tableName string
	var createSQL string
	row := db.Db.Raw("SHOW CREATE TABLE `" + a.DbTableName + "`").Row()
	if err := row.Scan(&tableName, &createSQL); err != nil {
		return err
	}

	// 解析为通用结构
	tableInfo, cols, _, err := ParseCreateTable(createSQL)
	if err != nil {
		return err
	}

	// 保存或更新表记录
	var table DevDbtable
	if err := db.Db.Where("table_name = ?", a.DbTableName).First(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			table = DevDbtable{
				ID:          utils.GetUUID(),
				DbTableName: a.DbTableName,
				TableRemark: tableInfo.TableRemark,
				Engine:      tableInfo.Engine,
				Charset:     tableInfo.Charset,
				CreateTime:  time.Now().Unix(),
			}
			if err := db.Db.Create(&table).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		table.TableRemark = tableInfo.TableRemark
		table.Engine = tableInfo.Engine
		table.Charset = tableInfo.Charset
		if err := db.Db.Save(&table).Error; err != nil {
			return err
		}
	}

	// 先删除已有字段然后插入新的字段
	if err := db.Db.Where("table_id = ?", table.ID).Delete(&DevDbtableColumns{}).Error; err != nil {
		return err
	}
	for i := range cols {
		cols[i].TableId = table.ID
		if err := db.Db.Create(&cols[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

// 执行注册数据库
func regDevDbtableTable() error {
	if db.Db == nil {
		return errors.New("db not initialized")
	}
	return db.Db.AutoMigrate(&DevDbtable{}, &DevDbtableColumns{})
}

type DevDbtableColumns struct {
	Id           string `json:"id" form:"id" gorm:"primaryKey"`      // 主键
	TableId      string `json:"tableId" form:"tableId" gorm:"index"` // 表格ID
	Column       string `json:"column" form:"column"`                // 字段名称
	ColumnType   string `json:"columnType" form:"columnType"`        // 字段类型
	ColumnLong   int    `json:"columnLong" form:"columnLong"`        // 字段长度
	ColumnScale  int    `json:"columnScale" form:"columnScale"`      // 小数位或 scale
	IsUnsigned   int    `json:"isUnsigned" form:"isUnsigned"`        // 是否 unsigned
	NotNull      int    `json:"notNull" form:"notNull"`              // 是否 NOT NULL
	ColumnValue  string `json:"columnValue" form:"columnValue"`      // 字段默认值
	ColumnKey    string `json:"columnKey" form:"columnKey"`          // KEY 信息，如 PRI
	Extra        string `json:"extra" form:"extra"`                  // 额外信息，如 auto_increment
	ColumnCom    string `json:"columnCom" form:"columnCom"`          // 字段注释
	Sort         int    `json:"sort" form:"sort"`                    // 排序
	IsSearch     int    `json:"isSearch" form:"isSearch"`            // 是否搜索字段
	IsFormshow   int    `json:"isFormshow" form:"isFormshow"`        // 表单是否显示
	IsDict       int    `json:"isDict" form:"isDict"`                // 是否是字典
	DictCode     string `json:"dictCode" form:"dictCode"`            // 字典code
	FormComplate string `json:"formComplate" form:"formComplate"`    // 表单组件
	FormProps    string `json:"formProps" form:"formProps"`          // 表单组件参数
	CreateTime   int64  `json:"createTime" form:"createTime"`        // 创建时间
}

func (a *DevDbtable) GetColumns() []DevDbtableColumns {
	var items []DevDbtableColumns
	db.Db.Model(&DevDbtableColumns{}).Where("table_id = ?", a.ID).Order("sort asc").Find(&items)
	return items
}

// 新增
func (a *DevDbtableColumns) Add() error {
	a.Id = utils.GetUUID()
	err := db.Db.Create(a).Error
	return err
}

// 编辑
func (a *DevDbtableColumns) Edit() error {
	if len(a.Id) != 32 {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Updates(a).Error
	return err
}

// 删除
func (a *DevDbtableColumns) Del() error {
	if len(a.Id) != 32 {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Model(a).Delete(a).Error
	return err
}

// 获取分页
func (a *DevDbtableColumns) GetPage(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []DevDbtableColumns
	db.Db.Model(&DevDbtableColumns{}).Where(where, v...).Count(&page.Total).Order("create_time desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	page.List = items
	return page
}

// 获取列表
func (a *DevDbtableColumns) GetList() []DevDbtableColumns {
	var items []DevDbtableColumns
	db.Db.Model(&DevDbtableColumns{}).Order("create_time desc").Find(&items)
	return items
}
