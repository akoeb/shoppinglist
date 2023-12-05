package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	// Here we check for any db errors then exit
	if err != nil {
		panic(err)
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db nil")
	}
	return db
}

func migrate(db *sql.DB) {
	sql := `
    CREATE TABLE IF NOT EXISTS items(
        uid VARCHAR NOT NULL PRIMARY KEY,
        title VARCHAR NOT NULL,
        status VARCHAR NOT NULL,
        orderno INTEGER NOT NULL,
        shop_id VARCHAR REFERENCES shops(uid)
    );
    CREATE TABLE IF NOT EXISTS shops(
		uid VARCHAR NOT NULL PRIMARY KEY,
        name VARCHAR NOT NULL,
        color VARCHAR,
		orderno INTEGER NOT NULL 
    );
	CREATE TABLE IF NOT EXISTS versions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		items INTEGER,
		shops INTEGER 
	);
	INSERT INTO versions (id, items, shops) VALUES (1, 0, 0)
  		ON CONFLICT(id) DO NOTHING;
    `

	_, err := db.Exec(sql)
	// Exit if something goes wrong with our SQL statement above
	if err != nil {
		panic(err)
	}
}

// ********************************** //
//     database access functions:     //
// ********************************** //

// versions:
func GetVersions(db *sql.DB) (Versions, error) {
	result := Versions{}
	sql := "SELECT items, shops FROM versions"
	rows, err := db.Query(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&result.ItemVersion, &result.ShopVersion)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// items

// GetAllItems from database
func GetAllItems(db *sql.DB) (ItemCollection, error) {
	result := ItemCollection{}
	result.Items = make([]Item, 0)

	sql := "SELECT uid, title, status, orderno, shop_id FROM items ORDER BY orderno, uid"
	rows, err := db.Query(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		item := Item{}
		var shopId string
		err = rows.Scan(&item.UId, &item.Title, &item.Status, &item.Orderno, &shopId)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
		if len(shopId) > 0 {
			shop, err := GetShopByID(db, shopId)
			if err != nil {
				return result, err
			}
			item.Shop = &shop
		}
		result.Items = append(result.Items, item)
	}

	// add version:
	versions, err := GetVersions(db)
	if err != nil {
		return result, err
	}
	result.Version = versions.ItemVersion

	return result, nil
}

// GetItemByID loads one item from database, identified by its id
func GetItemByID(db *sql.DB, uid string) (Item, error) {
	result := Item{}
	var shopId string
	sql := "SELECT uid, title, status, orderno, shop_id FROM items WHERE uid = ?"
	rows, err := db.Query(sql, uid)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&result.UId, &result.Title, &result.Status, &result.Orderno, &shopId)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
		if len(shopId) > 0 {
			shop, err := GetShopByID(db, shopId)
			if err != nil {
				return result, err
			}
			result.Shop = &shop
		}

	}
	return result, nil
}

