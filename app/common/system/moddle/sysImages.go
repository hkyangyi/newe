package moddle

import (
	"newe/common/base"
	"newe/common/utils"
)

type SysImages struct {
	ID         string `gorm:"primary_key" json:"id"` //
	Path       string `json:"path"`                  //
	URL        string `json:"url"`                   //
	Size       int64  `json:"size"`                  //
	CreateTime int64  `json:"createTime"`            //
	CreateName string `json:"createName"`            //
	DepartId   string `json:"departId"`              //
}

// 获取图片列表
func ImagesGetList(pl utils.PageList, departId string) (utils.PageList, error) {
	var items []SysImages
	err := base.MYDB.Model(SysImages{}).Order("create_time desc").Count(&pl.Total).Offset(pl.GetOffice()).Limit(pl.PageSize).Find(&items).Error
	pl.List = items
	return pl, err
}

// 添加图片
func ImagesAdd(data SysImages) error {
	data.ID = utils.GetUUID()
	err := base.MYDB.Create(&data).Error
	return err
}
