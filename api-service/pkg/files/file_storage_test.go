package files

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func setupTestStorage(t *testing.T) (*FileStorage, string) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_storage_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Reset singleton for testing
	instance = nil
	once = sync.Once{}

	// Initialize file storage
	fs := InitFileStorage()

	return fs, tempDir
}

func cleanupTestStorage(tempDir string) {
	os.RemoveAll(tempDir)
}

func TestInitFileStorage(t *testing.T) {
	fs, tempDir := setupTestStorage(t)
	defer cleanupTestStorage(tempDir)

	if fs == nil {
		t.Fatal("FileStorage instance should not be nil")
	}

	// Test singleton pattern
	fs2 := GetInstance()
	if fs != fs2 {
		t.Error("GetInstance should return the same instance")
	}
}

func TestSaveFile(t *testing.T) {
	fs, tempDir := setupTestStorage(t)
	defer cleanupTestStorage(tempDir)

	tests := []struct {
		name        string
		content     string
		storagePath string
		isOverwrite bool
		wantErr     bool
	}{
		{
			name:        "Save new file",
			content:     "test content",
			storagePath: filepath.Join(tempDir, "test.txt"),
			isOverwrite: false,
			wantErr:     false,
		},
		{
			name:        "Save file with subdirectory",
			content:     "test content",
			storagePath: filepath.Join(tempDir, "subdir/test.txt"),
			isOverwrite: false,
			wantErr:     false,
		},
		{
			name:        "Empty storage path",
			content:     "test content",
			storagePath: "",
			isOverwrite: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.content))
			fileName, err := fs.SaveFile(reader, tt.storagePath, tt.isOverwrite)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && fileName == "" {
				t.Error("SaveFile() should return file name")
			}

			// Verify file was created
			if !tt.wantErr {
				if _, err := os.Stat(tt.storagePath); os.IsNotExist(err) {
					t.Error("File was not created")
				}
			}
		})
	}
}

