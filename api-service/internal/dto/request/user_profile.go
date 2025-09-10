package request

// UserProfileRequest 定义获取用户个人资料的请求参数
type UserProfileRequest struct {
	// UserID 指定要查询的用户ID，如果不提供，则获取当前登录用户的资料
	UserID uint `form:"userid" json:"userid" example:"1"`
}

// UserProfileUpdateRequest 更新用户个人资料的请求
type UserProfileUpdateRequest struct {
	UserID    uint    `json:"userid" form:"userid" binding:"required" example:"1"` // 要更新资料的用户ID
	Nickname  *string `json:"nickname" example:"John Doe"`
	Avatar    *string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Phone     *string `json:"phone" example:"+1234567890"`
	Gender    *int    `json:"gender" example:"1" binding:"omitempty,oneof=0 1 2"` // 0-未知，1-男，2-女
	Signature *string `json:"signature" example:"This is my signature"`
	Timezone  *string `json:"timezone" example:"Asia/Shanghai"`
	Language  *string `json:"language" example:"zh-CN" binding:"omitempty,len=5"`
}

// ProfileChangePasswordRequest 用户个人中心修改密码的请求
type ProfileChangePasswordRequest struct {
	UserID          uint   `json:"userid" form:"userid" binding:"required" example:"1"` // 要修改密码的用户ID
	OldPassword     string `json:"old_password" binding:"required" example:"oldpass123"`
	NewPassword     string `json:"new_password" binding:"required,min=6" example:"newpass123"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword" example:"newpass123"`
}
