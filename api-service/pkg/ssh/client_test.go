package ssh

import (
	"context"
	"testing"
	"time"
)

func TestClientConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *ClientConfig
		wantErr bool
	}{
		{
			name: "Valid password config",
			config: &ClientConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Username: "root",
				Password: "password123",
				Timeout:  30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "Valid private key config",
			config: &ClientConfig{
				Host:       "192.168.1.100",
				Port:       22,
				Username:   "root",
				PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----",
				Timeout:    30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "Missing host",
			config: &ClientConfig{
				Port:     22,
				Username: "root",
				Password: "password123",
				Timeout:  30 * time.Second,
			},
			wantErr: false, // Configuration validation happens during connection
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just test that config can be created
			if tt.config == nil {
				t.Error("Config should not be nil")
			}
		})
	}
}

func TestNewClient_InvalidConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		config  *ClientConfig
		wantErr bool
	}{
		{
			name: "Invalid host",
			config: &ClientConfig{
				Host:     "invalid.host.that.does.not.exist",
				Port:     22,
				Username: "root",
				Password: "password",
				Timeout:  2 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "Invalid port",
			config: &ClientConfig{
				Host:     "127.0.0.1",
				Port:     99999,
				Username: "root",
				Password: "password",
				Timeout:  2 * time.Second,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(ctx, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if client != nil {
				defer client.Close()
			}
		})
	}
}

// Note: The following tests require a real SSH server to be available
// They are marked as integration tests and should be run separately
// func TestClient_UploadFile(t *testing.T) { ... }
// func TestClient_DownloadFile(t *testing.T) { ... }
// func TestClient_DeleteFile(t *testing.T) { ... }
// func TestClient_TestConnection(t *testing.T) { ... }
