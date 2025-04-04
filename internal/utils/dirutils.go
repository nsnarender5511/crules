package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo represents information about a file
type FileInfo struct {
	Path         string
	RelativePath string
	Size         int64
	ModTime      time.Time
	IsDir        bool
}

// ListDirectoryContents returns a list of all files in a directory recursively
func ListDirectoryContents(dir string) ([]FileInfo, error) {
	Debug("Listing directory contents | dir=" + dir)
	var files []FileInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == dir {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		files = append(files, FileInfo{
			Path:         path,
			RelativePath: relPath,
			Size:         info.Size(),
			ModTime:      info.ModTime(),
			IsDir:        info.IsDir(),
		})

		return nil
	})

	if err != nil {
		Error("Failed to list directory contents | dir=" + dir + ", error=" + err.Error())
		return nil, fmt.Errorf("failed to list directory contents: %v", err)
	}

	Debug("Listed directory contents | dir=" + dir + ", fileCount=" + fmt.Sprintf("%d", len(files)))
	return files, nil
}

// HasMDCFiles checks if a directory contains any .mdc files
func HasMDCFiles(dir string) (bool, error) {
	Debug("Checking for .mdc files | dir=" + dir)

	files, err := ListDirectoryContents(dir)
	if err != nil {
		Error("Failed to list directory contents | dir=" + dir + ", error=" + err.Error())
		return false, fmt.Errorf("failed to list directory contents: %v", err)
	}

	for _, file := range files {
		if !file.IsDir && strings.HasSuffix(file.RelativePath, ".mdc") {
			Debug("Found .mdc file | file=" + file.RelativePath)
			return true, nil
		}
	}

	Debug("No .mdc files found | dir=" + dir)
	return false, nil
}

// FormatFileSize returns a human-readable file size
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// CountFiles returns the number of files in a directory (non-recursive)
func CountFiles(dir string) (int, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() {
			count++
		}
	}
	return count, nil
}

// CountFilesByExt returns the number of files with a specific extension (non-recursive)
func CountFilesByExt(dir, ext string) (int, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ext {
			count++
		}
	}
	return count, nil
}

// CountFilesRecursive returns the total number of files in a directory and its subdirectories
func CountFilesRecursive(dir string) int {
	total := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			total++
		}
		return nil
	})
	return total
}
