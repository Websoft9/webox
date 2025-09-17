package service

import (
	"api-service/internal/dto/request"
	"api-service/internal/dto/response"
	"api-service/internal/interface/repository"
	"api-service/internal/interface/service"
	"api-service/internal/model"
	"api-service/pkg/auth"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"context"
	"time"

	"gorm.io/gorm"
)

// User status constants
const (
	UserStatusActive   = 1
	UserStatusInactive = 0
)

type userService struct {
	userRepo repository.UserRepository
	logger   logger.Logger
}

func NewUserService(userRepo repository.UserRepository, logger logger.Logger) service.UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// ChangePassword changes the user's password
func (s *userService) ChangePassword(ctx context.Context, userID uint, req *request.UserChangePasswordRequest) error {
	s.logger.InfoContext(ctx, "Changing user password", logger.Uint("user_id", userID))

	// 1. Get user information
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeRecordNotFound, "User not found")
		}
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to get user information")
	}

	// 2. Verify old password
	if user.PasswordHash != auth.HashToken(req.OldPassword) {
		s.logger.WarnContext(ctx, "Old password verification failed", logger.Uint("user_id", userID))
		return errors.ErrInvalidCredentials
	}

	// 3. Encrypt new password
	hashedPassword := auth.HashToken(req.NewPassword)

	// 4. Update password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update password", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, "Failed to update password")
	}

	s.logger.InfoContext(ctx, "Password changed successfully", logger.Uint("user_id", userID))

	return nil
}

// ListUsers gets the user list
func (s *userService) ListUsers(ctx context.Context,
	req *request.UserListRequest) (*response.UserListResponse, int64, error) {
	s.logger.InfoContext(ctx, "Getting user list",
		logger.Int("page", req.Page),
		logger.Int("pageSize", req.PageSize))

	// Build filters
	filters := make(map[string]interface{})
	if req.Status != nil {
		filters["status"] = *req.Status
	}
	if req.Gender != nil {
		filters["gender"] = *req.Gender
	}
	if req.Language != nil {
		filters["language"] = *req.Language
	}
	if req.Keyword != nil && *req.Keyword != "" {
		filters["keyword"] = *req.Keyword
	}
	// Add role_id filter
	if req.RoleID != nil {
		filters["role_id"] = *req.RoleID
	}

	// Get user list
	users, total, err := s.userRepo.ListWithRelations(ctx, req.GetOffset(), req.GetPageSize(), filters)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get user list", logger.ErrorField(err))
		return nil, 0, errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to get user list")
	}

	// Convert to response format
	userList := make([]response.UserResponse, 0, len(users))
	for _, user := range users {
		userList = append(userList, *s.buildUserResponse(user))
	}

	return &response.UserListResponse{
		Users: userList,
		Total: total,
	}, total, nil
}

// GetUser gets user details
func (s *userService) GetUser(ctx context.Context, userID uint) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "Getting user details", logger.Uint("user_id", userID))

	user, err := s.userRepo.GetByIDWithRelations(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordQueryFailed, "User not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to get user")
	}

	return s.buildUserResponse(user), nil
}

// CreateUser creates a user
func (s *userService) CreateUser(ctx context.Context, currentUserID uint, req *request.UserCreateRequest) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "Creating user", logger.String("username", req.Username))

	// 1. Validate user creation
	if err := s.validateUserCreation(ctx, req); err != nil {
		return nil, err
	}

	// 2. Build user object from request
	user := s.createUserFromRequest(req)

	// 3. Save user to database (user.ID will be set after creation)
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, "Failed to create user")
	}

	// 4. Assign roles if role_ids provided
	if len(req.RoleIDs) > 0 {
		for _, roleID := range req.RoleIDs {
			userRole := &model.UserRole{
				UserID:    user.ID, // Use user.ID after creation
				RoleID:    roleID,
				GrantedBy: currentUserID, // Set as needed
				GrantedAt: time.Now(),
				Status:    1, // enabled
			}
			if err := s.userRepo.CreateUserRole(ctx, userRole); err != nil {
				s.logger.ErrorContext(ctx, "Failed to assign role to user", logger.Uint("user_id", user.ID), logger.Uint("role_id", roleID), logger.ErrorField(err))
				return nil, errors.WrapError(err, errors.CodeRecordCreateFailed, "Failed to assign role to user")
			}
		}
	}

	s.logger.InfoContext(ctx, "User created successfully", logger.Uint("user_id", user.ID))

	return s.buildUserResponse(user), nil
}

