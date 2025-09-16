package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// API演示客户端
type NeweClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// 登录请求结构
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 登录响应结构
type LoginResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
	Success string      `json:"success"`
}

// 用户信息结构
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int    `json:"status"`
}

// 创建新的API客户端
func NewNeweClient(baseURL string) *NeweClient {
	return &NeweClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// 用户登录
func (c *NeweClient) Login(username, password string) error {
	loginData := LoginRequest{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/admin/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return fmt.Errorf("JSON解析失败: %v", err)
	}

	if loginResp.Code != 200 {
		return fmt.Errorf("登录失败: %s", loginResp.Message)
	}

	// 从响应中提取token（根据实际响应结构调整）
	if result, ok := loginResp.Result.(map[string]interface{}); ok {
		if token, ok := result["token"].(string); ok {
			c.Token = token
			fmt.Printf("登录成功，Token: %s\n", token)
			return nil
		}
	}

	return fmt.Errorf("无法获取token")
}

// 获取用户信息
func (c *NeweClient) GetUserInfo() (*UserInfo, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/admin/getUserInfo", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var response struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Result  UserInfo `json:"result"`
		Success string   `json:"success"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("获取用户信息失败: %s", response.Message)
	}

	return &response.Result, nil
}

// 获取权限代码
func (c *NeweClient) GetPermCode() ([]string, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/admin/getPermCode", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var response struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Result  []string `json:"result"`
		Success string   `json:"success"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	if response.Code != 200 {
		return nil, fmt.Errorf("获取权限代码失败: %s", response.Message)
	}

	return response.Result, nil
}

// 健康检查
func (c *NeweClient) HealthCheck() error {
	req, err := http.NewRequest("GET", c.BaseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("健康检查失败: HTTP %d", resp.StatusCode)
	}

	fmt.Println("健康检查: 服务正常")
	return nil
}

// 退出登录
func (c *NeweClient) Logout() error {
	req, err := http.NewRequest("GET", c.BaseURL+"/admin/logout", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	var response LoginResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("JSON解析失败: %v", err)
	}

	if response.Code != 200 {
		return fmt.Errorf("退出登录失败: %s", response.Message)
	}

	fmt.Println("退出登录成功")
	c.Token = ""
	return nil
}

// 演示函数
func main() {
	fmt.Println("🚀 Newe API 调用演示")
	fmt.Println("====================")

	// 创建客户端
	client := NewNeweClient("http://localhost:8080")

	// 健康检查
	fmt.Println("\n1. 健康检查...")
	if err := client.HealthCheck(); err != nil {
		fmt.Printf("健康检查失败: %v\n", err)
		return
	}

	// 用户登录
	fmt.Println("\n2. 用户登录...")
	if err := client.Login("admin", "admin123"); err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}

	// 获取用户信息
	fmt.Println("\n3. 获取用户信息...")
	userInfo, err := client.GetUserInfo()
	if err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
		return
	}
	fmt.Printf("用户信息: %+v\n", userInfo)

	// 获取权限代码
	fmt.Println("\n4. 获取权限代码...")
	permCodes, err := client.GetPermCode()
	if err != nil {
		fmt.Printf("获取权限代码失败: %v\n", err)
		return
	}
	fmt.Printf("权限代码: %v\n", permCodes)

	// 退出登录
	fmt.Println("\n5. 退出登录...")
	if err := client.Logout(); err != nil {
		fmt.Printf("退出登录失败: %v\n", err)
	}

	fmt.Println("\n✅ API调用演示完成！")
}

// 示例：文件上传
func (c *NeweClient) UploadFile(filename string, fileData []byte) error {
	// 这里实现文件上传逻辑
	// 需要使用multipart/form-data格式
	return nil
}

// 示例：创建用户
func (c *NeweClient) CreateUser(userData map[string]interface{}) error {
	jsonData, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/admin/system/memberadd", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 处理响应...
	return nil
}