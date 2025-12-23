package main

import (
	"fmt"
)

type migration struct {
	version    int
	statements []string
}

var migrations = []migration{
	{
		version: 1,
		statements: []string{
			`CREATE TABLE IF NOT EXISTS schema_migrations (
				version INTEGER PRIMARY KEY,
				applied_at TEXT DEFAULT CURRENT_TIMESTAMP
			);`,
			`CREATE TABLE IF NOT EXISTS manufacturers (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL UNIQUE
			);`,
			`CREATE TABLE IF NOT EXISTS storage_locations (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				friendly_name TEXT NOT NULL UNIQUE,
				note TEXT
			);`,
			`CREATE TABLE IF NOT EXISTS boxes (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				location_id INTEGER NOT NULL REFERENCES storage_locations(id) ON DELETE CASCADE,
				code TEXT NOT NULL UNIQUE,
				name TEXT,
				CHECK (length(trim(code)) > 0)
			);`,
			`CREATE TABLE IF NOT EXISTS bags (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				box_id INTEGER NOT NULL REFERENCES boxes(id) ON DELETE CASCADE,
				serial_no TEXT NOT NULL,
				UNIQUE (box_id, serial_no),
				CHECK (length(trim(serial_no)) > 0)
			);`,
			`CREATE TABLE IF NOT EXISTS sets (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				bag_id INTEGER NOT NULL UNIQUE REFERENCES bags(id) ON DELETE CASCADE,
				manufacturer_id INTEGER REFERENCES manufacturers(id) ON DELETE SET NULL,
				name TEXT NOT NULL,
				photo_path TEXT,
				photo_source TEXT,
				CHECK (length(trim(name)) > 0)
			);`,
			`CREATE TABLE IF NOT EXISTS elements (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				set_id INTEGER NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
				name TEXT NOT NULL,
				kind TEXT,
				CHECK (length(trim(name)) > 0),
				CHECK (kind IN ('stempel','stanze') OR kind IS NULL OR length(kind)=0)
			);`,
			`CREATE TABLE IF NOT EXISTS tags (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL UNIQUE,
				CHECK (length(trim(name)) > 0)
			);`,
			`CREATE TABLE IF NOT EXISTS set_tags (
				set_id INTEGER NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
				tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
				PRIMARY KEY (set_id, tag_id)
			);`,
			`CREATE INDEX IF NOT EXISTS idx_sets_name ON sets(name);`,
			`CREATE INDEX IF NOT EXISTS idx_boxes_code ON boxes(code);`,
			`CREATE INDEX IF NOT EXISTS idx_boxes_name ON boxes(name);`,
			`CREATE INDEX IF NOT EXISTS idx_bags_serial ON bags(serial_no);`,
			`CREATE INDEX IF NOT EXISTS idx_elements_name ON elements(name);`,
			`CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);`,
			`CREATE INDEX IF NOT EXISTS idx_set_tags_set_id ON set_tags(set_id);`,
		},
	},
	{
		version: 2,
		statements: []string{
			`CREATE TABLE IF NOT EXISTS types (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL UNIQUE,
				CHECK (length(trim(name)) > 0)
			);`,
			`ALTER TABLE sets ADD COLUMN type_id INTEGER REFERENCES types(id) ON DELETE SET NULL;`,
			`CREATE INDEX IF NOT EXISTS idx_types_name ON types(name);`,
		},
	},
	{
		version: 3,
		statements: []string{
			`ALTER TABLE storage_locations ADD COLUMN room TEXT;`,
			`ALTER TABLE storage_locations ADD COLUMN shelf TEXT;`,
			`ALTER TABLE storage_locations ADD COLUMN compartment TEXT;`,
		},
	},
}

func (a *App) runMigrations() error {
	if a.db == nil {
		return fmt.Errorf("database not initialised")
	}

	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TEXT DEFAULT CURRENT_TIMESTAMP
	);`); err != nil {
		return fmt.Errorf("init migration table: %w", err)
	}

	var current int
	if err = tx.QueryRow(`SELECT IFNULL(MAX(version), 0) FROM schema_migrations`).Scan(&current); err != nil {
		return fmt.Errorf("read migration version: %w", err)
	}

	for _, m := range migrations {
		if m.version <= current {
			continue
		}

		for _, stmt := range m.statements {
			if _, err = tx.Exec(stmt); err != nil {
				return fmt.Errorf("apply migration %d: %w", m.version, err)
			}
		}

		if _, err = tx.Exec(`INSERT INTO schema_migrations(version) VALUES (?)`, m.version); err != nil {
			return fmt.Errorf("record migration %d: %w", m.version, err)
		}
	}

	err = tx.Commit()
	return err
}
