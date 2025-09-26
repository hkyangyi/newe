// Package message 提供简单的站内消息队列：入库并在用户在线时通过 WS 推送。
package message

import (
	"errors"
	"time"

	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/common/ws"
	"github.com/hkyangyi/newe/model"
)

// AddMessage 入库一条消息；若用户在线则立即通过 WS 推送
// 返回保存后的消息对象（含ID与时间）
func AddMessage(uid, title, intro, url, picurl, msgType string) (*model.SysMessage, error) {
	if uid == "" || title == "" {
		return nil, errors.New("userID and title required")
	}
	m := &model.SysMessage{
		ID:        utils.GetUUID(),
		Uid:       uid,
		Title:     title,
		Intro:     intro,
		URL:       url,
		Picurl:    picurl,
		MsgType:   msgType,
		CreatedAt: time.Now().Unix(),
		ReadAt:    0,
	}
	if err := m.Add(); err != nil {
		return nil, err
	}

	wskey := "Ws_NeweSysAdmin_" + uid

	//获取ws管理器地址
	wsgmag := ws.MainHub.Get("NeweSysAdmin")
	// 若 WS 在线则推送
	if wsgmag.Has(wskey) {
		_ = wsgmag.SendJSON(wskey, map[string]any{
			"type": "notify",
			"data": m,
		})
	}

	return m, nil
}

func MarkRead(id string) error {
	if id == "" {
		return errors.New("id required")
	}
	row := model.SysMessage{ID: id}
	return row.ReadEnd()
}
