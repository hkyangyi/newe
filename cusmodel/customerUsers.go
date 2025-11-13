package cusmodel

import "github.com/hkyangyi/newe/common/db"

type CustomerUsers struct {
	Id                 string `json:"id" form:"id"  gorm:"primary_key"`
	Username           string `json:"username" form:"username" `                     // 主理人账号
	Password           string `json:"password" form:"password" `                     // 主理人密码
	Uid                string `json:"uid" form:"uid" `                               // 主理人用户ID绑定
	Realname           string `json:"realname" form:"realname" `                     // 真实姓名
	Headimgurl         string `json:"headimgurl" form:"headimgurl" `                 // 主理人照片
	Mp                 string `json:"mp" form:"mp" `                                 // 手机号码
	Slogan             string `json:"slogan" form:"slogan" `                         // 标语
	Nickname           string `json:"nickname" form:"nickname" `                     // 昵称
	CustomerName       string `json:"customerName" form:"customerName" `             // 企业名称
	CustomerNickname   string `json:"customerNickname" form:"customerNickname" `     // 店铺名称
	CustomerHeadimgurl string `json:"customerHeadimgurl" form:"customerHeadimgurl" ` // 店铺头像
	CustomerBgurl      string `json:"customerBgurl" form:"customerBgurl" `           // 店铺背景图
	CustomerVedio      string `json:"customerVedio" form:"customerVedio" `           // 店铺宣传视频
	Status             int    `json:"status" form:"status" `                         // 状态 1正常 -1 停用
	CreateTime         int64  `json:"createTime" form:"createTime" `                 // 创建时间
	UpdateTime         int64  `json:"updateTime" form:"updateTime" `                 // 更新时间
	UpdateUid          string `json:"updateUid" form:"updateUid" `                   // 更新人
	UpdateRealname     string `json:"updateRealname" form:"updateRealname" `         // 最后更新人
}

func FindCusUserByUsername(user, pass string) (CustomerUsers, error) {
	var data CustomerUsers
	err := db.Db.Model(&data).Where("username = ? and password = ?", user, pass).First(&data).Error
	return data, err
}