// UpdateUser updates user information and synchronizes roles.
func (s *userService) UpdateUser(ctx context.Context, currentUserID, userID uint, req *request.UserUpdateRequest) (*response.UserResponse, error) {
	s.logger.InfoContext(ctx, "Updating user", logger.Uint("user_id", userID))

	// 1. Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(errors.CodeRecordNotFound, "User not found")
		}
		s.logger.ErrorContext(ctx, "Failed to get user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to get user")
	}

	// 2. If updating email, check uniqueness
	if req.Email != nil && *req.Email != user.Email {
		if err := s.validateEmailUniqueness(ctx, *req.Email, userID); err != nil {
			return nil, err
		}
	}

	// 3. Update user info
	s.updateUserFromRequest(user, req)

	// 4. Save updates
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user", logger.ErrorField(err))
		return nil, errors.WrapError(err, errors.CodeRecordUpdateFailed, "Failed to update user")
	}

	// 5. Synchronize roles if RoleIDs provided
	if req.RoleIDs != nil {
		if err := s.syncUserRoles(ctx, currentUserID, userID, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	s.logger.InfoContext(ctx, "User updated successfully", logger.Uint("user_id", userID))
	return s.buildUserResponse(user), nil
}

// syncUserRoles synchronizes user roles according to the new role IDs.
func (s *userService) syncUserRoles(ctx context.Context, currentUserID, userID uint, newIDs []uint) error {
	newRoleIDs := make(map[uint]struct{}, len(newIDs))
	for _, id := range newIDs {
		newRoleIDs[id] = struct{}{}
	}

	currentRoleIDs, err := s.userRepo.GetRoleIDsByUserID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to get current user roles", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to get current user roles")
	}
	currentRoleIDSet := make(map[uint]struct{}, len(currentRoleIDs))
	for _, id := range currentRoleIDs {
		currentRoleIDSet[id] = struct{}{}
	}

	unionRoleIDs := make(map[uint]struct{})
	for id := range newRoleIDs {
		unionRoleIDs[id] = struct{}{}
	}
	for id := range currentRoleIDSet {
		unionRoleIDs[id] = struct{}{}
	}

	for roleID := range unionRoleIDs {
		_, inNew := newRoleIDs[roleID]
		_, inCurrent := currentRoleIDSet[roleID]
		if inNew && inCurrent {
			continue
		}
		if inNew && !inCurrent {
			userRole := &model.UserRole{
				UserID:    userID,
				RoleID:    roleID,
				GrantedBy: currentUserID,
				GrantedAt: time.Now(),
				Status:    1,
			}
			if err := s.userRepo.CreateUserRole(ctx, userRole); err != nil {
				s.logger.ErrorContext(ctx, "Failed to assign role to user", logger.Uint("user_id", userID), logger.Uint("role_id", roleID), logger.ErrorField(err))
				return errors.WrapError(err, errors.CodeRecordCreateFailed, "Failed to assign role to user")
			}
		}
		if !inNew && inCurrent {
			if err := s.userRepo.DeleteUserRole(ctx, userID, roleID); err != nil {
				s.logger.ErrorContext(ctx, "Failed to remove role from user", logger.Uint("user_id", userID), logger.Uint("role_id", roleID), logger.ErrorField(err))
				return errors.WrapError(err, errors.CodeRecordDeleteFailed, "Failed to remove role from user")
			}
		}
	}
	return nil
}

