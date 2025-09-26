package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hkyangyi/newe/common/db"
	"github.com/hkyangyi/newe/common/utils"
	"gorm.io/gorm"
)

type SysMenus struct {
	ID         string `gorm:"primaryKey;size:64" json:"id" form:"id"` // 主键ID
	Pid        string `json:"pid" form:"pid"`                         // 父级ID（根节点为空或0）
	Level      int8   `json:"level" form:"level"`                     // 层级（0开始，写入时自动计算）
	Name       string `json:"name" form:"name"`                       // 路由名称（前端 route name，需唯一）
	Title      string `json:"title" form:"title"`                     // 菜单标题（显示用）
	Icon       string `json:"icon" form:"icon"`                       // 图标
	ActiveIcon string `json:"activeIcon" form:"activeIcon"`           // 选中时图标

	Path      string `json:"path" form:"path"`           // 相对路径（以/开头，同级唯一）
	FullPath  string `json:"fullPath" form:"fullPath"`   // 递归拼接后的完整路径（冗余存储加速）
	Component string `json:"component" form:"component"` // 前端组件路径
	Redirect  string `json:"redirect" form:"redirect"`   // 重定向路径

	Type      int `json:"type" form:"type"`           // 类型 1目录 2菜单 3按钮
	OrderNo   int `json:"orderNo" form:"orderNo"`     // 排序号（越小越前）
	Status    int `json:"status" form:"status"`       // 1启用 -1禁用
	IsDeleted int `json:"isDeleted" form:"isDeleted"` // 1正常 -1已删除（逻辑删除）

	KeepAlive          int `json:"keepAlive" form:"keepAlive"`                   // 1是 -1否 是否缓存 keep-alive
	HideInMenu         int `json:"hideInMenu" form:"hideInMenu"`                 // 1是 -1否 菜单中隐藏
	HideInTab          int `json:"hideInTab" form:"hideInTab"`                   // 1是 -1否 多标签中隐藏
	HideInBreadcrumb   int `json:"hideInBreadcrumb" form:"hideInBreadcrumb"`     // 1是 -1否 面包屑隐藏
	HideChildrenInMenu int `json:"hideChildrenInMenu" form:"hideChildrenInMenu"` // 1是 -1否 隐藏子菜单

	AuthorityRaw string   `gorm:"column:authority" json:"-" form:"authority"` // 权限数组原始JSON（存到库）
	Authority    []string `gorm:"-" json:"authority"`                         // 展示用权限数组
	Permission   string   `json:"permission" form:"permission"`               // 主权限编码（按钮/菜单标识）

	Badge         string `json:"badge" form:"badge"`                 // 徽标文本
	BadgeType     string `json:"badgeType" form:"badgeType"`         // 徽标类型（样式类型）
	BadgeVariants string `json:"badgeVariants" form:"badgeVariants"` // 徽标风格变体

	FullPathKey int    `json:"fullPathKey" form:"fullPathKey"` // 1是 -1否 以 fullPath 作为激活匹配 key
	ActivePath  string `json:"activePath" form:"activePath"`   // 强制激活路径（用于高亮）

	AffixTab      int `json:"affixTab" form:"affixTab"`           // 1是 -1否 固定在标签页
	AffixTabOrder int `json:"affixTabOrder" form:"affixTabOrder"` // 固定标签排序优先级

	IframeSrc                string `json:"iframeSrc" form:"iframeSrc"`                               // iframe 地址
	IgnoreAccess             int    `json:"ignoreAccess" form:"ignoreAccess"`                         // 1是 -1否 忽略权限
	Link                     string `json:"link" form:"link"`                                         // 外部链接
	MaxNumOfOpenTab          int    `json:"maxNumOfOpenTab" form:"maxNumOfOpenTab"`                   // 最大标签数
	MenuVisibleWithForbidden int    `json:"menuVisibleWithForbidden" form:"menuVisibleWithForbidden"` // 1是 -1否 无权限仍展示
	OpenInNewWindow          int    `json:"openInNewWindow" form:"openInNewWindow"`                   // 1是 -1否 新窗口打开
	NoBasicLayout            int    `json:"noBasicLayout" form:"noBasicLayout"`                       // 1是 -1否 不使用基础布局

	QueryParamsRaw string            `gorm:"column:query_params" json:"-" form:"queryParams"` // 默认附加 query 参数原始 JSON
	Query          map[string]string `gorm:"-" json:"queryParams"`                            // 解析后的 query 参数

	BindApi int `json:"bindApi" form:"bindApi"` // 绑定的后端接口 ID（预留）

	CreatedAt int64 `json:"createdAt"` // 创建时间 Unix 时间戳
	UpdatedAt int64 `json:"updatedAt"` // 更新时间 Unix 时间戳

	Children []SysMenus `gorm:"-" json:"children,omitempty"` // 子节点（递归组装，不入库）
}

