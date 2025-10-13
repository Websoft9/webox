package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	// DefaultDirPerm defines the default permission for created directories
	DefaultDirPerm os.FileMode = 0755
)

// FileStorage provides file storage management functionality
// It implements singleton pattern for global access
type FileStorage struct{}

var (
	instance *FileStorage
	once     sync.Once
	mu       sync.RWMutex
)

// InitFileStorage initializes the file storage instance
// This function should be called once during application startup
func InitFileStorage() *FileStorage {
	once.Do(func() {
		instance = &FileStorage{}
	})

	return instance
}

// GetInstance returns the singleton instance of FileStorage
// Panics if InitFileStorage has not been called
func GetInstance() *FileStorage {
	mu.RLock()
	defer mu.RUnlock()

	if instance == nil {
		panic("FileStorage not initialized. Call InitFileStorage first.")
	}

	return instance
}

// SaveFile saves an uploaded file to the specified storage path
// Parameters:
//   - file: io.Reader containing the file data
//   - storagePath: absolute path where the file should be stored
//   - isOverwrite: if true, overwrites existing file; if false, returns error if file exists
//
// Returns:
//   - string: the saved file name
//   - error: error if operation fails
func (fs *FileStorage) SaveFile(file io.Reader, storagePath string, isOverwrite bool) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file reader cannot be nil")
	}

	if storagePath == "" {
		return "", fmt.Errorf("storage path cannot be empty")
	}

	// Clean the path
	fullPath := filepath.Clean(storagePath)

	// Check if file exists
	if !isOverwrite {
		if _, err := os.Stat(fullPath); err == nil {
			return "", fmt.Errorf("file already exists: %s", storagePath)
		}
	}

	// Ensure parent directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, DefaultDirPerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		// Clean up on failure
		if removeErr := os.Remove(fullPath); removeErr != nil {
			return "", fmt.Errorf("failed to write file: %w, cleanup error: %v", err, removeErr)
		}
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filepath.Base(fullPath), nil
}

// GetFile retrieves a file from the specified storage path
// Parameters:
//   - fileStoragePath: absolute path of the file to retrieve
//
// Returns:
//   - []byte: file content as byte array
//   - error: error if operation fails
func (fs *FileStorage) GetFile(fileStoragePath string) ([]byte, error) {
	if fileStoragePath == "" {
		return nil, fmt.Errorf("file storage path cannot be empty")
	}

	// Clean the path
	fullPath := filepath.Clean(fileStoragePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", fileStoragePath)
	}

	// Read file content
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// DeleteFile deletes a file from the specified storage path
// Parameters:
//   - fileStoragePath: absolute path of the file to delete
//   - isForce: if true, ignores errors if file doesn't exist
//
// Returns:
//   - bool: true if file was deleted, false if file didn't exist (only when isForce=true)
//   - error: error if operation fails
func (fs *FileStorage) DeleteFile(fileStoragePath string, isForce bool) (bool, error) {
	if fileStoragePath == "" {
		return false, fmt.Errorf("file storage path cannot be empty")
	}

	// Clean the path
	fullPath := filepath.Clean(fileStoragePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if isForce {
			return false, nil
		}
		return false, fmt.Errorf("file not found: %s", fileStoragePath)
	}

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		return false, fmt.Errorf("failed to delete file: %w", err)
	}

	return true, nil
}

// RenameFile renames a file in the specified storage path
// Parameters:
//   - fileName: current file name
//   - newName: new file name
//   - storagePath: directory path where the file is located
//
// Returns:
//   - bool: true if file was renamed successfully
//   - error: error if operation fails
func (fs *FileStorage) RenameFile(fileName, newName, storagePath string) (bool, error) {
	if fileName == "" || newName == "" {
		return false, fmt.Errorf("file names cannot be empty")
	}

	if storagePath == "" {
		return false, fmt.Errorf("storage path cannot be empty")
	}

	// Construct full paths
	oldPath := filepath.Join(storagePath, fileName)
	newPath := filepath.Join(storagePath, newName)

	// Clean the paths
	oldPath = filepath.Clean(oldPath)
	newPath = filepath.Clean(newPath)

	// Check if source file exists
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return false, fmt.Errorf("source file not found: %s", fileName)
	}

	// Check if destination file already exists
	if _, err := os.Stat(newPath); err == nil {
		return false, fmt.Errorf("destination file already exists: %s", newName)
	}

	// Rename file
	if err := os.Rename(oldPath, newPath); err != nil {
		return false, fmt.Errorf("failed to rename file: %w", err)
	}

	return true, nil
}

// Exists checks if a file exists at the specified storage path
// Parameters:
//   - fileStoragePath: absolute path of the file to check
//
// Returns:
//   - bool: true if file exists, false otherwise
//   - error: error if operation fails (e.g., permission denied)
func (fs *FileStorage) Exists(fileStoragePath string) (bool, error) {
	if fileStoragePath == "" {
		return false, fmt.Errorf("file storage path cannot be empty")
	}

	// Clean the path
	fullPath := filepath.Clean(fileStoragePath)

	// Check if file exists
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	// Ensure it's a file, not a directory
	if info.IsDir() {
		return false, fmt.Errorf("path is a directory, not a file: %s", fileStoragePath)
	}

	return true, nil
}