// UpdateUserStatus updates user status
func (s *userService) UpdateUserStatus(ctx context.Context, userID uint, req *request.UserUpdateStatusRequest) error {
	s.logger.InfoContext(ctx, "Updating user status", logger.Uint("user_id", userID), logger.Int("status", req.Status))

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeRecordNotFound, "User not found")
		}
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to get user information")
	}

	user.Status = req.Status
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update user status", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, "Failed to update user status")
	}

	s.logger.InfoContext(ctx, "User status updated successfully", logger.Uint("user_id", userID))
	return nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, userID uint) error {
	s.logger.InfoContext(ctx, "Deleting user", logger.Uint("user_id", userID))

	// 1. Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeRecordNotFound, "User not found")
		}
		s.logger.ErrorContext(ctx, "Failed to check user", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to check user")
	}

	// 2. Delete user
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete user", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordDeleteFailed, "Failed to delete user")
	}

	s.logger.InfoContext(ctx, "User deleted successfully", logger.Uint("user_id", userID))
	return nil
}

// UpdateUserPassword updates user password (admin operation)
func (s *userService) UpdateUserPassword(
	ctx context.Context, userID uint, req *request.UserPasswordUpdateRequest,
) error {
	s.logger.InfoContext(ctx, "Updating user password", logger.Uint("user_id", userID))

	// 1. Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(errors.CodeRecordNotFound, "User not found")
		}
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to get user information")
	}

	// 2. Encrypt new password
	hashedPassword := auth.HashToken(req.NewPassword)

	// 3. Update password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Failed to update password", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordUpdateFailed, "Failed to update password")
	}

	s.logger.InfoContext(ctx, "User password updated successfully", logger.Uint("user_id", userID))
	return nil
}

// ValidateUserAccess validates user access permissions
func (s *userService) ValidateUserAccess(ctx context.Context, userID uint, resource string) error {
	// TODO: Implement specific access validation logic
	return nil
}

// CheckUserQuota checks user quota
func (s *userService) CheckUserQuota(ctx context.Context, userID uint, resourceType string) error {
	// TODO: Implement specific quota check logic
	return nil
}

// Helper functions

// validateEmailUniqueness validates email uniqueness
func (s *userService) validateEmailUniqueness(ctx context.Context, email string, userID uint) error {
	exists, err := s.userRepo.ExistsByEmailExcludeID(ctx, email, userID)
	if err != nil {
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to check email uniqueness")
	}
	if exists {
		return errors.NewAppError(errors.CodeEmailAlreadyExists, "Email already exists")
	}
	return nil
}

// validateUserCreation validates user creation
func (s *userService) validateUserCreation(ctx context.Context, req *request.UserCreateRequest) error {
	// Check if username exists
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check username", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to check username")
	}
	if exists {
		return errors.ErrUserAlreadyExists
	}

	// Check if email exists
	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to check email", logger.ErrorField(err))
		return errors.WrapError(err, errors.CodeRecordQueryFailed, "Failed to check email")
	}
	if exists {
		return errors.NewAppError(errors.CodeEmailAlreadyExists, "Email already exists")
	}
	return nil
}

// createUserFromRequest builds user object from request
func (s *userService) createUserFromRequest(req *request.UserCreateRequest) *model.User {
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: auth.HashToken(req.Password),
		Status:       1, // default active
	}

	// Set optional fields
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Signature != nil {
		user.Signature = *req.Signature
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Timezone != nil {
		user.Timezone = *req.Timezone
	}
	if req.Language != nil {
		user.Language = *req.Language
	}

	return user
}

func (s *userService) updateUserFromRequest(user *model.User, req *request.UserUpdateRequest) {
	// Update user information
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.Signature != nil {
		user.Signature = *req.Signature
	}
	if req.Timezone != nil {
		user.Timezone = *req.Timezone
	}
	if req.Language != nil {
		user.Language = *req.Language
	}
}

// buildUserResponse builds user response
func (s *userService) buildUserResponse(user *model.User) *response.UserResponse {
	resp := &response.UserResponse{
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
		resp.Roles = make([]response.RoleResponse, len(user.Roles))
		for i := range user.Roles {
			role := &user.Roles[i]
			resp.Roles[i] = response.RoleResponse{
				ID:          role.ID,
				Name:        role.Name,
				Code:        role.Code,
				Description: role.Description,
			}
		}
	}

	return resp
}
