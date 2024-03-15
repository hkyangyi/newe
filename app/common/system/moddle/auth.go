package moddle

import "errors"

type AdminAuth struct {
	Token  string    `json:"token"`
	Merdb  SysMember `json:"merdb"`  //用户数据
	Role   SysRole   `json:"role"`   //角色数据
	Depart SysDepart `json:"depart"` //组织结构
}

// 刷新数据
func (a *AdminAuth) RefreshByMerdb() error {
	if len(a.Merdb.ID) != 32 {
		return errors.New("无用户数据")
	}

	a.Merdb.Refresh()
	a.Role.ID = a.Merdb.RoleId
	if len(a.Role.ID) != 32 {
		return errors.New("用户角色出错，请重新设置该用户")
	}
	err := a.Role.Refresh()
	if err != nil {
		return errors.New("用户角色出差或已删除，请重新设置该用户")
	}

	a.Depart.ID = a.Merdb.DepartId
	if len(a.Depart.ID) != 32 {
		return errors.New("用户部门结构未设置，请从新设置")
	}

	err = a.Depart.Refresh()
	if err != nil {
		return errors.New("用户部门结构未设置或已删除，请重新设置该用户")
	}
	return nil
}
