package cusmodel

import "github.com/hkyangyi/newe/common/db"

type CusAuth struct {
	Isadmin bool           `json:"isadmin" gorm:""` // 是否管理员
	Token   string         `json:"accessToken"`
	CDB     CustomerUsers  `json:"cdb" gorm:"-"`
	MerDb   CustomerMember `json:"merdb" gorm:"-"`
}

// 刷新用户信息
func (a *CusAuth) RefreshByMerdb() error {
	// if a.Isadmin {
	// 	err := db.Db.Table("customer_users").Where("id = ?", a.CDB.Id).First(&a.CDB).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// 	a.CDB.Password = ""
	// 	a.MerDb = CustomerMember{}
	// 	return nil
	// } else {

	err := db.Db.Table("customer_member").Where("id = ?", a.MerDb.ID).First(&a.MerDb).Error
	if err != nil {
		return err
	}

	a.MerDb.Password = ""
	err = db.Db.Table("customer_users").Where("id = ?", a.MerDb.Cid).First(&a.CDB).Error
	if err != nil {
		return err
	}
	a.CDB.Password = ""
	return nil
	//}
}
