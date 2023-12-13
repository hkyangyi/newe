package utils

type PageList struct {
	List     interface{} `gorm:"-" json:"list"`
	Page     int         `gorm:"-" json:"page" form:"page"`         //当前页
	PageSize int         `gorm:"-" json:"pageSize" form:"pageSize"` //每页数量
	Total    int64       `gorm:"-" json:"total"`                    //总条数
}

func (p *PageList) GetOffice() int {
	result := 0
	if p.Page > 0 {
		result = (p.Page - 1) * p.PageSize
	}
	return result
}