// UpsertItem writes an item to database.
// Whether to INSRT or UPDATE is determined by the existence if its ID field
// modifies the item, adds the ID on creates.
func UpsertItem(db *sql.DB, item *Item) error {

	DeleteItemByID(db, item.UId)

	var query = "INSERT INTO items(uid, title, status, orderno, shop_id ) VALUES(?, ?, ?, ?, ?)"

	// Create a prepared SQL statement
	stmt, err := db.Prepare(query)
	// Exit if we get an error
	if err != nil {
		return err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// if shop is not nil, extract shop id
	var shopId string
	if item.Shop != nil {
		shopId = item.Shop.UId
	}

	_, err = stmt.Exec(item.UId, item.Title, item.Status, item.Orderno, shopId)

	return err
}

// ReplaceItemList completely replaces the List in the database
func ReplaceItemList(db *sql.DB, list *ItemCollection) error {
	var errList []error
	// get the original list:
	orig, err := GetAllItems(db)
	// immediate Exit if we get an error
	if err != nil {
		return err
	}
	// create map for easier lookup of ids:
	itemMap := make(map[string]bool)
	for _, item := range orig.Items {
		itemMap[item.UId] = true
	}

	// loop over input list and do upserts on all items recognizing id has been processed
	// if errors occur, they are collected for later
	for _, item := range list.Items {
		err = UpsertItem(db, &item)
		if err != nil {
			errList = append(errList, err)
		} else {
			if itemMap[item.UId] {
				delete(itemMap, item.UId)
			}
		}
	}

	// delete the remaining ones:
	for id := range itemMap {
		_, err = DeleteItemByID(db, id)
		if err != nil {
			errList = append(errList, err)
		}
	}

	// set version:
	var query = "UPDATE versions SET items = ? WHERE id = 1"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(list.Version)
	if err != nil {
		errList = append(errList, err)
	}

	// wrap errors
	if len(errList) > 0 {
		err = errList[0]
		for i := 1; i < len(errList); i++ {
			err = fmt.Errorf("%v; %w", err, errList[i])
		}
		return err
	}
	return nil
}

// DeleteItemByID deletes one item from the database, identified by its id
func DeleteItemByID(db *sql.DB, id string) (int, error) {

	sql := "DELETE FROM items WHERE uid = ?"

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

// shops

// GetAllShops from database
func GetAllShops(db *sql.DB) (ShopCollection, error) {
	result := ShopCollection{}
	result.Shops = make([]Shop, 0)

	sql := "SELECT uid, name, color, orderno FROM shops ORDER BY  orderno, uid"
	rows, err := db.Query(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		shop := Shop{}
		err = rows.Scan(&shop.UId, &shop.Name, &shop.Color, &shop.Orderno)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
		result.Shops = append(result.Shops, shop)
	}

	// add version:
	versions, err := GetVersions(db)
	if err != nil {
		return result, err
	}
	result.Version = versions.ShopVersion

	return result, nil
}

// GetShopByID loads one item from database, identified by its id
func GetShopByID(db *sql.DB, uid string) (Shop, error) {
	result := Shop{}
	sql := "SELECT uid, name, color, orderno FROM shops WHERE uid = ?"
	rows, err := db.Query(sql, uid)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&result.UId, &result.Name, &result.Color, &result.Orderno)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// UpsertShop writes an item to database.
// Whether to INSERT or UPDATE is determined by the existence if its ID field
// modifies the item, adds the ID on creates.
func UpsertShop(db *sql.DB, shop *Shop) error {
	DeleteShopByID(db, shop.UId)

	query := "INSERT INTO shops(uid, name, color, orderno ) VALUES(?, ?, ?, ?)"

	// Create a prepared SQL statement
	stmt, err := db.Prepare(query)
	// Exit if we get an error
	if err != nil {
		return err
	}
	// Make sure to cleanup after the program exits
	defer stmt.Close()

	// Execute
	_, err = stmt.Exec(shop.UId, shop.Name, shop.Color, shop.Orderno)

	return err
}

// ReplaceItemList completely replaces the List in the database
func ReplaceShopList(db *sql.DB, list *ShopCollection) error {
	var errList []error
	// get the original list:
	orig, err := GetAllShops(db)
	// immediate Exit if we get an error
	if err != nil {
		return err
	}
	// create map for easier lookup of ids:
	shopMap := make(map[string]bool)
	for _, item := range orig.Shops {
		shopMap[item.UId] = true
	}

	// loop over input list and do upserts on all items recognizing id has been processed
	// if errors occur, they are collected for later
	for _, shop := range list.Shops {
		err = UpsertShop(db, &shop)
		if err != nil {
			errList = append(errList, err)
		} else {
			if shopMap[shop.UId] {
				delete(shopMap, shop.UId)
			}
		}
	}

	// delete the remaining ones:
	for id := range shopMap {
		_, err = DeleteShopByID(db, id)
		if err != nil {
			errList = append(errList, err)
		}
	}
	// set version:
	var query = "UPDATE versions SET shops = ? WHERE id = 1"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(list.Version)
	if err != nil {
		errList = append(errList, err)
	}

	// wrap errors
	if len(errList) > 0 {
		err = errList[0]
		for i := 1; i < len(errList); i++ {
			err = fmt.Errorf("%v; %w", err, errList[i])
		}
		return err
	}
	return nil
}

// DeleteItemByID deletes one item from the database, identified by its id
func DeleteShopByID(db *sql.DB, id string) (int, error) {

	sql := "DELETE FROM shops WHERE uid = ?"

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
