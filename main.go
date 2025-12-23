package main

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Samla",
		Width:     1200,
		Height:    800,
		MinWidth:  900,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileHandler(app),
		},
		BackgroundColour: &options.RGBA{R: 240, G: 253, B: 244, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// FileHandler serves local files from the app's data directory
type FileHandler struct {
	app *App
}

func NewFileHandler(app *App) *FileHandler {
	return &FileHandler{app: app}
}

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only handle /localfile/ requests
	if !strings.HasPrefix(r.URL.Path, "/localfile/") {
		return
	}

	// Get the relative path
	relPath := strings.TrimPrefix(r.URL.Path, "/localfile/")
	if relPath == "" {
		http.NotFound(w, r)
		return
	}

	// Build full path
	fullPath := filepath.Join(h.app.paths.BaseDir, relPath)

	// Security check: ensure path is within BaseDir
	cleanPath := filepath.Clean(fullPath)
	if !strings.HasPrefix(cleanPath, filepath.Clean(h.app.paths.BaseDir)) {
		http.NotFound(w, r)
		return
	}

	// Check file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Serve the file
	http.ServeFile(w, r, fullPath)
}
