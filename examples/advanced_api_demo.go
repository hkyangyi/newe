package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// 高级API演示客户端
type AdvancedNeweClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// 通用响应结构
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
	Success string      `json:"success"`
}

// 用户信息
type UserInfo struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	RealName  string `json:"realname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    int    `json:"status"`
	DepartId  string `json:"depart_id"`
	OrgCode   string `json:"org_code"`
	RoleId    string `json:"role_id"`
	CreatedAt int64  `json:"create_time"`
}

// 字典项
type DictItem struct {
	ID       string `json:"id"`
	DictCode string `json:"dict_code"`
	DictName string `json:"dict_name"`
	DictValue string `json:"dict_value"`
	ParentId string `json:"parent_id"`
	SortNo   int    `json:"sort_no"`
	Status   int    `json:"status"`
}

// 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// 创建高级客户端
func NewAdvancedNeweClient(baseURL string) *AdvancedNeweClient {
	return &AdvancedNeweClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				IdleConnTimeout:     30 * time.Second,
				DisableCompression: false,
			},
		},
	}
}

// 发送请求辅助函数
func (c *AdvancedNeweClient) sendRequest(method, endpoint string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", c.Token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Newe-API-Client/1.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// 用户登录
func (c *AdvancedNeweClient) Login(username, password string) error {
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	body, err := c.sendRequest("POST", "/admin/login", loginData)
	if err != nil {
		return err
	}

	var response ApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Code != 200 {
		return fmt.Errorf("登录失败: %s", response.Message)
	}

	// 提取token
	if result, ok := response.Result.(map[string]interface{}); ok {
		if token, ok := result["token"].(string); ok {
			c.Token = token
			fmt.Printf("✅ 登录成功，Token: %s\n", token[:20]+"...")
			return nil
		}
	}

	return fmt.Errorf("无法获取token")
}

// 获取用户列表
func (c *AdvancedNeweClient) GetUserList(page, pageSize int) ([]UserInfo, int, error) {
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("pageSize", fmt.Sprintf("%d", pageSize))

	endpoint := "/admin/system/membergetlist?" + params.Encode()
	body, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, 0, err
	}

	var response struct {
		Code    int          `json:"code"`
		Message string       `json:"message"`
		Result  PageResponse `json:"result"`
		Success string       `json:"success"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, 0, err
	}

	if response.Code != 200 {
		return nil, 0, fmt.Errorf("获取用户列表失败: %s", response.Message)
	}

	// 转换结果
	var users []UserInfo
	if list, ok := response.Result.List.([]interface{}); ok {
		for _, item := range list {
			if userData, ok := item.(map[string]interface{}); ok {
				user := UserInfo{
					ID:        getString(userData, "id"),
					Username:  getString(userData, "username"),
					RealName:  getString(userData, "realname"),
					Email:     getString(userData, "email"),
					Phone:     getString(userData, "phone"),
					Status:    getInt(userData, "status"),
					DepartId:  getString(userData, "depart_id"),
					OrgCode:   getString(userData, "org_code"),
					RoleId:    getString(userData, "role_id"),
					CreatedAt: getInt64(userData, "create_time"),
				}
				users = append(users, user)
			}
		}
	}

	return users, response.Result.Total, nil
}

