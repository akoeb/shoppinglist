package main

import (
	"database/sql"
	"fmt"
	"strings"
)

// AllowedStatusCodes for checking that statuses are always correct
var AllowedStatusCodes = []string{"OPEN", "CHECKED"}

// Item is our shopping list item
type Item struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Orderno int    `json:"orderno"`
}

// Valid tells you whether an item is valid
func (i *Item) Valid() (bool, []string) {
	// required:
	var errors []string
	if i.Title == "" {
		errors = append(errors, "Title is missing")
	}
	if i.Status != AllowedStatusCodes[0] && i.Status != AllowedStatusCodes[1] {
		errors = append(errors, fmt.Sprintf("Status is of wrong format, only following are allowed: %s", strings.Join(AllowedStatusCodes, ", ")))
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

// GetAllItems from database
func GetAllItems(db *sql.DB) (ItemCollection, error) {
	result := ItemCollection{}
	sql := "SELECT id, title, status, orderno FROM items"
	rows, err := db.Query(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ID, &item.Title, &item.Status, &item.Orderno)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}

// GetItemByID loads one item from database, identified by its id
func GetItemByID(db *sql.DB, id int) (Item, error) {
	result := Item{}
	sql := "SELECT id, title, status, orderno FROM items WHERE id = ?"
	rows, err := db.Query(sql, id)
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
// Whether to INSRT or UPDATE is determined by the existence if its ID field
// modifies the item, adds the ID on creates.
func UpsertItem(db *sql.DB, item *Item) error {

	doInsert := true
	if item.ID > 0 {
		doInsert = false
	}
	var query string
	if doInsert {
		query = "INSERT INTO items(title, status, orderno ) VALUES(?, ?, ?)"
	} else {
		query = `UPDATE items set title = ?, status = ?, orderno = ? WHERE id = ?`
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
		result, err = stmt.Exec(item.Title, item.Status, item.Orderno)
	} else {
		result, err = stmt.Exec(item.Title, item.Status, item.Orderno, item.ID)
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
func DeleteItemByID(db *sql.DB, id int) (int, error) {

	sql := "DELETE FROM items WHERE id = ?"

	// Create a prepared SQL statement
	stmt, err := db.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		return 0, err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	result, err := stmt.Exec(id)
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
func deleteManyItemsByStatus(db *sql.DB, status string) (int, error) {

	sql := "DELETE FROM items WHERE status = ?"

	// Create a prepared SQL statement
	stmt, err := db.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		return 0, err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	result, err := stmt.Exec(status)
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
