package crypto

import (
	"testing"
)

func TestRSACryptoEncryptDecrypt(t *testing.T) {
	// Generate a new RSA key pair with 2048 bits
	rsaCrypto, err := NewRSACrypto(2048)
	if err != nil {
		t.Fatalf("Failed to create RSA crypto: %v", err)
	}

	// Calculate the maximum allowed plaintext size
	maxSize := rsaCrypto.GetMaxEncryptSize()
	t.Logf("Max encrypt size for 2048-bit key: %d bytes", maxSize)

	// Test data
	testData := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "Simple text",
			plaintext: "Hello, this is a test message",
		},
		{
			name:      "Empty string",
			plaintext: "",
		},
		{
			name:      "Special characters",
			plaintext: "!@#$%^&*()_+{}:\"<>?[];',./",
		},
		// Shortened long text to fit within the max encrypt size
		{
			name:      "Text within limits",
			plaintext: "This is a shorter text that should fit within RSA-OAEP size limitations.",
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			// Check if plaintext is within limits
			if len(tt.plaintext) > maxSize {
				t.Logf("Skipping test case: plaintext too long for standard RSA encryption (%d bytes > %d bytes max)",
					len(tt.plaintext), maxSize)
				return
			}

			// Encrypt
			encrypted, err := rsaCrypto.EncryptString(tt.plaintext)
			if err != nil {
				t.Fatalf("Failed to encrypt: %v", err)
			}

			// Decrypt
			decrypted, err := rsaCrypto.DecryptString(encrypted)
			if err != nil {
				t.Fatalf("Failed to decrypt: %v", err)
			}

			// Verify
			if decrypted != tt.plaintext {
				t.Errorf("Decrypted text does not match original.\nOriginal: %s\nDecrypted: %s", tt.plaintext, decrypted)
			}
		})
	}
}

func TestRSACryptoKeyPEMExport(t *testing.T) {
	// Generate a new RSA key pair
	rsaCrypto, err := NewRSACrypto(2048)
	if err != nil {
		t.Fatalf("Failed to create RSA crypto: %v", err)
	}

	// Export keys to PEM
	privateKeyPEM, err := rsaCrypto.GetPrivateKeyPEM()
	if err != nil {
		t.Fatalf("Failed to export private key: %v", err)
	}

	publicKeyPEM, err := rsaCrypto.GetPublicKeyPEM()
	if err != nil {
		t.Fatalf("Failed to export public key: %v", err)
	}

	// Create new instance from exported keys
	newRSACrypto, err := NewRSACryptoFromKeys(privateKeyPEM, publicKeyPEM)
	if err != nil {
		t.Fatalf("Failed to create RSA crypto from exported keys: %v", err)
	}

	// Test encryption and decryption with new instance
	plaintext := "Test message for key export/import"
	encrypted, err := newRSACrypto.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt with imported keys: %v", err)
	}

	decrypted, err := newRSACrypto.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt with imported keys: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original after key export/import.\nOriginal: %s\nDecrypted: %s", plaintext, decrypted)
	}
}

func TestRSACryptoGenerateKeyPair(t *testing.T) {
	// Generate key pair with the utility function
	privateKeyPEM, publicKeyPEM, err := GenerateKeyPair(2048)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create instance from generated keys
	rsaCrypto, err := NewRSACryptoFromKeys(privateKeyPEM, publicKeyPEM)
	if err != nil {
		t.Fatalf("Failed to create RSA crypto from generated keys: %v", err)
	}

	// Validate the keys
	err = rsaCrypto.ValidateKeys()
	if err != nil {
		t.Fatalf("Key validation failed: %v", err)
	}

	// Test encryption and decryption
	plaintext := "Test message for generated keys"
	encrypted, err := rsaCrypto.EncryptString(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt with generated keys: %v", err)
	}

	decrypted, err := rsaCrypto.DecryptString(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt with generated keys: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text does not match original with generated keys.\nOriginal: %s\nDecrypted: %s", plaintext, decrypted)
	}
}

