package model

import (
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
)

// SysOperationLog 操作日志表模型
// 记录后台所有敏感操作、接口调用、异常等，便于审计与追溯
// 路由模板与实际路径均记录，兼容 RESTful 路由
// 分组、响应码等字段便于统计与分析

type SysOperationLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`                    // 日志主键
	MemberID     string    `gorm:"type:varchar(32);index" json:"memberId"`                          // 操作人ID
	Username     string    `gorm:"type:varchar(64)" json:"username" form:"username"`                // 操作人用户名
	Method       string    `gorm:"type:varchar(10)" json:"method" form:"method"`                    // 请求方法
	RoutePattern string    `gorm:"type:varchar(128);index" json:"routePattern" form:"routePattern"` // 路由模板（如 /api/v1/user/:id）
	ActualPath   string    `gorm:"type:varchar(128)" json:"actualPath" form:"actualPath"`           // 实际请求路径（如 /api/v1/user/123）
	GroupName    string    `gorm:"type:varchar(32);index" json:"groupName"`                         // 接口分组（如 user、order，可选）
	ApiID        int64     `gorm:"type:bigint;index" json:"apiId"`                                  // 关联 sys_api_list.id
	IP           string    `gorm:"type:varchar(45)" json:"ip"`                                      // 操作人IP
	Params       string    `gorm:"type:text" json:"params"`                                         // 请求参数（JSON字符串）
	Result       string    `gorm:"type:varchar(16)" json:"result" form:"result"`                    // 操作结果（success/fail）
	Message      string    `gorm:"type:varchar(255)" json:"message"`                                // 结果说明或错误信息
	StatusCode   int       `gorm:"type:int" json:"status_code"`                                     // HTTP响应码
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`                                // 操作时间
	StartTime    time.Time `gorm:"-" json:"startTime" form:"startTime"`                             // 开始时间（查询用）
	EndTime      time.Time `gorm:"-" json:"endTime" form:"endTime"`                                 // 结束时间（查询用）
	utils.PageList
}

func (SysOperationLog) TableName() string { return "sys_operation_log" }

func (a *SysOperationLog) Add() error {
	err := db.Db.Create(a).Error
	return err
}

// EnsureApiAndGetID 确保当前日志对应的接口在 sys_api_list 中存在，并返回接口ID
func (a *SysOperationLog) EnsureApiAndGetID(name string) (int64, error) {
	var api SysApiList
	// 优先按 route_pattern + method 匹配（模板维度）
	q := db.Db.Where("route = ? AND method = ?", a.RoutePattern, a.Method).First(&api)
	if q.Error == nil && api.ID > 0 {
		return api.ID, nil
	}
	// 不存在则创建
	api = SysApiList{
		ParentID:  0,
		GroupName: a.GroupName,
		Name:      name,
		Route:     a.RoutePattern,
		Method:    a.Method,
		Desc:      "",
	}
	if err := db.Db.Create(&api).Error; err != nil {
		return 0, err
	}
	return api.ID, nil
}

func (a *SysOperationLog) GetPage(page utils.PageList, where string, v ...interface{}) utils.PageList {
	var items []SysOperationLog
	db.Db.Model(&SysOperationLog{}).Where(where, v...).Count(&page.Total).Order("created_at desc").Offset(page.GetOffice()).Limit(page.PageSize).Find(&items)
	page.List = items
	return page
}
