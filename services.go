package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func normalizeName(val string) string {
	return strings.TrimSpace(val)
}

func normalizeLower(val string) string {
	return strings.ToLower(strings.TrimSpace(val))
}

// Locations
func (a *App) ListLocations() ([]StorageLocation, error) {
	rows, err := a.db.Query(`SELECT id, friendly_name, IFNULL(room,''), IFNULL(shelf,''), IFNULL(compartment,''), IFNULL(note,'') FROM storage_locations ORDER BY friendly_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []StorageLocation
	for rows.Next() {
		var loc StorageLocation
		if err := rows.Scan(&loc.ID, &loc.FriendlyName, &loc.Room, &loc.Shelf, &loc.Compartment, &loc.Note); err != nil {
			return nil, err
		}
		list = append(list, loc)
	}
	return list, rows.Err()
}

func (a *App) CreateLocation(name, room, shelf, compartment, note string) (int64, error) {
	name = normalizeName(name)
	if name == "" {
		return 0, errors.New("friendly name is required")
	}

	res, err := a.db.Exec(`INSERT INTO storage_locations(friendly_name, room, shelf, compartment, note) VALUES (?, ?, ?, ?, ?)`,
		name, strings.TrimSpace(room), strings.TrimSpace(shelf), strings.TrimSpace(compartment), strings.TrimSpace(note))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) UpdateLocation(id int64, name, room, shelf, compartment, note string) error {
	name = normalizeName(name)
	if name == "" {
		return errors.New("friendly name is required")
	}
	_, err := a.db.Exec(`UPDATE storage_locations SET friendly_name = ?, room = ?, shelf = ?, compartment = ?, note = ? WHERE id = ?`,
		name, strings.TrimSpace(room), strings.TrimSpace(shelf), strings.TrimSpace(compartment), strings.TrimSpace(note), id)
	return err
}

func (a *App) DeleteLocation(id int64) error {
	_, err := a.db.Exec(`DELETE FROM storage_locations WHERE id = ?`, id)
	return err
}

// Boxes
func (a *App) ListBoxes(locationID int64) ([]Box, error) {
	var rows *sql.Rows
	var err error
	if locationID > 0 {
		rows, err = a.db.Query(`SELECT id, location_id, code, IFNULL(name,'') FROM boxes WHERE location_id = ? ORDER BY code`, locationID)
	} else {
		rows, err = a.db.Query(`SELECT id, location_id, code, IFNULL(name,'') FROM boxes ORDER BY code`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Box
	for rows.Next() {
		var b Box
		if err := rows.Scan(&b.ID, &b.LocationID, &b.Code, &b.Name); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (a *App) CreateBox(locationID int64, code, name string) (int64, error) {
	code = normalizeName(code)
	if code == "" {
		return 0, errors.New("code is required")
	}
	if locationID <= 0 {
		return 0, errors.New("location is required")
	}
	res, err := a.db.Exec(`INSERT INTO boxes(location_id, code, name) VALUES (?, ?, ?)`, locationID, code, strings.TrimSpace(name))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) UpdateBox(id int64, locationID int64, code, name string) error {
	code = normalizeName(code)
	if code == "" {
		return errors.New("code is required")
	}
	if locationID <= 0 {
		return errors.New("location is required")
	}
	_, err := a.db.Exec(`UPDATE boxes SET location_id = ?, code = ?, name = ? WHERE id = ?`, locationID, code, strings.TrimSpace(name), id)
	return err
}

func (a *App) DeleteBox(id int64) error {
	_, err := a.db.Exec(`DELETE FROM boxes WHERE id = ?`, id)
	return err
}

// Manufacturers
func (a *App) ListManufacturers() ([]Manufacturer, error) {
	rows, err := a.db.Query(`SELECT id, name FROM manufacturers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Manufacturer
	for rows.Next() {
		var m Manufacturer
		if err := rows.Scan(&m.ID, &m.Name); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (a *App) CreateManufacturerIfMissing(name string) (int64, error) {
	name = normalizeName(name)
	if name == "" {
		return 0, nil
	}

	var id int64
	err := a.db.QueryRow(`SELECT id FROM manufacturers WHERE LOWER(name) = ?`, normalizeLower(name)).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	res, err := a.db.Exec(`INSERT INTO manufacturers(name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) CreateManufacturer(name string) (int64, error) {
	name = normalizeName(name)
	if name == "" {
		return 0, errors.New("name is required")
	}
	res, err := a.db.Exec(`INSERT INTO manufacturers(name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) UpdateManufacturer(id int64, name string) error {
	name = normalizeName(name)
	if name == "" {
		return errors.New("name is required")
	}
	_, err := a.db.Exec(`UPDATE manufacturers SET name = ? WHERE id = ?`, name, id)
	return err
}

func (a *App) DeleteManufacturer(id int64) error {
	// Set manufacturer_id to NULL for all sets using this manufacturer
	if _, err := a.db.Exec(`UPDATE sets SET manufacturer_id = NULL WHERE manufacturer_id = ?`, id); err != nil {
		return err
	}
	_, err := a.db.Exec(`DELETE FROM manufacturers WHERE id = ?`, id)
	return err
}

// Types
func (a *App) ListTypes() ([]Type, error) {
	rows, err := a.db.Query(`SELECT id, name FROM types ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Type
	for rows.Next() {
		var t Type
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		items = append(items, t)
	}
	return items, rows.Err()
}

func (a *App) CreateTypeIfMissing(name string) (int64, error) {
	name = normalizeName(name)
	if name == "" {
		return 0, nil
	}
	var id int64
	err := a.db.QueryRow(`SELECT id FROM types WHERE LOWER(name) = ?`, normalizeLower(name)).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	res, err := a.db.Exec(`INSERT INTO types(name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) CreateType(name string) (int64, error) {
	name = normalizeName(name)
	if name == "" {
		return 0, errors.New("name is required")
	}
	res, err := a.db.Exec(`INSERT INTO types(name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) UpdateType(id int64, name string) error {
	name = normalizeName(name)
	if name == "" {
		return errors.New("name is required")
	}
	_, err := a.db.Exec(`UPDATE types SET name = ? WHERE id = ?`, name, id)
	return err
}

func (a *App) DeleteType(id int64) error {
	// Set type_id to NULL for all sets using this type
	if _, err := a.db.Exec(`UPDATE sets SET type_id = NULL WHERE type_id = ?`, id); err != nil {
		return err
	}
	_, err := a.db.Exec(`DELETE FROM types WHERE id = ?`, id)
	return err
}

// GetNextBagSerial returns the next available bag serial number for a given box
func (a *App) GetNextBagSerial(boxID int64) (string, error) {
	if boxID <= 0 {
		return "0001", nil
	}

	var maxSerial sql.NullString
	err := a.db.QueryRow(`
		SELECT MAX(serial_no) 
		FROM bags 
		WHERE box_id = ? AND serial_no GLOB '[0-9]*'
	`, boxID).Scan(&maxSerial)

	if err != nil || !maxSerial.Valid || maxSerial.String == "" {
		return "0001", nil
	}

	// Try to parse as number and increment
	var num int
	_, err = fmt.Sscanf(maxSerial.String, "%d", &num)
	if err != nil {
		return "0001", nil
	}

	return fmt.Sprintf("%04d", num+1), nil
}

// Bags & Sets
func (a *App) CreateBagWithSet(boxID int64, serialNo, setName, manufacturerName, typeName string) (int64, error) {
	setName = normalizeName(setName)
	if setName == "" {
		return 0, errors.New("set name is required")
	}
	if boxID <= 0 {
		return 0, errors.New("box is required")
	}
	serialNo = normalizeName(serialNo)
	if serialNo == "" {
		return 0, errors.New("bag serial is required")
	}

	tx, err := a.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.Exec(`INSERT INTO bags(box_id, serial_no) VALUES (?, ?)`, boxID, serialNo)
	if err != nil {
		return 0, err
	}
	bagID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	var manufacturerID sql.NullInt64
	if manufacturerName != "" {
		id, err := ensureManufacturerTx(tx, manufacturerName)
		if err != nil {
			return 0, err
		}
		manufacturerID = sql.NullInt64{Int64: id, Valid: true}
	}
	var typeID sql.NullInt64
	if typeName != "" {
		id, err := ensureTypeTx(tx, typeName)
		if err != nil {
			return 0, err
		}
		typeID = sql.NullInt64{Int64: id, Valid: true}
	}

	res, err = tx.Exec(`INSERT INTO sets(bag_id, manufacturer_id, type_id, name) VALUES (?, ?, ?, ?)`, bagID, manufacturerID, typeID, setName)
	if err != nil {
		return 0, err
	}
	setID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	return setID, err
}

func ensureManufacturerTx(tx *sql.Tx, name string) (int64, error) {
	name = normalizeName(name)
	if name == "" {
		return 0, nil
	}

	var id int64
	err := tx.QueryRow(`SELECT id FROM manufacturers WHERE LOWER(name) = ?`, normalizeLower(name)).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	res, err := tx.Exec(`INSERT INTO manufacturers(name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ensureTypeTx(tx *sql.Tx, name string) (int64, error) {
	name = normalizeName(name)
	if name == "" {
		return 0, nil
	}

	var id int64
	err := tx.QueryRow(`SELECT id FROM types WHERE LOWER(name) = ?`, normalizeLower(name)).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	res, err := tx.Exec(`INSERT INTO types(name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) UpdateSet(setID int64, setName, manufacturerName, typeName string, boxID int64, bagSerial string) error {
	setName = normalizeName(setName)
	bagSerial = normalizeName(bagSerial)
	if setName == "" {
		return errors.New("set name is required")
	}
	if boxID <= 0 {
		return errors.New("box is required")
	}
	if bagSerial == "" {
		return errors.New("bag serial is required")
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

	var bagID int64
	if err = tx.QueryRow(`SELECT bag_id FROM sets WHERE id = ?`, setID).Scan(&bagID); err != nil {
		return err
	}

	var manufacturerID sql.NullInt64
	if manufacturerName != "" {
		mID, err := ensureManufacturerTx(tx, manufacturerName)
		if err != nil {
			return err
		}
		manufacturerID = sql.NullInt64{Int64: mID, Valid: true}
	}

	var typeID sql.NullInt64
	if typeName != "" {
		tID, err := ensureTypeTx(tx, typeName)
		if err != nil {
			return err
		}
		typeID = sql.NullInt64{Int64: tID, Valid: true}
	}

	if _, err = tx.Exec(`UPDATE sets SET name = ?, manufacturer_id = ?, type_id = ? WHERE id = ?`, setName, manufacturerID, typeID, setID); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE bags SET box_id = ?, serial_no = ? WHERE id = ?`, boxID, bagSerial, bagID); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (a *App) DeleteSet(setID int64) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var photoPath sql.NullString
	var bagID sql.NullInt64
	if err = tx.QueryRow(`SELECT photo_path, bag_id FROM sets WHERE id = ?`, setID).Scan(&photoPath, &bagID); err != nil {
		return err
	}

	if _, err = tx.Exec(`DELETE FROM sets WHERE id = ?`, setID); err != nil {
		return err
	}

	if bagID.Valid {
		if _, err = tx.Exec(`DELETE FROM bags WHERE id = ?`, bagID.Int64); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	if photoPath.Valid && photoPath.String != "" {
		_ = deleteLocalImage(a.paths.ImagesDir, photoPath.String)
	}
	return nil
}

func (a *App) GetSet(setID int64) (SetDetails, error) {
	var details SetDetails
	row := a.db.QueryRow(`
		SELECT s.id, s.name, s.manufacturer_id, IFNULL(m.name,''), s.type_id, IFNULL(tp.name,''), IFNULL(s.photo_path,''), IFNULL(s.photo_source,''),
		       b.id, b.serial_no, bx.id, bx.code, IFNULL(bx.name,''), loc.id, IFNULL(loc.friendly_name,''), IFNULL(loc.note,''),
		       IFNULL(loc.room,''), IFNULL(loc.shelf,''), IFNULL(loc.compartment,'')
		FROM sets s
		JOIN bags b ON b.id = s.bag_id
		JOIN boxes bx ON bx.id = b.box_id
		LEFT JOIN storage_locations loc ON loc.id = bx.location_id
		LEFT JOIN manufacturers m ON m.id = s.manufacturer_id
		LEFT JOIN types tp ON tp.id = s.type_id
		WHERE s.id = ?`, setID)

	var manufacturerID sql.NullInt64
	var typeID sql.NullInt64
	var bag BagInfo
	if err := row.Scan(
		&details.ID, &details.Name, &manufacturerID, &details.ManufacturerName, &typeID, &details.TypeName,
		&details.PhotoPath, &details.PhotoSource,
		&bag.ID, &bag.SerialNo, &bag.BoxID, &bag.BoxCode, &bag.BoxName, &bag.LocationID, &bag.LocationName, &bag.LocationNote,
		&bag.LocationRoom, &bag.LocationShelf, &bag.LocationCompartment,
	); err != nil {
		return details, err
	}
	if manufacturerID.Valid {
		details.ManufacturerID = &manufacturerID.Int64
	}
	if typeID.Valid {
		details.TypeID = &typeID.Int64
	}
	details.Bag = bag

	// tags
	tagRows, err := a.db.Query(`SELECT t.name FROM set_tags st JOIN tags t ON t.id = st.tag_id WHERE st.set_id = ? ORDER BY t.name`, setID)
	if err != nil {
		return details, err
	}
	defer tagRows.Close()
	for tagRows.Next() {
		var t string
		if err := tagRows.Scan(&t); err != nil {
			return details, err
		}
		details.Tags = append(details.Tags, t)
	}
	if err := tagRows.Err(); err != nil {
		return details, err
	}

	eRows, err := a.db.Query(`SELECT id, set_id, name, IFNULL(kind,'') FROM elements WHERE set_id = ? ORDER BY id`, setID)
	if err != nil {
		return details, err
	}
	defer eRows.Close()
	for eRows.Next() {
		var e Product
		if err := eRows.Scan(&e.ID, &e.SetID, &e.Name, &e.Kind); err != nil {
			return details, err
		}
		details.Products = append(details.Products, e)
	}
	return details, eRows.Err()
}

// Produkte
func (a *App) ListProductsBySet(setID int64) ([]Product, error) {
	rows, err := a.db.Query(`SELECT id, set_id, name, IFNULL(kind,'') FROM elements WHERE set_id = ? ORDER BY id`, setID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var elems []Product
	for rows.Next() {
		var e Product
		if err := rows.Scan(&e.ID, &e.SetID, &e.Name, &e.Kind); err != nil {
			return nil, err
		}
		elems = append(elems, e)
	}
	return elems, rows.Err()
}

func (a *App) AddProduct(setID int64, name, kind string) (int64, error) {
	name = normalizeName(name)
	if name == "" {
		return 0, errors.New("produkt name required")
	}
	kind = normalizeLower(kind)
	if kind == "" {
		kind = ""
	}
	if kind != "" && kind != "stempel" && kind != "stanze" {
		return 0, errors.New("invalid produkt kind")
	}
	res, err := a.db.Exec(`INSERT INTO elements(set_id, name, kind) VALUES (?, ?, NULLIF(?, ''))`, setID, name, kind)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) UpdateProduct(id int64, name, kind string) error {
	name = normalizeName(name)
	if name == "" {
		return errors.New("produkt name required")
	}
	kind = normalizeLower(kind)
	if kind != "" && kind != "stempel" && kind != "stanze" {
		return errors.New("invalid produkt kind")
	}
	_, err := a.db.Exec(`UPDATE elements SET name = ?, kind = NULLIF(?, '') WHERE id = ?`, name, kind, id)
	return err
}

func (a *App) DeleteProduct(id int64) error {
	_, err := a.db.Exec(`DELETE FROM elements WHERE id = ?`, id)
	return err
}

// Tags
func (a *App) CreateTagIfMissing(name string) (int64, error) {
	name = normalizeLower(name)
	if name == "" {
		return 0, errors.New("tag cannot be empty")
	}
	var id int64
	err := a.db.QueryRow(`SELECT id FROM tags WHERE LOWER(name) = ?`, name).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	res, err := a.db.Exec(`INSERT INTO tags(name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) SetTags(setID int64, tagNames []string) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec(`DELETE FROM set_tags WHERE set_id = ?`, setID); err != nil {
		return err
	}

	for _, t := range tagNames {
		tag := normalizeLower(t)
		if tag == "" {
			continue
		}
		tagID, err := ensureTagTx(tx, tag)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(`INSERT OR IGNORE INTO set_tags(set_id, tag_id) VALUES (?, ?)`, setID, tagID); err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

func ensureTagTx(tx *sql.Tx, name string) (int64, error) {
	var id int64
	err := tx.QueryRow(`SELECT id FROM tags WHERE LOWER(name) = ?`, normalizeLower(name)).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	res, err := tx.Exec(`INSERT INTO tags(name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) ListTags() ([]string, error) {
	rows, err := a.db.Query(`SELECT name FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

// Tag represents a tag entity
type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (a *App) ListTagsFull() ([]Tag, error) {
	rows, err := a.db.Query(`SELECT id, name FROM tags ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (a *App) CreateTag(name string) (int64, error) {
	name = normalizeLower(name)
	if name == "" {
		return 0, errors.New("name is required")
	}
	res, err := a.db.Exec(`INSERT INTO tags(name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (a *App) UpdateTag(id int64, name string) error {
	name = normalizeLower(name)
	if name == "" {
		return errors.New("name is required")
	}
	_, err := a.db.Exec(`UPDATE tags SET name = ? WHERE id = ?`, name, id)
	return err
}

func (a *App) DeleteTag(id int64) error {
	// Remove all set_tags associations first
	if _, err := a.db.Exec(`DELETE FROM set_tags WHERE tag_id = ?`, id); err != nil {
		return err
	}
	_, err := a.db.Exec(`DELETE FROM tags WHERE id = ?`, id)
	return err
}

// Fuzzy search helper - checks if query parts match target with fuzzy matching
func fuzzyMatch(query, target string) bool {
	query = strings.ToLower(query)
	target = strings.ToLower(target)
	// Exact substring match
	if strings.Contains(target, query) {
		return true
	}
	// Fuzzy match: each query character should appear in order
	qi := 0
	for ti := 0; ti < len(target) && qi < len(query); ti++ {
		if target[ti] == query[qi] {
			qi++
		}
	}
	return qi == len(query)
}

// parseSearchQuery extracts special filters from the query
// Supports: @Box, @Produkt, @Hersteller, @Tag, @Ort
func parseSearchQuery(query string) (searchTerm string, filters map[string]string) {
	filters = make(map[string]string)
	query = strings.TrimSpace(query)

	// Check for @ prefixes (case insensitive)
	prefixes := []struct {
		prefix string
		key    string
	}{
		{"@box ", "box"},
		{"@Box ", "box"},
		{"@BOX ", "box"},
		{"@produkt ", "product"},
		{"@Produkt ", "product"},
		{"@PRODUKT ", "product"},
		{"@product ", "product"},
		{"@Product ", "product"},
		{"@hersteller ", "manufacturer"},
		{"@Hersteller ", "manufacturer"},
		{"@HERSTELLER ", "manufacturer"},
		{"@tag ", "tag"},
		{"@Tag ", "tag"},
		{"@TAG ", "tag"},
		{"@ort ", "location"},
		{"@Ort ", "location"},
		{"@ORT ", "location"},
		{"@standort ", "location"},
		{"@Standort ", "location"},
	}

	for _, p := range prefixes {
		if strings.HasPrefix(query, p.prefix) {
			filters[p.key] = strings.TrimSpace(query[len(p.prefix):])
			return "", filters
		}
	}

	return query, filters
}

// Search with sorting options and special filters
// sortBy: "name" (default), "box", "location", "added"
// Supports @Box, @Produkt, @Hersteller, @Tag, @Ort prefixes
func (a *App) SearchSets(query string, sortBy string) ([]SetSearchResult, error) {
	searchTerm, filters := parseSearchQuery(query)

	// Determine ORDER BY clause
	orderClause := "s.name"
	switch sortBy {
	case "box":
		orderClause = "bx.code, b.serial_no"
	case "location":
		orderClause = "IFNULL(loc.friendly_name,'zzz'), bx.code"
	case "added":
		orderClause = "s.id DESC" // Newer sets have higher IDs
	default:
		orderClause = "s.name"
	}

	var rows *sql.Rows
	var err error

	baseQuery := `
		SELECT s.id, s.name, IFNULL(m.name,''), bx.code, IFNULL(bx.name,''), b.serial_no,
		       IFNULL(loc.friendly_name,''), IFNULL(GROUP_CONCAT(DISTINCT t.name),''), IFNULL(s.photo_path,'')
		FROM sets s
		JOIN bags b ON b.id = s.bag_id
		JOIN boxes bx ON bx.id = b.box_id
		LEFT JOIN storage_locations loc ON loc.id = bx.location_id
		LEFT JOIN manufacturers m ON m.id = s.manufacturer_id
		LEFT JOIN set_tags st ON st.set_id = s.id
		LEFT JOIN tags t ON t.id = st.tag_id
		%s
		GROUP BY s.id
		ORDER BY %s LIMIT 200`

	// Build WHERE clause based on filters or general search
	if len(filters) > 0 {
		// Specific filter search
		var whereClause string
		var args []interface{}

		if boxFilter, ok := filters["box"]; ok {
			like := "%" + strings.ToLower(boxFilter) + "%"
			whereClause = `WHERE LOWER(bx.code) LIKE ? OR LOWER(bx.name) LIKE ?`
			args = []interface{}{like, like}
		} else if productFilter, ok := filters["product"]; ok {
			like := "%" + strings.ToLower(productFilter) + "%"
			whereClause = `LEFT JOIN elements e ON e.set_id = s.id WHERE LOWER(e.name) LIKE ?`
			args = []interface{}{like}
		} else if manuFilter, ok := filters["manufacturer"]; ok {
			like := "%" + strings.ToLower(manuFilter) + "%"
			whereClause = `WHERE LOWER(m.name) LIKE ?`
			args = []interface{}{like}
		} else if tagFilter, ok := filters["tag"]; ok {
			like := "%" + strings.ToLower(tagFilter) + "%"
			whereClause = `WHERE LOWER(t.name) LIKE ?`
			args = []interface{}{like}
		} else if locFilter, ok := filters["location"]; ok {
			like := "%" + strings.ToLower(locFilter) + "%"
			whereClause = `WHERE LOWER(loc.friendly_name) LIKE ? OR LOWER(loc.room) LIKE ?`
			args = []interface{}{like, like}
		}

		rows, err = a.db.Query(fmt.Sprintf(baseQuery, whereClause, orderClause), args...)
	} else if searchTerm == "" {
		// No search - show all
		rows, err = a.db.Query(fmt.Sprintf(baseQuery, "", orderClause))
	} else {
		// General search across all fields
		like := "%" + strings.ToLower(searchTerm) + "%"
		whereClause := `
			LEFT JOIN elements e ON e.set_id = s.id
			WHERE LOWER(s.name) LIKE ? OR LOWER(t.name) LIKE ? OR LOWER(e.name) LIKE ? 
			      OR LOWER(bx.code) LIKE ? OR LOWER(bx.name) LIKE ? OR LOWER(b.serial_no) LIKE ?
			      OR LOWER(loc.friendly_name) LIKE ? OR LOWER(m.name) LIKE ?`
		rows, err = a.db.Query(fmt.Sprintf(baseQuery, whereClause, orderClause),
			like, like, like, like, like, like, like, like)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SetSearchResult
	for rows.Next() {
		var r SetSearchResult
		var tagList string
		if err := rows.Scan(&r.SetID, &r.SetName, &r.ManufacturerName, &r.BoxCode, &r.BoxName, &r.BagSerial, &r.LocationName, &tagList, &r.ThumbnailPath); err != nil {
			return nil, err
		}
		if tagList != "" {
			parts := strings.Split(tagList, ",")
			for _, p := range parts {
				r.Tags = append(r.Tags, strings.TrimSpace(p))
			}
		}
		results = append(results, r)
	}

	// Apply additional fuzzy filtering for general search
	if searchTerm != "" && len(filters) == 0 {
		queryLower := strings.ToLower(searchTerm)
		var filtered []SetSearchResult
		for _, r := range results {
			// Check if any field fuzzy-matches
			if fuzzyMatch(queryLower, r.SetName) ||
				fuzzyMatch(queryLower, r.BoxCode) ||
				fuzzyMatch(queryLower, r.BoxName) ||
				fuzzyMatch(queryLower, r.BagSerial) ||
				fuzzyMatch(queryLower, r.LocationName) ||
				fuzzyMatch(queryLower, r.ManufacturerName) {
				filtered = append(filtered, r)
				continue
			}
			// Check tags
			for _, tag := range r.Tags {
				if fuzzyMatch(queryLower, tag) {
					filtered = append(filtered, r)
					break
				}
			}
		}
		results = filtered
	}

	return results, rows.Err()
}
