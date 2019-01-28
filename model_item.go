package main

import (
	"database/sql"
	"fmt"
	"strings"
)

// AllowedStatusCodes for checking that statuses are always correct
var AllowedStatusCodes = []string{"OPEN", "CHECKED"}

func isAllowedStatusCode(code string) bool {
	for _, item := range AllowedStatusCodes {
		if code == item {
			return true
		}
	}
	return false
}

// Item is our shopping list item
type Item struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Status     string `json:"status"`
	Orderno    int    `json:"orderno"`
	LocationID int    `json:"location"`
}

// Valid tells you whether an item is valid
func (i *Item) Valid() (bool, []string) {
	// required:
	var errors []string
	if i.Title == "" {
		errors = append(errors, "Title is missing")
	}
	if !isAllowedStatusCode(i.Status) {
		errors = append(errors, fmt.Sprintf("Status is of wrong format (%s), only following are allowed: %s", i.Status, strings.Join(AllowedStatusCodes, ", ")))
	}
	if len(errors) > 0 {
		return false, errors
	}
	return true, errors
}

// ItemCollection is a collection of shopping list items
type ItemCollection struct {
	Items []Item `json:"items"`
}

// Valid tells you whether all items in the collection are valid
func (i *ItemCollection) Valid() bool {
	for _, item := range i.Items {
		if ok, _ := item.Valid(); !ok {
			return false
		}
	}
	return true
}

// GetItemsInCategory from database
func GetItemsInCategory(db *sql.DB, categoryID int, locationID int) (ItemCollection, error) {
	result := ItemCollection{}
	var err error
	var rows *sql.Rows

	if locationID > 0 {
		sql := "SELECT id, title, status, orderno, location_id FROM items WHERE category_id = ? AND location_id = ? ORDER BY orderno, id"
		rows, err = db.Query(sql, categoryID, locationID)
	} else {
		sql := "SELECT id, title, status, orderno, location_id FROM items WHERE category_id = ? ORDER BY orderno, id"
		rows, err = db.Query(sql, categoryID)
	}

	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ID, &item.Title, &item.Status, &item.Orderno, &item.LocationID)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}

// GetItemByIDAndCategory loads one item from database, identified by its id
func GetItemByIDAndCategory(db *sql.DB, categoryID int, itemID int) (Item, error) {
	result := Item{}
	sql := "SELECT id, title, status, orderno, location_id FROM items WHERE id = ? and category_id = ?"
	rows, err := db.Query(sql, itemID, categoryID)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Title, &result.Status, &result.Orderno)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// UpsertItem writes an item to database.
// Whether to INSERT or UPDATE is determined by the existence if its ID field
// modifies the item, adds the ID on creates.
func UpsertItem(db *sql.DB, categoryID int, item *Item) error {
	// do category and/or location exist?
	category, err := GetCategoryByID(db, categoryID)
	if err == nil {
		return fmt.Errorf("Can not get category by id %v: %v", categoryID, err)
	}
	if category.ID == 0 {
		return fmt.Errorf("Can not get category by id %v: Does not Exist", categoryID)
	}

	if item.LocationID > 0 {
		location, err := GetLocationByID(db, item.LocationID)
		if err == nil {
			return fmt.Errorf("Can not get location by id %v: %v", item.LocationID, err)
		}
		if location.ID == 0 {
			return fmt.Errorf("Can not get location by id %v: Does not Exist", item.LocationID)
		}
	}

	// insert or update
	doInsert := true
	if item.ID > 0 {
		doInsert = false
	}
	var query string
	if doInsert {
		query = "INSERT INTO items(title, status, orderno, location_id, category_id ) VALUES(?, ?, ?, ?, ?)"
	} else {
		query = `UPDATE items set title = ?, status = ?, orderno = ?, location_id = ?, category_id = ? WHERE id = ?`
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
		result, err = stmt.Exec(item.Title, item.Status, item.Orderno, item.LocationID, categoryID)
	} else {
		result, err = stmt.Exec(item.Title, item.Status, item.Orderno, item.LocationID, categoryID, item.ID)
	}
	// Exit if we get an error
	if err != nil {
		return err
	}

	// in insert, read the autoincremented id back into struct
	if doInsert {
		id64, err := result.LastInsertId()
		item.ID = int(id64)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteItemByID deletes one item from the database, identified by its id
func DeleteItemByID(db *sql.DB, categoryID int, itemID int) (int, error) {

	sql := "DELETE FROM items WHERE id = ? and category_id = ?"

	// Create a prepared SQL statement
	stmt, err := db.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		return 0, err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	result, err := stmt.Exec(itemID, categoryID)
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
