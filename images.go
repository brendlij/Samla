package main

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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
