package model

import (
	"errors"
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

type SysMessage struct {
	ID        string `gorm:"primary_key;size:32" json:"id" form:"id"`              // 消息ID
	Uid       string `gorm:"index:idx_user_created;size:64" json:"uid" form:"uid"` // 用户ID（与登录uuid一致）
	Title     string `gorm:"size:255" json:"title"`                                // 标题
	Intro     string `gorm:"size:512" json:"intro"`                                // 简介
	URL       string `gorm:"size:512" json:"url"`                                  // 跳转地址
	Picurl    string `gorm:"size:512" json:"picurl"`                               // 图片地址
	MsgType   string `gorm:"size:64" json:"msgType" form:"msgType"`                // 消息类型
	CreatedAt int64  `gorm:"index:idx_user_created" json:"createdAt"`              // 创建时间（秒）
	ReadAt    int64  `json:"readAt"`                                               // 阅读时间（0 表示未读）
	utils.PageList
}

func (SysMessage) TableName() string { return "sys_message" }

// 执行注册数据库
func regSysMessageTable() error {
	if db.Db == nil {
		return errors.New("db not initialized")
	}
	return db.Db.AutoMigrate(&SysMessage{})
}

func (a *SysMessage) Add() error {
	err := db.Db.Create(a).Error
	return err
}
func (a *SysMessage) Edit() error {
	if len(a.ID) != 32 {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Model(a).Updates(a).Error
	return err
}

func (a *SysMessage) ReadEnd() error {
	if len(a.ID) != 32 {
		return errors.New("缺少参数ID")
	}
	err := db.Db.Model(a).Where("read_at=0").Update("read_at", time.Now().Unix()).Error
	return err
}

func (a *SysMessage) Del() error {
	err := db.Db.Model(a).Delete(a).Error
	return err
}

func (a *SysMessage) GetList(where string, v ...interface{}) ([]SysMessage, error) {
	var list []SysMessage
	err := db.Db.Where(where, v...).Order("created_at desc").Find(&list).Error
	return list, err
}

func (a *SysMessage) GetPage(page utils.PageList, where string, v ...interface{}) (utils.PageList, error) {
	var list []SysMessage
	err := db.Db.Model(&SysMessage{}).Where(where, v...).Order("created_at desc").Count(&page.Total).Offset(page.GetOffice()).Limit(page.PageSize).Find(&list).Error
	page.List = list
	return page, err
}

func (a *SysMessage) GetOne(where string, v ...interface{}) (SysMessage, error) {
	var one SysMessage
	err := db.Db.Where(where, v...).First(&one).Error
	return one, err
}

func (a *SysMessage) SetRead() error {
	err := db.Db.Model(a).Update("read_at", time.Now().Unix()).Error
	return err
}
