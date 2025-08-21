package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 测试用户管理API接口

const baseURL = "http://localhost:8080/api/v1"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 测试创建用户
func testCreateUser() {
	url := baseURL + "/users"
	payload := map[string]interface{}{
		"group_id":  1,
		"username":  "testuser",
		"email":     "test@example.com",
		"password":  "TestPass123",
		"nickname":  "测试用户",
		"phone":     "13800138000",
		"gender":    1,
		"signature": "这是一个测试用户",
		"timezone":  "Asia/Shanghai",
		"language":  "zh-CN",
		"status":    1,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("创建用户请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("创建用户响应: %+v\n", result)
}

// 测试获取用户列表
func testListUsers() {
	url := baseURL + "/users?page=1&page_size=10"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("获取用户列表请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("用户列表响应: %+v\n", result)
}

// 测试获取用户详情
func testGetUser() {
	url := baseURL + "/users/1"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("获取用户详情请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("用户详情响应: %+v\n", result)
}

// 测试更新用户
func testUpdateUser() {
	url := baseURL + "/users/1"
	payload := map[string]interface{}{
		"nickname":  "更新后的昵称",
		"signature": "更新后的签名",
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("更新用户请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("更新用户响应: %+v\n", result)
}

// 测试修改用户密码
func testChangePassword() {
	url := baseURL + "/users/1/password"
	payload := map[string]interface{}{
		"new_password":     "NewPass123",
		"confirm_password": "NewPass123",
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("修改密码请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("修改密码响应: %+v\n", result)
}

// 测试更新用户状态
func testUpdateUserStatus() {
	url := baseURL + "/users/1/status"
	payload := map[string]interface{}{
		"status": 0, // 禁用用户
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("更新用户状态请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("更新用户状态响应: %+v\n", result)
}

// 测试删除用户
func testDeleteUser() {
	url := baseURL + "/users/1"
	req, _ := http.NewRequest("DELETE", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("删除用户请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Printf("删除用户响应: %+v\n", result)
}

func main() {
	fmt.Println("=== 用户管理API测试 ===")
	fmt.Printf("测试时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// 注意：这些测试需要JWT认证和管理员权限
	// 在实际测试时需要先登录获取token，并在请求头中添加Authorization
	fmt.Println("注意：实际测试时需要添加JWT认证头")
	fmt.Println()

	fmt.Println("1. 测试创建用户")
	testCreateUser()
	fmt.Println()

	fmt.Println("2. 测试获取用户列表")
	testListUsers()
	fmt.Println()

	fmt.Println("3. 测试获取用户详情")
	testGetUser()
	fmt.Println()

	fmt.Println("4. 测试更新用户")
	testUpdateUser()
	fmt.Println()

	fmt.Println("5. 测试修改用户密码")
	testChangePassword()
	fmt.Println()

	fmt.Println("6. 测试更新用户状态")
	testUpdateUserStatus()
	fmt.Println()

	fmt.Println("7. 测试删除用户")
	testDeleteUser()
	fmt.Println()
}
