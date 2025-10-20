package ssh

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Client represents an SSH client with SFTP support
type Client struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
	config     *ClientConfig
}

// ClientConfig contains SSH connection configuration
type ClientConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey string
	Timeout    time.Duration
}

// NewClient creates a new SSH client with SFTP support
func NewClient(ctx context.Context, config *ClientConfig) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second // nolint:mnd // Default SSH timeout
	}

	// Build SSH client configuration
	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // #nosec G106 -- Development mode, TODO: implement proper host key verification in production
		Timeout:         config.Timeout,
	}

	// Add authentication method
	if config.PrivateKey != "" {
		// Use private key authentication
		signer, err := ssh.ParsePrivateKey([]byte(config.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if config.Password != "" {
		// Use password authentication
		sshConfig.Auth = []ssh.AuthMethod{ssh.Password(config.Password)}
	} else {
		return nil, fmt.Errorf("either password or private key must be provided")
	}

	// Connect to SSH server
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		if closeErr := sshClient.Close(); closeErr != nil {
			return nil, fmt.Errorf("failed to create SFTP client: %w (close error: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("failed to create SFTP client: %w", err)
	}

	return &Client{
		sshClient:  sshClient,
		sftpClient: sftpClient,
		config:     config,
	}, nil
}

// Close closes the SSH and SFTP connections
func (c *Client) Close() error {
	var sftpErr, sshErr error

	if c.sftpClient != nil {
		sftpErr = c.sftpClient.Close()
	}

	if c.sshClient != nil {
		sshErr = c.sshClient.Close()
	}

	if sftpErr != nil {
		return fmt.Errorf("failed to close SFTP client: %w", sftpErr)
	}

	if sshErr != nil {
		return fmt.Errorf("failed to close SSH client: %w", sshErr)
	}

	return nil
}

// UploadFile uploads a file to the remote server
func (c *Client) UploadFile(localData []byte, remotePath string) error {
	// Create remote file
	remoteFile, err := c.sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer func() { _ = remoteFile.Close() }() // #nosec G104 -- defer close error can be safely ignored

	// Write data to remote file
	_, err = remoteFile.Write(localData)
	if err != nil {
		return fmt.Errorf("failed to write to remote file: %w", err)
	}

	return nil
}

// DownloadFile downloads a file from the remote server
func (c *Client) DownloadFile(remotePath string) ([]byte, error) {
	// Open remote file
	remoteFile, err := c.sftpClient.Open(remotePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open remote file: %w", err)
	}
	defer func() { _ = remoteFile.Close() }() // #nosec G104 -- defer close error can be safely ignored

	// Read file content
	data, err := io.ReadAll(remoteFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read remote file: %w", err)
	}

	return data, nil
}

// DeleteFile deletes a file from the remote server
func (c *Client) DeleteFile(remotePath string) error {
	err := c.sftpClient.Remove(remotePath)
	if err != nil {
		return fmt.Errorf("failed to delete remote file: %w", err)
	}
	return nil
}

// FileExists checks if a file exists on the remote server
func (c *Client) FileExists(remotePath string) (bool, error) {
	_, err := c.sftpClient.Stat(remotePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat remote file: %w", err)
	}
	return true, nil
}

// GetFileSize returns the size of a remote file
func (c *Client) GetFileSize(remotePath string) (int64, error) {
	stat, err := c.sftpClient.Stat(remotePath)
	if err != nil {
		return 0, fmt.Errorf("failed to stat remote file: %w", err)
	}
	return stat.Size(), nil
}

// TestConnection tests the SSH connection
func (c *Client) TestConnection() error {
	// Create a new session to test the connection
	session, err := c.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() { _ = session.Close() }() // #nosec G104 -- defer close error can be safely ignored

	// Run a simple command
	output, err := session.Output("echo 'connection_test'")
	if err != nil {
		return fmt.Errorf("failed to execute test command: %w", err)
	}

	if len(output) == 0 {
		return fmt.Errorf("empty response from test command")
	}

	return nil
}

// ExecuteCommand executes a command on the remote server
func (c *Client) ExecuteCommand(command string) (string, error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() { _ = session.Close() }() // #nosec G104 -- defer close error can be safely ignored

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return string(output), nil
}

// ReadFileContent reads the content of a remote file
func (c *Client) ReadFileContent(remotePath string) (string, error) {
	data, err := c.DownloadFile(remotePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFileContent writes content to a remote file
func (c *Client) WriteFileContent(remotePath, content string) error {
	return c.UploadFile([]byte(content), remotePath)
}
