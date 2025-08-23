package auth

import (
	"api-service/internal/constants"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	secretKey  string
	expireTime int
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTAuth(secretKey string, expireTime int) *JWTAuth {
	return &JWTAuth{
		secretKey:  secretKey,
		expireTime: expireTime,
	}
}

func (j *JWTAuth) GenerateToken(userID uint) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(j.expireTime) * time.Second)

	claims := Claims{
		UserID:   userID,
		Username: "", // 可以从数据库查询或者从参数传入
		Role:     "", // 可以从数据库查询或者从参数传入
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// GenerateTokenWithUserInfo 生成包含用户信息的Token
func (j *JWTAuth) GenerateTokenWithUserInfo(userID uint, username, role string) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(j.expireTime) * time.Second)

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (j *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 全局JWT实例 - 需要在初始化时设置
var globalJWT *JWTAuth

// InitJWT 初始化全局JWT实例
func InitJWT(secretKey string, expireTime int) {
	globalJWT = NewJWTAuth(secretKey, expireTime)
}

// GenerateToken 生成Token（全局函数）
func GenerateToken(userID uint, username string) (string, error) {
	if globalJWT == nil {
		// 如果没有初始化，使用默认配置
		globalJWT = NewJWTAuth("default-secret-key", constants.DefaultJWTExpireTime)
	}

	token, _, err := globalJWT.GenerateTokenWithUserInfo(userID, username, "user")
	return token, err
}

// ValidateToken 验证Token（全局函数）
func ValidateToken(tokenString string) (*Claims, error) {
	if globalJWT == nil {
		globalJWT = NewJWTAuth("default-secret-key", constants.DefaultJWTExpireTime)
	}

	return globalJWT.ValidateToken(tokenString)
}
