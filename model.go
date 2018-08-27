package main

import (
    "database/sql"
)

type Item struct {
    ID   int    `json:"id"`
    Title string `json:"title"`
    Status string `json:"status"`
    Orderno int `json:"orderno"`
}

type ItemCollection struct {
    Items []Item `json:"items"`
}


// Get all Items from database
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

// Get all Items from database
func GetItemById(db *sql.DB, id int) (Item, error) {
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

func UpsertItem(db *sql.DB, item Item) (Item, error) {

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
        return item, err
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
        return item, err
    }

	// in insert, read the autoincremented id back into struct
    if doInsert {
        id64, err := result.LastInsertId()
        item.ID = int(id64)
        if err != nil {
            return item, err
        }
    }
	return item,nil
}
func DeleteItemById(db *sql.DB, id int) (int, error) {

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

