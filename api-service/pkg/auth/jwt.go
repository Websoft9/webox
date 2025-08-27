package auth

import (
	"api-service/internal/constants"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth represents JWT authentication handler
type JWTAuth struct {
	secretKey  string
	expireTime int
}

// Claims represents JWT claims structure
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// NewJWTAuth creates a new JWT authentication instance
func NewJWTAuth(secretKey string, expireTime int) *JWTAuth {
	return &JWTAuth{
		secretKey:  secretKey,
		expireTime: expireTime,
	}
}

// GenerateToken generates a JWT token for user authentication
func (j *JWTAuth) GenerateToken(userID uint) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(j.expireTime) * time.Second)

	claims := Claims{
		UserID:   userID,
		Username: "", // Can be queried from database or passed as parameter
		Role:     "", // Can be queried from database or passed as parameter
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

// GenerateTokenWithUserInfo generates JWT token with complete user information
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

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTAuth) GenerateTokenPair(userID uint, username, role string, refreshExpireTime int) (*TokenPair, error) {
	// Generate access token
	accessToken, accessExpiresAt, err := j.GenerateTokenWithUserInfo(userID, username, role)
	if err != nil {
		return nil, err
	}

	// Generate refresh token with longer expiration time
	refreshExpiresAt := time.Now().Add(time.Duration(refreshExpireTime) * time.Second)
	refreshClaims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiresAt,
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken validates and parses JWT token
func (j *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Additional validation: check if token is not expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken refreshes an access token using a refresh token
func (j *JWTAuth) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new token pair
	return j.GenerateTokenPair(claims.UserID, claims.Username, claims.Role, j.expireTime*constants.RefreshTokenMultiplier) // Refresh token lasts 24x longer
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[len(bearerPrefix):], nil
}

// Global JWT instance - needs to be set during initialization
var globalJWT *JWTAuth

// InitJWT initializes global JWT instance
func InitJWT(secretKey string, expireTime int) {
	globalJWT = NewJWTAuth(secretKey, expireTime)
}

// GenerateToken generates token using global JWT instance (backward compatibility)
func GenerateToken(userID uint, username string) (string, error) {
	if globalJWT == nil {
		// If not initialized, use default configuration
		globalJWT = NewJWTAuth("default-secret-key", constants.DefaultJWTExpireTime)
	}

	token, _, err := globalJWT.GenerateTokenWithUserInfo(userID, username, "user")
	return token, err
}

// ValidateToken validates token using global JWT instance (backward compatibility)
func ValidateToken(tokenString string) (*Claims, error) {
	if globalJWT == nil {
		globalJWT = NewJWTAuth("default-secret-key", constants.DefaultJWTExpireTime)
	}

	return globalJWT.ValidateToken(tokenString)
}

// GetGlobalJWT returns the global JWT instance
func GetGlobalJWT() *JWTAuth {
	return globalJWT
}

// UpdateJWTConfig updates the global JWT configuration
func UpdateJWTConfig(secretKey string, expireTime int) {
	globalJWT = NewJWTAuth(secretKey, expireTime)
}
