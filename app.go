package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	_ "modernc.org/sqlite"
)

type App struct {
	ctx   context.Context
	db    *sql.DB
	paths AppPaths
}

type AppPaths struct {
	BaseDir   string `json:"baseDir"`
	DataDir   string `json:"dataDir"`
	ImagesDir string `json:"imagesDir"`
	DBPath    string `json:"dbPath"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	paths, err := resolveAppPaths()
	if err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("failed to resolve app paths: %v", err))
		return
	}
	a.paths = paths

	if err := ensureDirs(paths); err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("failed to prepare app folders: %v", err))
		return
	}

	db, err := openDatabase(paths.DBPath)
	if err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("failed to open database: %v", err))
		return
	}

	a.db = db

	if err := a.runMigrations(); err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("failed to run migrations: %v", err))
		return
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		_ = a.db.Close()
	}
}

func resolveAppPaths() (AppPaths, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return AppPaths{}, err
	}

	base := filepath.Join(configDir, "Samla")
	dataDir := filepath.Join(base, "Data")
	imagesDir := filepath.Join(base, "Images")
	dbPath := filepath.Join(dataDir, "samla.db")

	return AppPaths{
		BaseDir:   base,
		DataDir:   dataDir,
		ImagesDir: imagesDir,
		DBPath:    dbPath,
	}, nil
}

func ensureDirs(paths AppPaths) error {
	for _, dir := range []string{paths.BaseDir, paths.DataDir, paths.ImagesDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func openDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Enforce foreign keys and use WAL for better concurrent reads.
	if _, err := db.Exec(`PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;`); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply pragmas: %w", err)
	}

	return db, nil
}

// Utility: expose app folders to the UI.
func (a *App) GetAppPaths() AppPaths {
	return a.paths
}

// Utility: open the base folder in the platform file explorer.
func (a *App) OpenAppFolder() error {
	if a.paths.BaseDir == "" {
		return errors.New("paths not initialised")
	}
	return openFolder(a.paths.BaseDir)
}