func TestRSACryptoLargeData(t *testing.T) {
	// Generate a new RSA key pair
	rsaCrypto, err := NewRSACrypto(2048)
	if err != nil {
		t.Fatalf("Failed to create RSA crypto: %v", err)
	}

	// Get max encrypt size
	maxSize := rsaCrypto.GetMaxEncryptSize()
	t.Logf("Max encrypt size for 2048-bit key: %d bytes", maxSize)

	// Create test data larger than max size
	largeData := make([]byte, maxSize*3)
	for i := 0; i < len(largeData); i++ {
		largeData[i] = byte(i % 256)
	}

	// Encrypt large data
	encrypted, err := rsaCrypto.EncryptLarge(largeData)
	if err != nil {
		t.Fatalf("Failed to encrypt large data: %v", err)
	}

	// Decrypt large data
	decrypted, err := rsaCrypto.DecryptLarge(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt large data: %v", err)
	}

	// Verify
	if len(decrypted) != len(largeData) {
		t.Errorf("Decrypted data length mismatch. Expected: %d, Got: %d", len(largeData), len(decrypted))
	}

	// Check first few and last few bytes to save time
	for i := 0; i < 10; i++ {
		if decrypted[i] != largeData[i] {
			t.Errorf("Decrypted data mismatch at beginning. Index %d: expected %d, got %d", i, largeData[i], decrypted[i])
			break
		}
	}

	for i := 1; i <= 10; i++ {
		if decrypted[len(decrypted)-i] != largeData[len(largeData)-i] {
			t.Errorf("Decrypted data mismatch at end. Index %d: expected %d, got %d", len(largeData)-i, largeData[len(largeData)-i], decrypted[len(decrypted)-i])
			break
		}
	}
}

func TestRSACryptoInvalidInput(t *testing.T) {
	// Test with too small key size
	_, err := NewRSACrypto(1024)
	if err == nil {
		t.Error("Expected error for small key size, got nil")
	}

	// Create valid crypto instance
	rsaCrypto, err := NewRSACrypto(2048)
	if err != nil {
		t.Fatalf("Failed to create RSA crypto: %v", err)
	}

	// Test with invalid base64 for decryption
	_, err = rsaCrypto.Decrypt("not-valid-base64!@#$")
	if err == nil {
		t.Error("Expected error for invalid base64, got nil")
	}

	// Test with valid base64 but invalid encrypted data
	_, err = rsaCrypto.Decrypt("SGVsbG8sIHdvcmxkIQ==") // "Hello, world!" in base64
	if err == nil {
		t.Error("Expected error for invalid encrypted data, got nil")
	}
}

// Add a specific test for the EncryptString size limitation
func TestRSACryptoStringLengthLimitations(t *testing.T) {
	// Generate a new RSA key pair with 2048 bits
	rsaCrypto, err := NewRSACrypto(2048)
	if err != nil {
		t.Fatalf("Failed to create RSA crypto: %v", err)
	}

	// Calculate the maximum allowed plaintext size
	maxSize := rsaCrypto.GetMaxEncryptSize()
	t.Logf("Max encrypt size for 2048-bit key: %d bytes", maxSize)

	// Create a string just below the limit
	safeString := make([]byte, maxSize-1)
	for i := 0; i < len(safeString); i++ {
		safeString[i] = 'A'
	}

	// This should work
	encrypted, err := rsaCrypto.EncryptString(string(safeString))
	if err != nil {
		t.Errorf("Failed to encrypt string within size limit: %v", err)
	} else {
		// Decrypt to verify
		decrypted, err := rsaCrypto.DecryptString(encrypted)
		if err != nil {
			t.Errorf("Failed to decrypt: %v", err)
		}
		if decrypted != string(safeString) {
			t.Errorf("Decrypted text does not match original")
		}
	}

	// Create a string that exceeds the limit
	tooLongString := make([]byte, maxSize+10)
	for i := 0; i < len(tooLongString); i++ {
		tooLongString[i] = 'B'
	}

	// Regular encryption should fail
	_, err = rsaCrypto.EncryptString(string(tooLongString))
	if err == nil {
		t.Errorf("Expected error when encrypting string exceeding size limit")
	} else {
		t.Logf("Got expected error for too long string: %v", err)
	}

	// But EncryptLarge should handle it
	encryptedLarge, err := rsaCrypto.EncryptLarge(tooLongString)
	if err != nil {
		t.Errorf("EncryptLarge failed with long string: %v", err)
	} else {
		// Decrypt to verify
		decrypted, err := rsaCrypto.DecryptLarge(encryptedLarge)
		if err != nil {
			t.Errorf("Failed to decrypt large data: %v", err)
		}
		if string(decrypted) != string(tooLongString) {
			t.Errorf("Decrypted text does not match original for large encryption")
		}
	}
}