// 获取字典列表
func (c *AdvancedNeweClient) GetDictList(parentId string, page, pageSize int) ([]DictItem, int, error) {
	params := url.Values{}
	if parentId != "" {
		params.Add("parent_id", parentId)
	}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("pageSize", fmt.Sprintf("%d", pageSize))

	endpoint := "/admin/system/dictgetlist?" + params.Encode()
	body, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, 0, err
	}

	var response struct {
		Code    int          `json:"code"`
		Message string       `json:"message"`
		Result  PageResponse `json:"result"`
		Success string       `json:"success"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, 0, err
	}

	if response.Code != 200 {
		return nil, 0, fmt.Errorf("获取字典列表失败: %s", response.Message)
	}

	// 转换结果
	var dictItems []DictItem
	if list, ok := response.Result.List.([]interface{}); ok {
		for _, item := range list {
			if dictData, ok := item.(map[string]interface{}); ok {
				dict := DictItem{
					ID:        getString(dictData, "id"),
					DictCode:  getString(dictData, "dict_code"),
					DictName:  getString(dictData, "dict_name"),
					DictValue: getString(dictData, "dict_value"),
					ParentId:  getString(dictData, "parent_id"),
					SortNo:    getInt(dictData, "sort_no"),
					Status:    getInt(dictData, "status"),
				}
				dictItems = append(dictItems, dict)
			}
		}
	}

	return dictItems, response.Result.Total, nil
}

// 根据字典代码获取字典项
func (c *AdvancedNeweClient) GetDictByCode(dictCode string) ([]DictItem, error) {
	params := url.Values{}
	params.Add("dict_code", dictCode)

	endpoint := "/admin/dict?" + params.Encode()
	body, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response ApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("获取字典失败: %s", response.Message)
	}

	var dictItems []DictItem
	if result, ok := response.Result.([]interface{}); ok {
		for _, item := range result {
			if dictData, ok := item.(map[string]interface{}); ok {
				dict := DictItem{
					ID:        getString(dictData, "id"),
					DictCode:  getString(dictData, "dict_code"),
					DictName:  getString(dictData, "dict_name"),
					DictValue: getString(dictData, "dict_value"),
					ParentId:  getString(dictData, "parent_id"),
					SortNo:    getInt(dictData, "sort_no"),
					Status:    getInt(dictData, "status"),
				}
				dictItems = append(dictItems, dict)
			}
		}
	}

	return dictItems, nil
}

// 创建用户
func (c *AdvancedNeweClient) CreateUser(userData map[string]interface{}) error {
	body, err := c.sendRequest("POST", "/admin/system/memberadd", userData)
	if err != nil {
		return err
	}

	var response ApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Code != 200 {
		return fmt.Errorf("创建用户失败: %s", response.Message)
	}

	fmt.Println("✅ 用户创建成功")
	return nil
}

// 创建字典项
func (c *AdvancedNeweClient) CreateDictItem(dictData map[string]interface{}) error {
	body, err := c.sendRequest("POST", "/admin/system/dictadd", dictData)
	if err != nil {
		return err
	}

	var response ApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Code != 200 {
		return fmt.Errorf("创建字典项失败: %s", response.Message)
	}

	fmt.Println("✅ 字典项创建成功")
	return nil
}

// 辅助函数：从map中获取字符串值
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// 辅助函数：从map中获取整数值
func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
}

// 辅助函数：从map中获取int64值
func getInt64(data map[string]interface{}, key string) int64 {
	if val, ok := data[key]; ok {
		if num, ok := val.(float64); ok {
			return int64(num)
		}
	}
	return 0
}

// 演示函数
func main() {
	fmt.Println("🚀 Newe 高级API调用演示")
	fmt.Println("========================")

	// 创建客户端
	client := NewAdvancedNeweClient("http://localhost:8080")

	// 用户登录
	fmt.Println("\n1. 用户登录...")
	if err := client.Login("admin", "admin123"); err != nil {
		fmt.Printf("❌ 登录失败: %v\n", err)
		return
	}

	// 获取用户列表
	fmt.Println("\n2. 获取用户列表...")
	users, total, err := client.GetUserList(1, 10)
	if err != nil {
		fmt.Printf("❌ 获取用户列表失败: %v\n", err)
	} else {
		fmt.Printf("✅ 获取到 %d 个用户 (总共 %d 个):\n", len(users), total)
		for i, user := range users {
			fmt.Printf("  %d. %s (%s)\n", i+1, user.Username, user.RealName)
		}
	}

	// 获取字典列表
	fmt.Println("\n3. 获取字典列表...")
	dictItems, total, err := client.GetDictList("", 1, 10)
	if err != nil {
		fmt.Printf("❌ 获取字典列表失败: %v\n", err)
	} else {
		fmt.Printf("✅ 获取到 %d 个字典项 (总共 %d 个):\n", len(dictItems), total)
		for i, item := range dictItems {
			fmt.Printf("  %d. %s: %s\n", i+1, item.DictCode, item.DictName)
		}
	}

	// 根据代码获取字典
	fmt.Println("\n4. 根据代码获取字典...")
	genderDict, err := client.GetDictByCode("gender")
	if err != nil {
		fmt.Printf("❌ 获取性别字典失败: %v\n", err)
	} else {
		fmt.Println("✅ 性别字典:")
		for _, item := range genderDict {
			fmt.Printf("  - %s: %s\n", item.DictValue, item.DictName)
		}
	}

	// 示例：创建新用户
	fmt.Println("\n5. 示例：创建新用户...")
	newUser := map[string]interface{}{
		"username":  "testuser",
		"realname":  "测试用户",
		"password":  "123456",
		"email":     "test@example.com",
		"phone":     "13800138000",
		"depart_id": "部门ID",
		"role_id":   "角色ID",
		"status":    1,
	}

	// 取消注释以下代码来实际创建用户
	/*
	if err := client.CreateUser(newUser); err != nil {
		fmt.Printf("❌ 创建用户失败: %v\n", err)
	} else {
		fmt.Println("✅ 用户创建成功")
	}
	*/

	fmt.Println("\n📋 API端点总结:")
	fmt.Println("  - POST   /admin/login           - 用户登录")
	fmt.Println("  - GET    /admin/system/membergetlist - 获取用户列表")
	fmt.Println("  - GET    /admin/system/dictgetlist   - 获取字典列表")
	fmt.Println("  - GET    /admin/dict             - 根据代码获取字典")
	fmt.Println("  - POST   /admin/system/memberadd - 创建用户")
	fmt.Println("  - POST   /admin/system/dictadd    - 创建字典项")
	fmt.Println("  - GET    /admin/getuserinfo      - 获取用户信息")
	fmt.Println("  - GET    /admin/getpermcode      - 获取权限代码")
	fmt.Println("  - GET    /admin/logout           - 退出登录")

	fmt.Println("\n✅ 高级API调用演示完成！")
	fmt.Println("💡 提示: 修改代码中的用户名、密码和URL以适应您的环境")
}