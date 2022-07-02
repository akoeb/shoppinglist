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
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Orderno int    `json:"orderno"`
	Shop    int    `json:"shop_id"`
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

// GetAllItems from database
func GetAllItems(db *sql.DB) (ItemCollection, error) {
	result := ItemCollection{}
	sql := "SELECT id, title, status, orderno, shop_id FROM items ORDER BY orderno, id"
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
	sql := "SELECT id, title, status, orderno, shop_id FROM items WHERE id = ?"
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
		query = "INSERT INTO items(title, status, orderno, shop_id ) VALUES(?, ?, ?, ?)"
	} else {
		query = `UPDATE items set title = ?, status = ?, orderno = ?, shop_id = ? WHERE id = ?`
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

	var sql string
	if status == "" {
		sql = "DELETE FROM items"
	} else {
		sql = "DELETE FROM items WHERE status = ?"
	}
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

func reOrderItems(db *sql.DB, items map[int]int) error {

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

// Shop is the entity of a shop.
type Shop struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Orderno int    `json:"orderno"`
}

// ShopCollection is a list of Shops
type ShopCollection struct {
	Shops []Shop `json:"shops"`
}

// GetAllItems from database
func GetAllShops(db *sql.DB) (ShopCollection, error) {
	result := ShopCollection{}
	sql := "SELECT id, name, color, orderno FROM shops ORDER BY  orderno, id"
	rows, err := db.Query(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		shop := Shop{}
		err = rows.Scan(&shop.ID, &shop.Name, &shop.Color, &shop.Orderno)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
		result.Shops = append(result.Shops, shop)
	}
	return result, nil
}

// GetShopByID loads one item from database, identified by its id
func GetShopByID(db *sql.DB, id int) (Shop, error) {
	result := Shop{}
	sql := "SELECT id, name, color, orderno FROM shops WHERE id = ?"
	rows, err := db.Query(sql, id)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Name, &result.Color, &result.Orderno)
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
func UpsertShop(db *sql.DB, shop *Shop) error {

	doInsert := true
	if shop.ID > 0 {
		doInsert = false
	}
	var query string
	if doInsert {
		query = "INSERT INTO shops(name, color, orderno ) VALUES(?, ?, ?)"
	} else {
		query = `UPDATE shops set name = ?, color = ?, orderno = ? WHERE id = ?`
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
		result, err = stmt.Exec(shop.Name, shop.Color, shop.Orderno)
	} else {
		result, err = stmt.Exec(shop.Name, shop.Color, shop.Orderno, shop.ID)
	}
	// Exit if we get an error
	if err != nil {
		return err
	}

	// in insert, read the autoincremented id back into struct
	if doInsert {
		id64, err := result.LastInsertId()
		shop.ID = int(id64)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteItemByID deletes one item from the database, identified by its id
func DeleteShopByID(db *sql.DB, id int) (int, error) {

	sql := "DELETE FROM shops WHERE id = ?"

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
