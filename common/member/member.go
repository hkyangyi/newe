package member

import (
	"fmt"

	"github.com/hkyangyi/newe/common/ws"
	"github.com/hkyangyi/newe/model"
)

// 获取在线用户
func GetOnlineMembers(orgcode string) []model.SysMember {
	//获取ws管理器地址
	wsgmag := ws.MainHub.Get("NeweSysAdmin")
	var members []model.SysMember
	//获取所有在线用户key
	keys := wsgmag.OnlineKeys()
	for _, key := range keys {
		//解析出用户ID
		var uid string
		_, err := fmt.Sscanf(key, "Ws_NeweSysAdmin_%s", &uid)
		if err != nil {
			continue
		}
		//根据用户ID获取用户信息
		var member model.SysMember
		err = model.GetMemberByID(uid, &member)
		if err != nil {
			continue
		}
		members = append(members, member)
	}
	return members
}

// IsMemberOnline 检测用户是否在线
func IsMemberOnline(uid string) bool {
	wsgmag := ws.MainHub.Get("NeweSysAdmin")
	key := "Ws_NeweSysAdmin_" + uid
	return wsgmag.Has(key)
}
