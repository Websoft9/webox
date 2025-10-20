package response

import (
	"api-service/internal/model"
	"time"
)

// UserResponse 用户响应结构
type UserResponse struct {
	ID          uint       `json:"id" example:"1"`
	Username    string     `json:"username" example:"john_doe"`
	Email       string     `json:"email" example:"john@example.com"`
	Nickname    string     `json:"nickname" example:"John Doe"`
	Avatar      string     `json:"avatar" example:"https://example.com/avatar.jpg"`
	Phone       string     `json:"phone" example:"+1234567890"`
	Gender      int        `json:"gender" example:"1"`
	Signature   string     `json:"signature" example:"This is my signature"`
	Status      int        `json:"status" example:"1"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" example:"2023-01-01T12:00:00Z"`
	LastLoginIP string     `json:"last_login_ip" example:"192.168.1.1"`
	Timezone    string     `json:"timezone" example:"Asia/Shanghai"`
	Language    string     `json:"language" example:"zh-CN"`
	CreatedAt   time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`

	// 关联数据
	Roles []RoleResponse `json:"roles,omitempty"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total" example:"100"`
}

// BuildUserResponse builds a user response from a user model

func BuildUserResponse(user *model.User) *UserResponse {
	resp := &UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Phone:       user.Phone,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		Signature:   user.Signature,
		Status:      user.Status,
		LastLoginAt: user.LastLoginAt,
		LastLoginIP: user.LastLoginIP,
		Timezone:    user.Timezone,
		Language:    user.Language,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	if len(user.Roles) > 0 {
		resp.Roles = make([]RoleResponse, len(user.Roles))
		for i := range user.Roles {
			role := &user.Roles[i]
			resp.Roles[i] = RoleResponse{
				ID:          role.ID,
				Name:        role.Name,
				Code:        role.Code,
				Description: role.Description,
				CreatedAt:   role.CreatedAt,
				UpdatedAt:   role.UpdatedAt,
			}
		}
	}

	return resp
}
