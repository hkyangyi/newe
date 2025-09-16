# Newe API 调用文档

本文档提供 Newe 框架的 API 调用示例和详细说明。

## 📋 快速开始

### 基础客户端使用

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	// 创建客户端
	client := NewNeweClient("http://localhost:8080")
	
	// 登录
	err := client.Login("admin", "admin123")
	if err != nil {
		log.Fatal("登录失败:", err)
	}
	
	// 获取用户信息
	userInfo, err := client.GetUserInfo()
	if err != nil {
		log.Fatal("获取用户信息失败:", err)
	}
	
	fmt.Printf("用户: %s (%s)\n", userInfo.Username, userInfo.RealName)
}
```

### 高级客户端使用

```go
// 创建高级客户端
client := NewAdvancedNeweClient("http://localhost:8080")

// 获取分页数据
users, total, err := client.GetUserList(1, 10)
dictItems, total, err := client.GetDictList("", 1, 20)
```

## 🔐 认证接口

### 登录接口
- **端点**: `POST /admin/login`
- **参数**:
  ```json
  {
    "username": "admin",
    "password": "admin123"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "result": {
      "token": "jwt_token_here",
      "user": { ... }
    },
    "success": "success"
  }
  ```

### 退出登录
- **端点**: `GET /admin/logout`
- **头部**: `Authorization: Bearer <token>`

## 👥 用户管理接口

### 获取用户列表
- **端点**: `GET /admin/system/membergetlist`
- **查询参数**:
  - `page`: 页码 (默认: 1)
  - `pageSize`: 每页数量 (默认: 10)
  - `username`: 用户名筛选
  - `realname`: 真实姓名筛选
  - `depart_id`: 部门ID筛选

### 创建用户
- **端点**: `POST /admin/system/memberadd`
- **参数**:
  ```json
  {
    "username": "testuser",
    "realname": "测试用户",
    "password": "123456",
    "email": "test@example.com",
    "phone": "13800138000",
    "depart_id": "部门ID",
    "role_id": "角色ID",
    "status": 1
  }
  ```

### 修改用户
- **端点**: `PUT /admin/system/memberedit`
- **参数**: 同创建用户，需要包含用户ID

### 删除用户
- **端点**: `DELETE /admin/system/memberdel`
- **参数**: `{"id": "用户ID"}`

## 📚 字典管理接口

### 获取字典列表
- **端点**: `GET /admin/system/dictgetlist`
- **查询参数**:
  - `page`: 页码
  - `pageSize`: 每页数量
  - `parent_id`: 父级ID (为空时获取顶级字典)

### 根据代码获取字典
- **端点**: `GET /admin/dict`
- **查询参数**: `dict_code=字典代码`

### 创建字典项
- **端点**: `POST /admin/system/dictadd`
- **参数**:
  ```json
  {
    "dict_code": "gender",
    "dict_name": "性别",
    "dict_value": "1",
    "dict_desc": "男性",
    "parent_id": "",
    "sort_no": 1,
    "status": 1
  }
  ```

## 🩺 系统监控接口

### 健康检查
- **端点**: `GET /health`
- **响应**:
  ```json
  {
    "status": "healthy",
    "timestamp": "2023-12-01T10:00:00Z",
    "version": "1.0.0"
  }
  ```

### 就绪检查
- **端点**: `GET /ready`

### 性能指标
- **端点**: `GET /metrics`

## 🔧 错误处理

### 响应格式
所有API都返回统一的响应格式：

```json
{
  "code": 200,        // 状态码
  "message": "成功",   // 消息
  "result": {...},    // 数据结果
  "success": "success" // 成功标识
}
```

### 常见错误码
- `200`: 成功
- `400`: 请求参数错误
- `401`: 未授权/Token过期
- `403`: 权限不足
- `404`: 资源不存在
- `500`: 服务器内部错误

## 🛡️ 安全说明

### Token 使用
1. 登录后获取 JWT Token
2. 在所有需要认证的请求头中包含:
   ```
   Authorization: Bearer <your_token>
   ```
3. Token 默认有效期为 72 小时

### 权限控制
- 管理员 (`admin`) 拥有所有权限
- 普通用户根据角色分配权限
- 数据权限根据部门组织结构控制

## 📝 示例代码

### 完整用户管理示例
```go
func userManagementDemo() {
	client := NewAdvancedNeweClient("http://localhost:8080")
	
	// 登录
	err := client.Login("admin", "admin123")
	if err != nil {
		log.Fatal(err)
	}
	
	// 获取用户列表
	users, total, err := client.GetUserList(1, 10)
	if err != nil {
		log.Fatal(err)
	}
	
	// 创建新用户
	newUser := map[string]interface{}{
		"username": "newuser",
		"realname": "新用户",
		"password": "password123",
		"email":    "new@example.com",
		"status":   1,
	}
	err = client.CreateUser(newUser)
	if err != nil {
		log.Fatal(err)
	}
}
```

### 字典数据操作示例
```go
func dictManagementDemo() {
	client := NewAdvancedNeweClient("http://localhost:8080")
	
	// 登录
	client.Login("admin", "admin123")
	
	// 获取性别字典
	genderDict, err := client.GetDictByCode("gender")
	if err != nil {
		log.Fatal(err)
	}
	
	// 创建新的字典项
	newDict := map[string]interface{}{
		"dict_code": "status",
		"dict_name": "状态",
		"dict_value": "1",
		"dict_desc": "激活",
		"parent_id": "",
		"sort_no":   1,
		"status":    1,
	}
	err = client.CreateDictItem(newDict)
	if err != nil {
		log.Fatal(err)
	}
}
```

## 🔗 相关资源

- [主项目 README](../README.md)
- [Docker 部署指南](../docker-compose.yml)
- [配置说明](../common/config/config.go)

## 💡 最佳实践

1. **错误处理**: 总是检查API调用的错误返回
2. **Token 管理**: 妥善保存和更新Token
3. **请求超时**: 设置合理的请求超时时间
4. **重试机制**: 对临时性错误实现重试逻辑
5. **日志记录**: 记录重要的API调用和错误

## 🆘 技术支持

如有问题请查看:
1. 检查服务是否正常运行
2. 验证配置参数是否正确
3. 查看日志文件获取详细错误信息
4. 在GitHub提交Issue