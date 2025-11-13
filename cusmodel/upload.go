package cusmodel

import (
	"os"
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

// SysUpload 上传记录表（图片/文件统一记录）
// - 采用 1/-1 布尔语义：Status=1 正常；IsDeleted=1 正常，-1 删除；IsImage=1 图片，-1 文件
// - 便于按 Md5/Sha1 去重复用；同时保留存储提供商、路径、URL 等信息
type CustomerUpload struct {
	ID       string `gorm:"primaryKey;autoIncrement" json:"id"`     // 主键（自增）
	Cid      string `gorm:"type:varchar(64);index" json:"cid"`      // 客户ID
	MemberID string `gorm:"type:varchar(64);index" json:"memberId"` // 上传人ID（member_id）

	OriginalName string `gorm:"type:varchar(255)" json:"originalName"` // 原始文件名（含扩展名）
	FileName     string `gorm:"type:varchar(255)" json:"fileName"`     // 存储文件名（唯一/随机名）
	Ext          string `gorm:"type:varchar(16)" json:"ext"`           // 扩展名（不含.）
	MimeType     string `gorm:"type:varchar(128)" json:"mimeType"`     // MIME 类型
	Size         int64  `json:"size"`                                  // 文件大小（字节）
	Md5          string `gorm:"type:char(32);index" json:"md5"`        // MD5（用于去重）
	Sha1         string `gorm:"type:char(40);index" json:"sha1"`       // SHA1（用于去重）

	Path string `gorm:"type:varchar(255)" json:"path"` // 存储路径（相对路径）
	URL  string `gorm:"type:varchar(512)" json:"url"`  // 访问URL（CDN或直链）

	IsImage int `gorm:"type:tinyint" json:"isImage"` // 1 图片 -1 非图片
	Width   int `json:"width"`                       // 宽（仅图片）
	Height  int `json:"height"`                      // 高（仅图片）

	Status    int   `gorm:"type:tinyint;index" json:"status"` // 1 正常 -1 停用
	CreatedAt int64 `json:"createdAt"`                        // 创建时间（Unix）
	utils.PageList
}

func (a *CustomerUpload) Add() error {
	a.CreatedAt = time.Now().Unix()
	err := db.Db.Create(a).Error
	return err
}

func (a *CustomerUpload) Del() error {
	//删除文件
	os.Remove(a.Path)
	//删除记录
	err := db.Db.Delete(a).Error
	return err
}

// 批量删除
func (a *CustomerUpload) DelBatch(ids []int64) error {
	var list []CustomerUpload
	err := db.Db.Where("id in ?", ids).Find(&list).Error
	if err != nil {

		return err
	}
	for _, v := range list {
		//删除文件
		os.Remove(v.Path)
	}
	err = db.Db.Where("id in ?", ids).Delete(&CustomerUpload{}).Error
	return err
}

func (a *CustomerUpload) GetPage(page utils.PageList, where string, v ...interface{}) (utils.PageList, error) {
	var items []CustomerUpload
	err := db.Db.Model(&CustomerUpload{}).Where(where, v...).Count(&page.Total).Order("created_at desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items).Error
	page.List = items
	return page, err
}
