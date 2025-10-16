package controller

import (
	response "api-service/internal/dto/common"
	"api-service/pkg/errors"
	"api-service/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidateRequest binds JSON request and validates it
func BindAndValidateRequest(ctx *gin.Context, req interface{}, validator *validator.Validate, log logger.Logger) bool {
	return bindAndValidate(ctx, req, validator, log, ctx.ShouldBindJSON, "Invalid request format", "Request validation failed")
}

// BindAndValidateQuery binds query parameters and validates them
func BindAndValidateQuery(ctx *gin.Context, req interface{}, validator *validator.Validate, log logger.Logger) bool {
	return bindAndValidate(ctx, req, validator, log, ctx.ShouldBindQuery, "Invalid query parameters", "Query validation failed")
}

// bindAndValidate is a helper function that binds and validates request/query parameters
func bindAndValidate(
	ctx *gin.Context,
	req interface{},
	validator *validator.Validate,
	log logger.Logger,
	bindFunc func(interface{}) error,
	bindErrMsg, validateErrMsg string,
) bool {
	// Bind parameters
	if err := bindFunc(req); err != nil {
		log.ErrorContext(ctx.Request.Context(), bindErrMsg, logger.ErrorField(err))
		response.WithErrorAndCode(ctx, errors.CodeInvalidParameterFormat, err)
		return false
	}

	// Validate parameters
	if err := validator.Struct(req); err != nil {
		log.ErrorContext(ctx.Request.Context(), validateErrMsg, logger.ErrorField(err))
		response.WithErrorAndCode(ctx, errors.CodeValidationFailed, err)
		return false
	}

	return true
}

// GetUserID extracts and validates user ID from context
func GetUserID(ctx *gin.Context) (uint, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx)
		return 0, false
	}
	return userID.(uint), true
}

// ParseIDParam parses ID parameter from URL
func ParseIDParam(ctx *gin.Context, paramName string) (uint, bool) {
	idStr := ctx.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.WithErrorAndCode(ctx, errors.CodeInvalidParameterFormat, err)
		return 0, false
	}
	return uint(id), true
}
