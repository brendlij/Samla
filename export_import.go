package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ExportData exports the database and images folder to a zip file
func (a *App) ExportData() (string, error) {
	// Open save dialog
	defaultName := fmt.Sprintf("samla-backup-%s.zip", time.Now().Format("2006-01-02"))
	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Export Samla Data",
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{DisplayName: "Zip Files (*.zip)", Pattern: "*.zip"},
		},
	})
	if err != nil {
		return "", err
	}
	if savePath == "" {
		return "", nil // User cancelled
	}

	// Ensure .zip extension
	if !strings.HasSuffix(strings.ToLower(savePath), ".zip") {
		savePath += ".zip"
	}

	// Create zip file
	zipFile, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add database file
	dbPath := a.paths.DBPath
	if _, err := os.Stat(dbPath); err == nil {
		if err := addFileToZip(zipWriter, dbPath, "Data/samla.db"); err != nil {
			return "", fmt.Errorf("failed to add database to zip: %w", err)
		}
	}

	// Add images folder
	imagesDir := a.paths.ImagesDir
	if _, err := os.Stat(imagesDir); err == nil {
		err = filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(a.paths.BaseDir, path)
			if err != nil {
				return err
			}
			// Use forward slashes in zip
			relPath = strings.ReplaceAll(relPath, "\\", "/")

			return addFileToZip(zipWriter, path, relPath)
		})
		if err != nil {
			return "", fmt.Errorf("failed to add images to zip: %w", err)
		}
	}

	return savePath, nil
}

// ImportData imports data from a zip file
func (a *App) ImportData() (string, error) {
	// Open file dialog
	openPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Import Samla Data",
		Filters: []runtime.FileFilter{
			{DisplayName: "Zip Files (*.zip)", Pattern: "*.zip"},
		},
	})
	if err != nil {
		return "", err
	}
	if openPath == "" {
		return "", nil // User cancelled
	}

	// Close current database connection
	if a.db != nil {
		a.db.Close()
		a.db = nil
	}

	// Open zip file
	zipReader, err := zip.OpenReader(openPath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	// Extract files
	for _, file := range zipReader.File {
		destPath := filepath.Join(a.paths.BaseDir, file.Name)

		// Security check: ensure we don't write outside base dir
		if !strings.HasPrefix(destPath, a.paths.BaseDir) {
			continue
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(destPath, 0o755)
			continue
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return "", err
		}

		// Extract file
		if err := extractFileFromZip(file, destPath); err != nil {
			return "", fmt.Errorf("failed to extract %s: %w", file.Name, err)
		}
	}

	// Reopen database
	db, err := openDatabase(a.paths.DBPath)
	if err != nil {
		return "", fmt.Errorf("failed to reopen database: %w", err)
	}
	a.db = db

	return openPath, nil
}

// GetStats returns statistics about the data
func (a *App) GetStats() (map[string]int, error) {
	stats := make(map[string]int)

	// Count sets
	var setCount int
	if err := a.db.QueryRow(`SELECT COUNT(*) FROM sets`).Scan(&setCount); err == nil {
		stats["sets"] = setCount
	}

	// Count products
	var productCount int
	if err := a.db.QueryRow(`SELECT COUNT(*) FROM elements`).Scan(&productCount); err == nil {
		stats["products"] = productCount
	}

	// Count boxes
	var boxCount int
	if err := a.db.QueryRow(`SELECT COUNT(*) FROM boxes`).Scan(&boxCount); err == nil {
		stats["boxes"] = boxCount
	}

	// Count locations
	var locationCount int
	if err := a.db.QueryRow(`SELECT COUNT(*) FROM storage_locations`).Scan(&locationCount); err == nil {
		stats["locations"] = locationCount
	}

	// Count tags
	var tagCount int
	if err := a.db.QueryRow(`SELECT COUNT(*) FROM tags`).Scan(&tagCount); err == nil {
		stats["tags"] = tagCount
	}

	// Count images
	imageCount := 0
	filepath.Walk(a.paths.ImagesDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			imageCount++
		}
		return nil
	})
	stats["images"] = imageCount

	return stats, nil
}

func addFileToZip(zipWriter *zip.Writer, sourcePath, zipPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = zipPath
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func extractFileFromZip(file *zip.File, destPath string) error {
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	return err
}
