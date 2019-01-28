package main

import "database/sql"

// Category as model "category"
type Category struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Orderno int    `json:"orderno"`
}

// Valid tells you whether an item is valid
func (c *Category) Valid() (bool, []string) {
	// required:
	var errors []string
	if c.Title == "" {
		errors = append(errors, "Title is missing")
	}
	if len(errors) > 0 {
		return false, errors
	}
	return true, errors
}

// CategoryCollection is a collection of list categories
type CategoryCollection struct {
	Categories []Category `json:"categories"`
}

// Valid tells you whether all items in the collection are valid
func (c *CategoryCollection) Valid() bool {
	for _, cat := range c.Categories {
		if ok, _ := cat.Valid(); !ok {
			return false
		}
	}
	return true
}

// GetAllCategories from DB
func GetAllCategories(db *sql.DB) (CategoryCollection, error) {
	result := CategoryCollection{}
	sql := "SELECT id, title, orderno FROM categories ORDER BY orderno, id"
	rows, err := db.Query(sql)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		cat := Category{}
		err = rows.Scan(&cat.ID, &cat.Title, &cat.Orderno)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
		result.Categories = append(result.Categories, cat)
	}
	return result, nil
}

// GetCategoryByID loads one category from database, identified by its id
func GetCategoryByID(db *sql.DB, id int) (Category, error) {
	result := Category{}
	sql := "SELECT id, title, orderno FROM categories WHERE id = ?"
	rows, err := db.Query(sql, id)
	// Exit if the SQL doesn't work for some reason
	if err != nil {
		return result, err
	}
	// make sure to cleanup when the program exits
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&result.ID, &result.Title, &result.Orderno)
		// Exit if we get an error
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

// UpsertCategory writes a category to database.
// Whether to INSERT or UPDATE is determined by the existence if its ID field
// modifies the item, adds the ID on creates.
func UpsertCategory(db *sql.DB, category *Category) error {

	doInsert := true
	if category.ID > 0 {
		doInsert = false
	}
	var query string
	if doInsert {
		query = "INSERT INTO categories(title, orderno ) VALUES(?, ?)"
	} else {
		query = `UPDATE categories set title = ?, orderno = ? WHERE id = ?`
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
		result, err = stmt.Exec(category.Title, category.Orderno)
	} else {
		result, err = stmt.Exec(category.Title, category.Orderno, category.ID)
	}
	// Exit if we get an error
	if err != nil {
		return err
	}

	// in insert, read the autoincremented id back into struct
	if doInsert {
		id64, err := result.LastInsertId()
		category.ID = int(id64)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteCategoryByID deletes one category from the database, identified by its id
// TODO: Transaction
func DeleteCategoryByID(db *sql.DB, id int) (int, error) {

	// first all items in the category:
	sql := "DELETE FROM items WHERE category_id = ?"

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

	// now the category itself
	sql = "DELETE FROM categories WHERE id = ?"

	// Create a prepared SQL statement
	stmt, err = db.Prepare(sql)
	// Exit if we get an error
	if err != nil {
		return 0, err
	}

	// Execute
	result, err = stmt.Exec(id)

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
