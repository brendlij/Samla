package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ChooseImageFile opens a file dialog to select an image
func (a *App) ChooseImageFile() (string, error) {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Bild auswählen / Choose Image",
		Filters: []runtime.FileFilter{
			{DisplayName: "Images", Pattern: "*.png;*.jpg;*.jpeg;*.gif;*.bmp;*.webp"},
		},
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (a *App) AttachImageFromFile(setID int64, filePath string) (string, error) {
	if setID <= 0 {
		return "", errors.New("set is required")
	}
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return "", errors.New("file path is empty")
	}

	src, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		ext = ".png"
	}
	fileName := fmt.Sprintf("%s%s", uuid.NewString(), ext)
	relPath := filepath.ToSlash(filepath.Join("Images", fileName))
	destPath := filepath.Join(a.paths.BaseDir, relPath)

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return "", err
	}
	_ = dst.Close()

	if err := a.setImagePath(setID, relPath, "file"); err != nil {
		return "", err
	}
	return relPath, nil
}

func (a *App) AttachImageFromURL(setID int64, rawURL string) (string, error) {
	if setID <= 0 {
		return "", errors.New("set is required")
	}
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", errors.New("url is required")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(parsed.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to download image (status %d)", resp.StatusCode)
	}

	ext := strings.ToLower(filepath.Ext(parsed.Path))
	if ext == "" {
		contentType := resp.Header.Get("Content-Type")
		switch {
		case strings.Contains(contentType, "png"):
			ext = ".png"
		case strings.Contains(contentType, "jpeg"), strings.Contains(contentType, "jpg"):
			ext = ".jpg"
		case strings.Contains(contentType, "gif"):
			ext = ".gif"
		default:
			ext = ".png"
		}
	}

	fileName := fmt.Sprintf("%s%s", uuid.NewString(), ext)
	relPath := filepath.ToSlash(filepath.Join("Images", fileName))
	destPath := filepath.Join(a.paths.BaseDir, relPath)

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return "", err
	}
	_ = out.Close()

	if err := a.setImagePath(setID, relPath, "url"); err != nil {
		return "", err
	}
	return relPath, nil
}

func (a *App) SaveCroppedImage(setID int64, base64Data string, ext string) (string, error) {
	if setID <= 0 {
		return "", errors.New("set is required")
	}
	if ext == "" {
		ext = ".png"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	data := strings.TrimSpace(base64Data)
	if strings.Contains(data, ",") {
		parts := strings.SplitN(data, ",", 2)
		data = parts[1]
	}

	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("unable to decode image: %w", err)
	}

	fileName := fmt.Sprintf("%s%s", uuid.NewString(), ext)
	relPath := filepath.ToSlash(filepath.Join("Images", fileName))
	destPath := filepath.Join(a.paths.BaseDir, relPath)

	if err := os.WriteFile(destPath, buf, 0o644); err != nil {
		return "", err
	}

	if err := a.setImagePath(setID, relPath, "cropped"); err != nil {
		return "", err
	}
	return relPath, nil
}

func (a *App) setImagePath(setID int64, relPath, source string) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var oldPath sql.NullString
	if err = tx.QueryRow(`SELECT photo_path FROM sets WHERE id = ?`, setID).Scan(&oldPath); err != nil {
		return err
	}

	if _, err = tx.Exec(`UPDATE sets SET photo_path = ?, photo_source = ? WHERE id = ?`, relPath, source, setID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	if oldPath.Valid && oldPath.String != "" && oldPath.String != relPath {
		_ = deleteLocalImage(a.paths.ImagesDir, oldPath.String)
	}
	return nil
}

func deleteLocalImage(imagesDir, relPath string) error {
	if relPath == "" {
		return nil
	}
	target := relPath
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(imagesDir), relPath)
	}
	// Basic safety: only delete inside Images directory.
	if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(imagesDir)) {
		return fmt.Errorf("refusing to delete outside images directory")
	}
	err := os.Remove(target)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// Log helper for backend notifications from image workflows.
func (a *App) logInfo(msg string) {
	if a.ctx != nil {
		runtime.LogInfo(a.ctx, msg)
	}
}

// ReadFileAsBase64 reads a file and returns it as a base64 data URL
func (a *App) ReadFileAsBase64(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	mimeType := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	case ".bmp":
		mimeType = "image/bmp"
	}

	base64Data := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data), nil
}