func TestSaveFileOverwrite(t *testing.T) {
	fs, tempDir := setupTestStorage(t)
	defer cleanupTestStorage(tempDir)

	storagePath := filepath.Join(tempDir, "overwrite_test.txt")

	// Save initial file
	reader1 := bytes.NewReader([]byte("initial content"))
	_, err := fs.SaveFile(reader1, storagePath, false)
	if err != nil {
		t.Fatalf("Failed to save initial file: %v", err)
	}

	// Try to save without overwrite (should fail)
	reader2 := bytes.NewReader([]byte("new content"))
	_, err = fs.SaveFile(reader2, storagePath, false)
	if err == nil {
		t.Error("SaveFile should fail when file exists and isOverwrite=false")
	}

	// Save with overwrite (should succeed)
	reader3 := bytes.NewReader([]byte("overwritten content"))
	_, err = fs.SaveFile(reader3, storagePath, true)
	if err != nil {
		t.Errorf("SaveFile with overwrite should succeed: %v", err)
	}

	// Verify content was overwritten
	data, err := fs.GetFile(storagePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(data) != "overwritten content" {
		t.Errorf("Expected 'overwritten content', got '%s'", string(data))
	}
}

func TestGetFile(t *testing.T) {
	fs, tempDir := setupTestStorage(t)
	defer cleanupTestStorage(tempDir)

	// Create test file
	testContent := "test file content"
	storagePath := filepath.Join(tempDir, "get_test.txt")
	reader := bytes.NewReader([]byte(testContent))
	_, err := fs.SaveFile(reader, storagePath, false)
	if err != nil {
		t.Fatalf("Failed to save test file: %v", err)
	}

	tests := []struct {
		name        string
		storagePath string
		wantContent string
		wantErr     bool
	}{
		{
			name:        "Get existing file",
			storagePath: storagePath,
			wantContent: testContent,
			wantErr:     false,
		},
		{
			name:        "Get non-existent file",
			storagePath: filepath.Join(tempDir, "nonexistent.txt"),
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "Empty storage path",
			storagePath: "",
			wantContent: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := fs.GetFile(tt.storagePath)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(data) != tt.wantContent {
				t.Errorf("GetFile() content = %s, want %s", string(data), tt.wantContent)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	fs, tempDir := setupTestStorage(t)
	defer cleanupTestStorage(tempDir)

	// Create test file
	storagePath := filepath.Join(tempDir, "delete_test.txt")
	reader := bytes.NewReader([]byte("test content"))
	_, err := fs.SaveFile(reader, storagePath, false)
	if err != nil {
		t.Fatalf("Failed to save test file: %v", err)
	}

	tests := []struct {
		name        string
		storagePath string
		isForce     bool
		wantDeleted bool
		wantErr     bool
	}{
		{
			name:        "Delete existing file",
			storagePath: storagePath,
			isForce:     false,
			wantDeleted: true,
			wantErr:     false,
		},
		{
			name:        "Delete non-existent file without force",
			storagePath: filepath.Join(tempDir, "nonexistent.txt"),
			isForce:     false,
			wantDeleted: false,
			wantErr:     true,
		},
		{
			name:        "Delete non-existent file with force",
			storagePath: filepath.Join(tempDir, "nonexistent2.txt"),
			isForce:     true,
			wantDeleted: false,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleted, err := fs.DeleteFile(tt.storagePath, tt.isForce)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if deleted != tt.wantDeleted {
				t.Errorf("DeleteFile() deleted = %v, want %v", deleted, tt.wantDeleted)
			}
		})
	}
}

func TestRenameFile(t *testing.T) {
	fs, tempDir := setupTestStorage(t)
	defer cleanupTestStorage(tempDir)

	// Create test file
	storagePath := filepath.Join(tempDir, "rename_dir")
	fileName := "original.txt"
	fullPath := filepath.Join(storagePath, fileName)
	reader := bytes.NewReader([]byte("test content"))
	_, err := fs.SaveFile(reader, fullPath, false)
	if err != nil {
		t.Fatalf("Failed to save test file: %v", err)
	}

	tests := []struct {
		name        string
		fileName    string
		newName     string
		storagePath string
		wantSuccess bool
		wantErr     bool
	}{
		{
			name:        "Rename existing file",
			fileName:    fileName,
			newName:     "renamed.txt",
			storagePath: storagePath,
			wantSuccess: true,
			wantErr:     false,
		},
		{
			name:        "Rename non-existent file",
			fileName:    "nonexistent.txt",
			newName:     "new.txt",
			storagePath: storagePath,
			wantSuccess: false,
			wantErr:     true,
		},
		{
			name:        "Empty file name",
			fileName:    "",
			newName:     "new.txt",
			storagePath: storagePath,
			wantSuccess: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success, err := fs.RenameFile(tt.fileName, tt.newName, tt.storagePath)

			if (err != nil) != tt.wantErr {
				t.Errorf("RenameFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if success != tt.wantSuccess {
				t.Errorf("RenameFile() success = %v, want %v", success, tt.wantSuccess)
			}

			// Verify file was renamed
			if tt.wantSuccess {
				newPath := filepath.Join(tt.storagePath, tt.newName)
				exists, err := fs.Exists(newPath)
				if err != nil {
					t.Errorf("Failed to check renamed file: %v", err)
				}
				if !exists {
					t.Error("Renamed file does not exist")
				}
			}
		})
	}
}

func TestExists(t *testing.T) {
	fs, tempDir := setupTestStorage(t)
	defer cleanupTestStorage(tempDir)

	// Create test file
	storagePath := filepath.Join(tempDir, "exists_test.txt")
	reader := bytes.NewReader([]byte("test content"))
	_, err := fs.SaveFile(reader, storagePath, false)
	if err != nil {
		t.Fatalf("Failed to save test file: %v", err)
	}

	// Create test directory
	dirPath := filepath.Join(tempDir, "test_dir")
	os.MkdirAll(dirPath, 0755)

	tests := []struct {
		name        string
		storagePath string
		wantExists  bool
		wantErr     bool
	}{
		{
			name:        "Existing file",
			storagePath: storagePath,
			wantExists:  true,
			wantErr:     false,
		},
		{
			name:        "Non-existent file",
			storagePath: filepath.Join(tempDir, "nonexistent.txt"),
			wantExists:  false,
			wantErr:     false,
		},
		{
			name:        "Directory path",
			storagePath: dirPath,
			wantExists:  false,
			wantErr:     true,
		},
		{
			name:        "Empty path",
			storagePath: "",
			wantExists:  false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := fs.Exists(tt.storagePath)

			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if exists != tt.wantExists {
				t.Errorf("Exists() = %v, want %v", exists, tt.wantExists)
			}
		})
	}
}
