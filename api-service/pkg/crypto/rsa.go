package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

const (
	MinRSAKeySize      = 2048
	RSAKeySizeBits     = 8
	RSAOAEPHashLen     = 32 // SHA256 hash length in bytes
	RSAOAEPHashPadding = 2
)

// RSACrypto is a utility for RSA encryption and decryption
type RSACrypto struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSACrypto creates a new RSA encryption instance with specified key size
func NewRSACrypto(keySize int) (*RSACrypto, error) {
	if keySize < MinRSAKeySize {
		return nil, errors.New("RSA key size must be at least 2048 bits")
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	return &RSACrypto{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// NewRSACryptoFromPrivateKey creates an RSA crypto instance from PEM-encoded private key
func NewRSACryptoFromPrivateKey(privateKeyPEM string) (*RSACrypto, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}
	}

	return &RSACrypto{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// NewRSACryptoFromKeys creates an RSA crypto instance from PEM-encoded private and public keys
func NewRSACryptoFromKeys(privateKeyPEM, publicKeyPEM string) (*RSACrypto, error) {
	// Parse private key
	privateBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if privateBlock == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		// Try PKCS8 format
		key, pkcs8Err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
		if pkcs8Err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", pkcs8Err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an RSA private key")
		}
	}

	// Parse public key
	publicBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if publicBlock == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return &RSACrypto{
		privateKey: privateKey,
		publicKey:  rsaPublicKey,
	}, nil
}

// Encrypt encrypts data using RSA-OAEP with SHA256
func (r *RSACrypto) Encrypt(plaintext []byte) (string, error) {
	if r.publicKey == nil {
		return "", errors.New("public key is nil")
	}

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, r.publicKey, plaintext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Use base64 encoding
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// EncryptString encrypts a string using RSA-OAEP with SHA256
func (r *RSACrypto) EncryptString(plaintext string) (string, error) {
	return r.Encrypt([]byte(plaintext))
}

// Decrypt decrypts data using RSA-OAEP with SHA256
func (r *RSACrypto) Decrypt(ciphertextBase64 string) ([]byte, error) {
	if r.privateKey == nil {
		return nil, errors.New("private key is nil")
	}

	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privateKey, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// DecryptString decrypts a base64 encoded string to plaintext
func (r *RSACrypto) DecryptString(ciphertextBase64 string) (string, error) {
	plaintext, err := r.Decrypt(ciphertextBase64)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// GetPrivateKeyPEM returns the private key in PEM format
func (r *RSACrypto) GetPrivateKeyPEM() (string, error) {
	if r.privateKey == nil {
		return "", errors.New("private key is nil")
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(r.privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return string(privateKeyPEM), nil
}

// GetPublicKeyPEM returns the public key in PEM format
func (r *RSACrypto) GetPublicKeyPEM() (string, error) {
	if r.publicKey == nil {
		return "", errors.New("public key is nil")
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(r.publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM), nil
}

// GetKeySize returns the RSA key size in bits
func (r *RSACrypto) GetKeySize() int {
	if r.privateKey == nil {
		return 0
	}
	return r.privateKey.Size() * RSAKeySizeBits
}

// ValidateKeys verifies if the key pair is valid and working
func (r *RSACrypto) ValidateKeys() error {
	if r.privateKey == nil || r.publicKey == nil {
		return errors.New("keys are nil")
	}

	// Test encryption and decryption
	testData := []byte("test-validation-data")
	encrypted, err := r.Encrypt(testData)
	if err != nil {
		return fmt.Errorf("encryption failed during validation: %w", err)
	}

	decrypted, err := r.Decrypt(encrypted)
	if err != nil {
		return fmt.Errorf("decryption failed during validation: %w", err)
	}

	if !bytes.Equal(decrypted, testData) {
		return errors.New("key pair validation failed: decrypted data does not match original")
	}

	return nil
}

// GenerateKeyPair generates a new RSA key pair and returns them in PEM format
func GenerateKeyPair(keySize int) (privateKeyPEM, publicKeyPEM string, err error) {
	rsaCrypto, err := NewRSACrypto(keySize)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM, err = rsaCrypto.GetPrivateKeyPEM()
	if err != nil {
		return "", "", err
	}

	publicKeyPEM, err = rsaCrypto.GetPublicKeyPEM()
	if err != nil {
		return "", "", err
	}

	return privateKeyPEM, publicKeyPEM, nil
}

// GetMaxEncryptSize returns the maximum size of data that can be encrypted in a single operation
// RSA-OAEP with SHA256 maximum size = key size (bytes) - 2 * hash length (bytes) - 2
func (r *RSACrypto) GetMaxEncryptSize() int {
	if r.publicKey == nil {
		return 0
	}
	// SHA256 hash length is 32 bytes
	return r.publicKey.Size() - 2*RSAOAEPHashLen - RSAOAEPHashPadding
}

// EncryptLarge encrypts large data by splitting it into chunks
func (r *RSACrypto) EncryptLarge(plaintext []byte) (string, error) {
	maxSize := r.GetMaxEncryptSize()
	if maxSize <= 0 {
		return "", errors.New("invalid key size for encryption")
	}

	var result []byte
	for i := 0; i < len(plaintext); i += maxSize {
		end := i + maxSize
		if end > len(plaintext) {
			end = len(plaintext)
		}

		chunk := plaintext[i:end]
		encrypted, err := r.Encrypt(chunk)
		if err != nil {
			return "", fmt.Errorf("failed to encrypt chunk: %w", err)
		}

		// Convert each encrypted chunk's base64 back to bytes
		encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
		if err != nil {
			return "", fmt.Errorf("failed to decode encrypted chunk: %w", err)
		}

		result = append(result, encryptedBytes...)
	}

	return base64.StdEncoding.EncodeToString(result), nil
}

// DecryptLarge decrypts large data by processing it in chunks
func (r *RSACrypto) DecryptLarge(ciphertextBase64 string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	keySize := r.privateKey.Size()
	var result []byte

	for i := 0; i < len(ciphertext); i += keySize {
		end := i + keySize
		if end > len(ciphertext) {
			end = len(ciphertext)
		}

		chunk := ciphertext[i:end]
		chunkBase64 := base64.StdEncoding.EncodeToString(chunk)

		decrypted, err := r.Decrypt(chunkBase64)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt chunk: %w", err)
		}

		result = append(result, decrypted...)
	}

	return result, nil
}
