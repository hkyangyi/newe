package model

import "errors"

type AdminAuth struct {
	Token  string    `json:"accessToken"`
	Merdb  SysMember `json:"userInfo"` //用户数据
	Depart SysDepart `json:"depart"`   //组织结构
}

// 刷新数据
func (a *AdminAuth) RefreshByMerdb() error {
	if len(a.Merdb.ID) != 32 {
		return errors.New("无用户数据")
	}

	a.Merdb.Refresh()
	a.Depart.ID = a.Merdb.DepartId
	if len(a.Depart.ID) != 32 {
		return errors.New("用户部门结构未设置，请从新设置")
	}

	err := a.Depart.Refresh()
	if err != nil {
		return errors.New("用户部门结构未设置或已删除，请重新设置该用户")
	}
	return nil
}
