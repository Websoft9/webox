package dto

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	GroupID   uint     `json:"group_id" binding:"required"`
	Username  string   `json:"username" binding:"required,min=3,max=32"`
	Email     string   `json:"email" binding:"required,email"`
	Password  string   `json:"password" binding:"required,min=8,max=32"`
	Nickname  string   `json:"nickname"`
	Phone     string   `json:"phone"`
	Gender    int8     `json:"gender"`
	Signature string   `json:"signature"`
	Timezone  string   `json:"timezone"`
	Language  string   `json:"language"`
	RoleIDs   []uint   `json:"role_ids"`
	Status    int8     `json:"status"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	GroupID   *uint    `json:"group_id"`
	Nickname  *string  `json:"nickname"`
	Email     *string  `json:"email"`
	Phone     *string  `json:"phone"`
	Gender    *int8    `json:"gender"`
	Signature *string  `json:"signature"`
	Timezone  *string  `json:"timezone"`
	Language  *string  `json:"language"`
	RoleIDs   []uint   `json:"role_ids"`
	Status    *int8    `json:"status"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=32"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// UserListQuery 用户列表查询参数
type UserListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	Keyword  string `form:"keyword"`
	Status   *int8  `form:"status"`
	RoleID   *uint  `form:"role_id"`
	Sort     string `form:"sort,default=created_at"`
	Order    string `form:"order,default=desc"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}
