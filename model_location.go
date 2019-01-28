package main

import (
	"database/sql"
)

// Location model "location"
type Location struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Color   string `json:"color"`
	Orderno int    `json:"orderno"`
}

// Valid tells you whether an item is valid
func (l *Location) Valid() (bool, []string) {
	// required:
	var errors []string
	if l.Title == "" {
		errors = append(errors, "Title is missing")
	}
	if len(errors) > 0 {
		return false, errors
	}
	return true, errors
}

// LocationCollection is a collection of list categories
type LocationCollection struct {
	Locations []Location `json:"categories"`
}

// Valid tells you whether all items in the collection are valid
func (l *LocationCollection) Valid() bool {
	for _, loc := range l.Locations {
		if ok, _ := loc.Valid(); !ok {
			return false
		}
	}
	return true
}

// GetLocations returns the list of locations, possibly filtered by items in this location with category == categoryID
func GetLocations(db *sql.DB, categoryID int) (LocationCollection, error) {
	result := LocationCollection{}
	var err error
	var rows *sql.Rows

	if categoryID > 0 {
		sql := "SELECT l.id, l.title, l.color, l.orderno FROM locations l, items i WHERE l.id = i.location_id and i.category_id = ? ORDER BY l.orderno, l.id"
		rows, err = db.Query(sql, categoryID)
	} else {
		sql := "SELECT l.id, l.title, l.color, l.orderno FROM locations l ORDER BY l.orderno, l.id"
		rows, err = db.Query(sql)
	}

	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		location := Location{}
		err = rows.Scan(&location.ID, &location.Title, &location.Color, &location.Orderno)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
		result.Locations = append(result.Locations, location)
	}
	return result, nil
}

// GetLocationByID returns the location for an ID
func GetLocationByID(db *sql.DB, locationID int) (Location, error) {
	result := Location{}
	sql := "SELECT l.id, l.title, l.color, l.orderno FROM locations l WHERE l.id = ? "
	rows, err := db.Query(sql, locationID)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&result.ID, &result.Title, &result.Color, &result.Orderno)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

//UpsertLocation inserts or updates a location in the database. insert/update is determined by the ID attribute
func UpsertLocation(db *sql.DB, location *Location) error {

	// insert or update
	doInsert := true
	if location.ID > 0 {
		doInsert = false
	}
	var query string
	if doInsert {
		query = "INSERT INTO locations(title, color, orderno ) VALUES(?, ?, ?, ?, ?)"
	} else {
		query = `UPDATE locations set title = ?, color = ?, orderno = ? WHERE id = ?`
	}

	// Create a prepared SQL statement
	stmt, err := db.Prepare(query)
	// Exit if we get an error
	if err != nil {
		return err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	var result sql.Result
	if doInsert {
		result, err = stmt.Exec(location.Title, location.Color, location.Orderno)
	} else {
		result, err = stmt.Exec(location.Title, location.Color, location.Orderno, location.ID)
	}
	// Exit if we get an error
	if err != nil {
		return err
	}

	// in insert, read the autoincremented id back into struct
	if doInsert {
		id64, err := result.LastInsertId()
		location.ID = int(id64)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteLocationByID deletes a location identified by its id. all items in this location will be unassigned.
// TODO: transaction
func DeleteLocationByID(db *sql.DB, locationID int) (int, error) {

	sql := "Update items set location_id = NULL WHERE location_id = ? "

	// Create a prepared SQL statement
	stmt, err := db.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		return 0, err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	result, err := stmt.Exec(locationID)

	// Exit if we get an error
	if err != nil {
		return 0, err
	}

	sql = "DELETE FROM locations WHERE id = ? "

	// Create a prepared SQL statement
	stmt, err = db.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		return 0, err
	}

	// Execute
	result, err = stmt.Exec(locationID)

	// Exit if we get an error
	if err != nil {
		return 0, err
	}

	numDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(numDeleted), nil
}

/*
TODO:
func reorderItems(db *sql.DB, items map[int]int) error {

	sql := "UPDATE items set orderno = ? where id = ?"

	// Create a prepared SQL statement
	tx, err := db.Begin()
	// Exit if we get an error
	if err != nil {
		return err
	}

	// prepare statement in this transaction
	stmt, err := tx.Prepare(sql)

	// Exit if we get an error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	for k, v := range items {
		_, err := stmt.Exec(v, k)

		// Exit if we get an error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// all good, commit the transaction
	tx.Commit()

	// and bye
	return nil
}
*/
