#!/bin/bash

# API测试脚本
# 测试用户管理相关的API接口（不包括注册和登录）

BASE_URL="http://localhost:8080/api/v1"
# 这里需要手动设置一个有效的JWT Token用于测试
TOKEN="your_jwt_token_here"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印测试结果
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

# 打印分隔线
print_separator() {
    echo -e "${YELLOW}======================================${NC}"
    echo -e "${YELLOW} $1${NC}"
    echo -e "${YELLOW}======================================${NC}"
}

# 检查Token是否设置
if [ "$TOKEN" = "your_jwt_token_here" ]; then
    echo -e "${RED}错误：请先设置有效的JWT Token！${NC}"
    echo "请在脚本中将 TOKEN 变量设置为有效的JWT Token"
    exit 1
fi

echo -e "${YELLOW}使用Token: $TOKEN${NC}"
echo ""

# 1. 测试创建用户
print_separator "测试创建用户"
echo "POST $BASE_URL/users"
CREATE_USER_RESPONSE=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "group_id": 1,
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "NewPass123",
    "nickname": "New User",
    "phone": "13800138000",
    "gender": 1,
    "signature": "Test signature",
    "timezone": "Asia/Shanghai",
    "language": "zh-CN",
    "status": 1
  }')

HTTP_CODE=${CREATE_USER_RESPONSE: -3}
RESPONSE_BODY=${CREATE_USER_RESPONSE%???}

if [ "$HTTP_CODE" = "200" ]; then
    print_result 0 "创建用户成功"
    echo "Response: $RESPONSE_BODY"
    USER_ID=$(echo "$RESPONSE_BODY" | jq -r '.data.id' 2>/dev/null)
    echo "创建的用户ID: $USER_ID"
else
    print_result 1 "创建用户失败 (HTTP: $HTTP_CODE)"
    echo "Response: $RESPONSE_BODY"
fi

echo ""

# 2. 测试获取用户列表
print_separator "测试获取用户列表"
echo "GET $BASE_URL/users?page=1&page_size=10"
LIST_USERS_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/users?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN")

HTTP_CODE=${LIST_USERS_RESPONSE: -3}
RESPONSE_BODY=${LIST_USERS_RESPONSE%???}

if [ "$HTTP_CODE" = "200" ]; then
    print_result 0 "获取用户列表成功"
    echo "Response: $RESPONSE_BODY"
else
    print_result 1 "获取用户列表失败 (HTTP: $HTTP_CODE)"
    echo "Response: $RESPONSE_BODY"
fi

echo ""

# 3. 测试获取用户详情 (如果有用户ID)
if [ ! -z "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
    print_separator "测试获取用户详情"
    echo "GET $BASE_URL/users/$USER_ID"
    USER_DETAIL_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/users/$USER_ID" \
      -H "Authorization: Bearer $TOKEN")

    HTTP_CODE=${USER_DETAIL_RESPONSE: -3}
    RESPONSE_BODY=${USER_DETAIL_RESPONSE%???}

    if [ "$HTTP_CODE" = "200" ]; then
        print_result 0 "获取用户详情成功"
        echo "Response: $RESPONSE_BODY"
    else
        print_result 1 "获取用户详情失败 (HTTP: $HTTP_CODE)"
        echo "Response: $RESPONSE_BODY"
    fi

    echo ""

    # 4. 测试更新用户
    print_separator "测试更新用户"
    echo "PUT $BASE_URL/users/$USER_ID"
    UPDATE_USER_RESPONSE=$(curl -s -w "%{http_code}" -X PUT "$BASE_URL/users/$USER_ID" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "nickname": "Updated User",
        "signature": "Updated signature"
      }')

    HTTP_CODE=${UPDATE_USER_RESPONSE: -3}
    RESPONSE_BODY=${UPDATE_USER_RESPONSE%???}

    if [ "$HTTP_CODE" = "200" ]; then
        print_result 0 "更新用户成功"
        echo "Response: $RESPONSE_BODY"
    else
        print_result 1 "更新用户失败 (HTTP: $HTTP_CODE)"
        echo "Response: $RESPONSE_BODY"
    fi

    echo ""

    # 5. 测试修改密码
    print_separator "测试修改密码"
    echo "PUT $BASE_URL/users/$USER_ID/password"
    CHANGE_PASSWORD_RESPONSE=$(curl -s -w "%{http_code}" -X PUT "$BASE_URL/users/$USER_ID/password" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "old_password": "NewPass123",
        "new_password": "NewPass456",
        "confirm_password": "NewPass456"
      }')

    HTTP_CODE=${CHANGE_PASSWORD_RESPONSE: -3}
    RESPONSE_BODY=${CHANGE_PASSWORD_RESPONSE%???}

    if [ "$HTTP_CODE" = "200" ]; then
        print_result 0 "修改密码成功"
        echo "Response: $RESPONSE_BODY"
    else
        print_result 1 "修改密码失败 (HTTP: $HTTP_CODE)"
        echo "Response: $RESPONSE_BODY"
    fi

    echo ""
fi

# 6. 测试搜索用户
print_separator "测试搜索用户"
echo "GET $BASE_URL/users?keyword=test&page=1&page_size=5"
SEARCH_USERS_RESPONSE=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/users?keyword=test&page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN")

HTTP_CODE=${SEARCH_USERS_RESPONSE: -3}
RESPONSE_BODY=${SEARCH_USERS_RESPONSE%???}

if [ "$HTTP_CODE" = "200" ]; then
    print_result 0 "搜索用户成功"
    echo "Response: $RESPONSE_BODY"
else
    print_result 1 "搜索用户失败 (HTTP: $HTTP_CODE)"
    echo "Response: $RESPONSE_BODY"
fi

echo ""

print_separator "测试完成"
echo "所有API接口测试完成！"