func (a *SysMenus) Add() error {
	a.ID = utils.GetUUID()
	now := time.Now().Unix()
	a.CreatedAt = now
	a.UpdatedAt = now
	if a.Status == 0 {
		a.Status = 1
	}
	if a.IsDeleted == 0 {
		a.IsDeleted = 1
	} // 默认正常=1
	a.normalizeBasic()
	a.normalizeFlags()
	return db.Db.Create(a).Error
}

func (a *SysMenus) Edit() error {
	if a.ID == "" {
		return errors.New("缺少参数ID")
	}
	a.normalizeBasic()
	a.UpdatedAt = time.Now().Unix()
	return db.Db.Model(a).Where("id = ?", a.ID).Updates(a).Error
}

func (a *SysMenus) Del() error {
	if a.ID == "" {
		return errors.New("缺少参数ID")
	}
	return db.Db.Model(&SysMenus{}).Where("id = ?", a.ID).Updates(map[string]interface{}{"is_deleted": -1, "updated_at": time.Now().Unix()}).Error
}

// ensure/normalize fields before write
func (m *SysMenus) BeforeCreate(tx *gorm.DB) error { return m.validate(tx, true) }
func (m *SysMenus) BeforeUpdate(tx *gorm.DB) error { return m.validate(tx, false) }

func (m *SysMenus) normalizeBasic() {
	m.Name = strings.TrimSpace(m.Name)
	m.Title = strings.TrimSpace(m.Title)
	m.Path = strings.TrimSpace(m.Path)
	if m.Path == "" {
		m.Path = "/" + sanitizeRouteName(m.Name)
	}
	if !strings.HasPrefix(m.Path, "/") {
		m.Path = "/" + m.Path
	}
	if m.Type == 1 && m.Component == "" { // 目录默认布局组件
		m.Component = ""
	}
	// FullPath 由上层查询时补，若没有父级直接等于 Path
	if m.Pid == "" || m.Pid == "0" {
		m.FullPath = m.Path
		m.Level = 0
	}
}

// 归一化所有 int 标志字段：不为1则设为-1
func (m *SysMenus) normalizeFlags() {
	flags := []*int{&m.KeepAlive, &m.HideInMenu, &m.HideInTab, &m.HideInBreadcrumb, &m.HideChildrenInMenu,
		&m.FullPathKey, &m.AffixTab, &m.IgnoreAccess, &m.MenuVisibleWithForbidden, &m.OpenInNewWindow, &m.NoBasicLayout}
	for _, f := range flags {
		if *f != 1 {
			*f = -1
		}
	}
}

func (m *SysMenus) validate(tx *gorm.DB, isCreate bool) error {
	fmt.Println(*m)
	m.normalizeBasic()
	m.normalizeFlags()
	m.prepareJSON()
	if m.Name == "" {
		return errors.New("路由名称不能为空")
	}
	if m.Title == "" {
		return errors.New("菜单标题不能为空")
	}
	if m.Type != 1 && m.Type != 2 && m.Type != 3 {
		return errors.New("类型必须为1/2/3")
	}
	if m.Status != 1 && m.Status != -1 {
		m.Status = 1
	}
	if m.IsDeleted != 1 && m.IsDeleted != -1 {
		m.IsDeleted = 1
	}

	// 唯一 name
	var cnt int64
	q := tx.Model(&SysMenus{}).Where("name = ? AND is_deleted = 1", m.Name)
	if !isCreate {
		q = q.Where("id <> ?", m.ID)
	}
	if err := q.Count(&cnt).Error; err == nil && cnt > 0 {
		return errors.New("路由名称已存在")
	}

	// 同级 path 唯一
	cnt = 0
	q2 := tx.Model(&SysMenus{}).Where("pid = ? AND path = ? AND is_deleted = 1", m.Pid, m.Path)
	if !isCreate {
		q2 = q2.Where("id <> ?", m.ID)
	}
	if err := q2.Count(&cnt).Error; err == nil && cnt > 0 {
		return errors.New("同级路径已存在")
	}

	// JSON 解析
	if m.AuthorityRaw != "" && !json.Valid([]byte(m.AuthorityRaw)) {
		return errors.New("authority 必须为有效JSON")
	}
	if m.QueryParamsRaw != "" && !json.Valid([]byte(m.QueryParamsRaw)) {
		return errors.New("queryParams 必须为有效JSON")
	}
	return nil
}