// GetImageAsBase64 loads an image by its relative path and returns base64
func (a *App) GetImageAsBase64(relPath string) (string, error) {
	if relPath == "" {
		return "", nil
	}
	
	var fullPath string
	if filepath.IsAbs(relPath) {
		fullPath = relPath
	} else {
		fullPath = filepath.Join(a.paths.BaseDir, relPath)
	}
	
	return a.ReadFileAsBase64(fullPath)
}

// Expose to frontend: return an accessible file URL for an image path.
func (a *App) ResolveImagePath(relPath string) string {
	if relPath == "" {
		return ""
	}
	if filepath.IsAbs(relPath) {
		return "file:///" + filepath.ToSlash(relPath)
	}
	full := filepath.Join(a.paths.BaseDir, relPath)
	return "file:///" + filepath.ToSlash(full)
}

// ScanImage uses Windows WIA (Windows Image Acquisition) to scan from a connected scanner
// Returns the path to the scanned image file
func (a *App) ScanImage() (string, error) {
	// Create temp file path for scanned image
	fileName := fmt.Sprintf("scan_%s.png", uuid.NewString())
	tempPath := filepath.Join(a.paths.ImagesDir, fileName)

	// PowerShell script to use WIA for scanning
	psScript := fmt.Sprintf(`
Add-Type -AssemblyName System.Drawing

# Create WIA DeviceManager
$deviceManager = New-Object -ComObject WIA.DeviceManager

# Get first scanner device
$device = $null
foreach ($d in $deviceManager.DeviceInfos) {
    if ($d.Type -eq 1) {  # Scanner type
        $device = $d.Connect()
        break
    }
}

if ($device -eq $null) {
    Write-Error "No scanner found"
    exit 1
}

# Get first item (scanner bed)
$item = $device.Items.Item(1)

# Configure for color scan at 200 DPI
$item.Properties.Item("6146").Value = 1  # Color intent
$item.Properties.Item("6147").Value = 200  # Horizontal DPI
$item.Properties.Item("6148").Value = 200  # Vertical DPI

# Scan
$imageFile = $item.Transfer("{B96B3CAE-0728-11D3-9D7B-0000F81EF32E}")  # PNG format

# Save to file
$imageFile.SaveFile("%s")

Write-Output "OK"
`, strings.ReplaceAll(tempPath, "\\", "\\\\"))

	// Execute PowerShell
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		if strings.Contains(errMsg, "No scanner found") {
			return "", errors.New("no scanner found - please connect a scanner via USB")
		}
		return "", fmt.Errorf("scan failed: %s", errMsg)
	}

	// Check if file was created
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		return "", errors.New("scan failed - no image created")
	}

	// Return relative path
	relPath := filepath.ToSlash(filepath.Join("Images", fileName))
	return relPath, nil
}

// ScanResult contains both the base64 data for preview and the relative path for storage
type ScanResult struct {
	Base64Data string `json:"base64Data"`
	RelPath    string `json:"relPath"`
}

// ScanImageToBase64 scans and returns the image as base64 for preview/cropping plus the path
func (a *App) ScanImageToBase64() (*ScanResult, error) {
	relPath, err := a.ScanImage()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(a.paths.BaseDir, relPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	// Return as data URL
	base64Data := base64.StdEncoding.EncodeToString(data)
	return &ScanResult{
		Base64Data: fmt.Sprintf("data:image/png;base64,%s", base64Data),
		RelPath:    relPath,
	}, nil
}

// AttachScannedImage attaches a previously scanned image to a set
func (a *App) AttachScannedImage(setID int64, scannedPath string) (string, error) {
	if setID <= 0 {
		return "", errors.New("set is required")
	}
	if scannedPath == "" {
		return "", errors.New("scanned path is required")
	}

	// The scannedPath should already be a relative path like "Images/scan_xxx.png"
	if err := a.setImagePath(setID, scannedPath, "scan"); err != nil {
		return "", err
	}
	return scannedPath, nil
}

// RemoveImage removes the image from a set and deletes the file
func (a *App) RemoveImage(setID int64) error {
	if setID <= 0 {
		return errors.New("set is required")
	}

	// Get current image path
	var oldPath sql.NullString
	err := a.db.QueryRow(`SELECT photo_path FROM sets WHERE id = ?`, setID).Scan(&oldPath)
	if err != nil {
		return err
	}

	// Clear the image path in DB
	_, err = a.db.Exec(`UPDATE sets SET photo_path = NULL, photo_source = NULL WHERE id = ?`, setID)
	if err != nil {
		return err
	}

	// Delete the file if it exists
	if oldPath.Valid && oldPath.String != "" {
		_ = deleteLocalImage(a.paths.ImagesDir, oldPath.String)
	}

	return nil
}
