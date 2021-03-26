-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE categories(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        title VARCHAR NOT NULL,
        orderno INTEGER NOT NULL
    );

CREATE TABLE locations(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        title VARCHAR NOT NULL,
        color VARCHAR,
        orderno INTEGER NOT NULL
    );
CREATE TABLE items(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        title VARCHAR NOT NULL,
        category_id integer not null references categories(id),
        location_id integer references location(id),
        status VARCHAR NOT NULL,
        orderno INTEGER NOT NULL
    );


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table items;
drop table locations;
drop table categories;