package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"
)

var (
	ErrInvalidKeySize    = errors.New("invalid key size: key must be 32 bytes for AES-256")
	ErrEmptyPlaintext    = errors.New("plaintext cannot be empty")
	ErrInvalidCiphertext = errors.New("invalid ciphertext format")
	ErrDecryptionFailed  = errors.New("decryption failed")
	ErrInvalidSecretKey  = errors.New("secret key cannot be empty")
)

// AESCrypto provides AES-256-GCM encryption and decryption functionality
type AESCrypto struct {
	secretKey []byte
}

// NewAESCrypto creates a new AESCrypto instance with the provided secret key
// The secret key will be hashed using SHA-256 to ensure it's exactly 32 bytes
func NewAESCrypto(secretKey string) (*AESCrypto, error) {
	if secretKey == "" {
		return nil, ErrInvalidSecretKey
	}

	// Hash the secret key using SHA-256 to get exactly 32 bytes for AES-256
	hash := sha256.Sum256([]byte(secretKey))

	return &AESCrypto{
		secretKey: hash[:],
	}, nil
}

// Encrypt encrypts the given plaintext using AES-256-GCM
// Returns base64-encoded string containing IV + encrypted data + auth tag
func (a *AESCrypto) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", ErrEmptyPlaintext
	}

	// Create AES cipher
	block, err := aes.NewCipher(a.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given base64-encoded ciphertext using AES-256-GCM
func (a *AESCrypto) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", ErrInvalidCiphertext
	}

	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(a.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Check minimum length
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	// Extract nonce and ciphertext
	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
	}

	return string(plaintext), nil
}

// RotateKey creates a new AESCrypto instance with a new secret key
// This can be used for key rotation scenarios
func (a *AESCrypto) RotateKey(newSecretKey string) (*AESCrypto, error) {
	return NewAESCrypto(newSecretKey)
}

// ValidateKey validates if the provided key can decrypt a test value
// This is useful for verifying key correctness during system startup
func (a *AESCrypto) ValidateKey() error {
	testValue := "test_encryption_validation"

	encrypted, err := a.Encrypt(testValue)
	if err != nil {
		return fmt.Errorf("validation failed during encryption: %w", err)
	}

	decrypted, err := a.Decrypt(encrypted)
	if err != nil {
		return fmt.Errorf("validation failed during decryption: %w", err)
	}

	if decrypted != testValue {
		return fmt.Errorf("validation failed: decrypted value does not match original")
	}

	return nil
}

// Singleton pattern for default crypto instance
var (
	defaultInstance *AESCrypto
	mu              sync.RWMutex
)

// InitDefaultCrypto initializes the default crypto instance with the provided key
func InitDefaultCrypto(encryptionKey string) (*AESCrypto, error) {
	crypto, err := NewAESCrypto(encryptionKey)
	if err != nil {
		return nil, err
	}

	mu.Lock()
	defer mu.Unlock()
	defaultInstance = crypto

	return defaultInstance, nil
}

// GetDefaultCrypto returns the default AES crypto instance (singleton)
func GetDefaultCrypto() *AESCrypto {
	mu.RLock()
	defer mu.RUnlock()

	if defaultInstance == nil {
		panic("Default crypto instance not initialized. Call InitDefaultCrypto first.")
	}

	return defaultInstance
}