// prepareJSON 保证写入数据库前 JSON 列为合法 JSON 文本（避免空字符串导致 MySQL 3140 错误）
func (m *SysMenus) prepareJSON() {
	// authority
	if len(m.Authority) > 0 {
		if b, err := json.Marshal(m.Authority); err == nil {
			m.AuthorityRaw = string(b)
		}
	} else if strings.TrimSpace(m.AuthorityRaw) == "" { // 未提供则设为空数组
		m.AuthorityRaw = "[]"
	} else if !json.Valid([]byte(m.AuthorityRaw)) { // 非法则回退
		m.AuthorityRaw = "[]"
	}

	// query params
	if len(m.Query) > 0 {
		if b, err := json.Marshal(m.Query); err == nil {
			m.QueryParamsRaw = string(b)
		}
	} else if strings.TrimSpace(m.QueryParamsRaw) == "" {
		m.QueryParamsRaw = "{}"
	} else if !json.Valid([]byte(m.QueryParamsRaw)) {
		m.QueryParamsRaw = "{}"
	}
}

// AfterFind 读取时反序列化 JSON
func (m *SysMenus) AfterFind(tx *gorm.DB) (err error) {
	if m.AuthorityRaw != "" && json.Valid([]byte(m.AuthorityRaw)) {
		var arr []string
		if e := json.Unmarshal([]byte(m.AuthorityRaw), &arr); e == nil {
			m.Authority = arr
		}
	}
	if m.QueryParamsRaw != "" && json.Valid([]byte(m.QueryParamsRaw)) {
		var mp map[string]string
		if e := json.Unmarshal([]byte(m.QueryParamsRaw), &mp); e == nil {
			m.Query = mp
		}
	}
	return nil
}

func sanitizeRouteName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "menu"
	}
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

func GetMenusList() []SysMenus {
	var items []SysMenus
	db.Db.Where("is_deleted = 1").Order("order_no asc").Find(&items)
	return items
}

// 根据部门ID获取菜单
func SysMenusListGetByDepart(departId string) []SysMenus {
	var rules []SysDepartRules
	db.Db.Where("depart_id = ?", departId).Find(&rules)
	ids := make([]string, 0, len(rules))
	for _, r := range rules {
		ids = append(ids, r.MenuId)
	}
	if len(ids) == 0 {
		return nil
	}
	var res []SysMenus
	db.Db.Where("is_deleted = 1 AND status = 1 AND id in ?", ids).Order("order_no asc").Find(&res)
	return res
}

func SysButtonGetList() []string {
	var items []SysMenus
	db.Db.Where("is_deleted = 1 AND status = 1 AND type = 3").Order("order_no asc").Find(&items)
	btns := make([]string, 0, len(items))
	for _, v := range items {
		if v.Permission != "" {
			btns = append(btns, v.Permission)
		}
	}
	return btns
}

// 获取按钮
func SysButtonGetByDepart(departId string) []string {
	var rules []SysDepartRules
	db.Db.Where("depart_id = ?", departId).Find(&rules)
	ids := make([]string, 0, len(rules))
	for _, r := range rules {
		ids = append(ids, r.MenuId)
	}
	if len(ids) == 0 {
		return nil
	}
	var res []SysMenus
	db.Db.Where("is_deleted = 1 AND status = 1 AND type = 3 AND id in ?", ids).Order("order_no asc").Find(&res)
	btns := make([]string, 0, len(res))
	for _, v := range res {
		if v.Permission != "" {
			btns = append(btns, v.Permission)
		}
	}
	return btns
}

func (a *SysMenus) GetList(where string, v ...interface{}) []SysMenus {
	var items []SysMenus
	db.Db.Where("is_deleted = 1").Where(where, v...).Order("order_no asc").Find(&items)
	return buildTree(items, "")
}

// GetMenuTree 返回全部正常且启用的菜单树（不含被禁用或删除的）
func GetMenuTree() []SysMenus {
	var items []SysMenus
	db.Db.Where("is_deleted = 1").Order("order_no asc").Find(&items)
	return buildTree(items, "")
}

func buildTree(items []SysMenus, pid string) []SysMenus {
	var res []SysMenus
	for i := range items {
		if items[i].Pid == pid {
			node := items[i]
			node.Children = buildTree(items, node.ID)
			// 计算 fullPath
			if node.Pid == "" || node.Pid == "0" {
				node.FullPath = node.Path
				node.Level = 0
			} else {
				// 查找父级 fullPath
				for _, p := range items {
					if p.ID == node.Pid {
						node.FullPath = strings.TrimRight(p.FullPath, "/") + node.Path
						node.Level = p.Level + 1
						break
					}
				}
			}
			res = append(res, node)
		}
	}
	return res
}

// 用户修改信息
func (a *SysMember) EditProfile() error {
	if a.ID == "" {
		return errors.New("缺少参数ID")
	}

	err := db.Db.Model(a).Where("id = ?", a.ID).Updates(map[string]interface{}{
		"nickname":  a.Nickname,
		"real_name": a.RealName,
		"sex":       a.Sex,
		"mp":        a.Mp,

		"update_time": time.Now().Unix(),
	}).Error

	return err
}
