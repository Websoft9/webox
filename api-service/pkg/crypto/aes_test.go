package crypto

import (
	"strings"
	"testing"
)

func TestNewAESCrypto(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		wantErr   bool
	}{
		{
			name:      "valid secret key",
			secretKey: "my-secret-key-123",
			wantErr:   false,
		},
		{
			name:      "empty secret key",
			secretKey: "",
			wantErr:   true,
		},
		{
			name:      "long secret key",
			secretKey: "this-is-a-very-long-secret-key-that-should-work-perfectly-fine",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crypto, err := NewAESCrypto(tt.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESCrypto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && crypto == nil {
				t.Errorf("NewAESCrypto() returned nil crypto instance")
			}
		})
	}
}

func TestAESCrypto_EncryptDecrypt(t *testing.T) {
	crypto, err := NewAESCrypto("test-secret-key")
	if err != nil {
		t.Fatalf("Failed to create AESCrypto: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
		wantErr   bool
	}{
		{
			name:      "normal text",
			plaintext: "Hello, World!",
			wantErr:   false,
		},
		{
			name:      "unicode text",
			plaintext: "你好，世界！",
			wantErr:   false,
		},
		{
			name:      "long text",
			plaintext: strings.Repeat("A", 1000),
			wantErr:   false,
		},
		{
			name:      "empty text",
			plaintext: "",
			wantErr:   true,
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+{}|:<>?[]\\;'\".,/",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := crypto.Encrypt(tt.plaintext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify encrypted text is different from plaintext
			if encrypted == tt.plaintext {
				t.Errorf("Encrypted text should not be the same as plaintext")
			}

			// Test decryption
			decrypted, err := crypto.Decrypt(encrypted)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}

			if decrypted != tt.plaintext {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.plaintext)
			}
		})
	}
}

func TestAESCrypto_ValidateKey(t *testing.T) {
	crypto, err := NewAESCrypto("test-secret-key")
	if err != nil {
		t.Fatalf("Failed to create AESCrypto: %v", err)
	}

	err = crypto.ValidateKey()
	if err != nil {
		t.Errorf("ValidateKey() error = %v", err)
	}
}

func TestAESCrypto_RotateKey(t *testing.T) {
	crypto, err := NewAESCrypto("old-secret-key")
	if err != nil {
		t.Fatalf("Failed to create AESCrypto: %v", err)
	}

	// Encrypt with old key
	plaintext := "test-data"
	encrypted, err := crypto.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt with old key: %v", err)
	}

	// Rotate to new key
	newCrypto, err := crypto.RotateKey("new-secret-key")
	if err != nil {
		t.Fatalf("Failed to rotate key: %v", err)
	}

	// Old encrypted data should not be decryptable with new key
	_, err = newCrypto.Decrypt(encrypted)
	if err == nil {
		t.Errorf("Expected decryption to fail with rotated key")
	}

	// New encryption should work with new key
	newEncrypted, err := newCrypto.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt with new key: %v", err)
	}

	decrypted, err := newCrypto.Decrypt(newEncrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt with new key: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original")
	}
}

func TestAESCrypto_ErrorCases(t *testing.T) {
	crypto, err := NewAESCrypto("test-secret-key")
	if err != nil {
		t.Fatalf("Failed to create AESCrypto: %v", err)
	}

	// Test decryption with invalid base64
	_, err = crypto.Decrypt("invalid-base64!")
	if err == nil {
		t.Errorf("Expected error for invalid base64")
	}

	// Test decryption with valid base64 but invalid ciphertext
	_, err = crypto.Decrypt("YWJjZGVmZ2g=") // "abcdefgh" in base64, too short
	if err == nil {
		t.Errorf("Expected error for invalid ciphertext")
	}

	// Test decryption with tampered data
	validEncrypted, _ := crypto.Encrypt("test")
	// Tamper with the last character to ensure we actually modify the data
	tamperedData := validEncrypted[:len(validEncrypted)-1] + "X"
	_, err = crypto.Decrypt(tamperedData)
	if err == nil {
		t.Errorf("Expected error for tampered ciphertext")
	}
}
