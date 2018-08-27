package main
import (
    "database/sql"
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
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        title VARCHAR NOT NULL,
        status VARCHAR NOT NULL,
        orderno INTEGER NOT NULL
    );
    `

    _, err := db.Exec(sql)
    // Exit if something goes wrong with our SQL statement above
    if err != nil {
        panic(err)
    }
}
