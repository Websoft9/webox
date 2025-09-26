package auth

import (
	"api-service/internal/config"
	"api-service/pkg/errors"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// JWTLeewaySeconds defines the leeway time in seconds for JWT validation
	JWTLeewaySeconds = 5
)

// Global JWT instance - needs to be set during initialization
var globalJWT *JWTAuth

// JWTAuth represents JWT authentication handler
type JWTAuth struct {
	authConfig *config.AuthConfig
}

// Claims represents JWT claims structure
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     []uint `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTAuth creates a new JWT authentication instance
func NewJWTAuth(authConfig *config.AuthConfig) *JWTAuth {
	return &JWTAuth{
		authConfig: authConfig,
	}
}

// InitJWT initializes global JWT instance
func InitJWT(authConfig *config.AuthConfig) {
	globalJWT = NewJWTAuth(authConfig)
}

// GetGlobalJWT returns the global JWT instance
func GetGlobalJWT() *JWTAuth {
	return globalJWT
}

// UpdateJWTConfig updates the global JWT configuration
func UpdateJWTConfig(authConfig *config.AuthConfig) {
	globalJWT = NewJWTAuth(authConfig)
}

// HashToken hashes a token for secure storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ParseWithClaims parses JWT token string and validates it with the provided secret
func ParseWithClaims(tokenString, secret string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewAppError(errors.CodeInvalidParameterFormat)
		}
		return []byte(secret), nil
	}, jwt.WithLeeway(JWTLeewaySeconds*time.Second))
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.NewAppError(errors.CodeRequiredParameterMissing)
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.NewAppError(errors.CodeInvalidParameterFormat)
	}

	return authHeader[len(bearerPrefix):], nil
}

// GenerateTokenWithUserInfo generates JWT token with complete user information
func (j *JWTAuth) GenerateTokenWithUserInfo(userID uint, username string, roleIDs []uint) (string, time.Time, error) {
	// Get configured expiration time from auth config
	expiresInSeconds := j.GetTokenExpiresInSeconds()
	expiresAt := time.Now().Add(time.Duration(expiresInSeconds) * time.Second)

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     roleIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(j.GetTokenSigningMethod(), claims)
	tokenString, err := token.SignedString([]byte(j.authConfig.APIAuth.TokenAuth.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates and parses JWT token
func (j *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
	token, err := ParseWithClaims(tokenString, j.authConfig.APIAuth.TokenAuth.Secret)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.NewAppError(errors.CodeValidationFailed)
}

// RefreshToken refreshes an access token using a refresh token
func (j *JWTAuth) RefreshToken(tokenString string) (string, time.Time, error) {
	token, err := ParseWithClaims(tokenString, j.authConfig.APIAuth.TokenAuth.Secret)
	if err == nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return j.GenerateTokenWithUserInfo(claims.UserID, claims.Username, claims.Role)
		}
	}
	return "", time.Time{}, errors.NewAppError(errors.CodeValidationFailed)
}

// GetTokenExpiresInSeconds returns the configured token expiration time in seconds
func (j *JWTAuth) GetTokenExpiresInSeconds() int {
	// Get configured expiration time from auth config
	expiresInSeconds := j.authConfig.APIAuth.TokenAuth.ExpiresIn

	// Use configured expiration time with fallback to default if not configured
	if expiresInSeconds <= 0 {
		expiresInSeconds = 3600 // Default 1 hour if config is missing or invalid
	}
	return expiresInSeconds
}

// GetTokenSigningMethod returns the HMAC signing method based on configured algorithm
func (j *JWTAuth) GetTokenSigningMethod() *jwt.SigningMethodHMAC {
	algorithm := j.authConfig.APIAuth.TokenAuth.Algorithm
	switch strings.ToUpper(algorithm) {
	case "HS256":
		return jwt.SigningMethodHS256
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}
