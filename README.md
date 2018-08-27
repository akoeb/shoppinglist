# shopping list #

## description ##

This is a simple web application to maintain a shopping list. You can view and modify this list from your computer at home or from your cell phone in the store. Or let your partner add things to the list while you are running to get the stuff ;-)

This is work in Progress and **NOT USABLE** yet!

You have been warned!

## Backend ##

The backend is written in golang with the echo framework

### Entities ###

* item (Name, Status (OPEN, CHECKED), created, orderno)

### API ###

* show all items, possibility to filter by status, order alphabetically or by dateor orderno  (default)
* add item
* change item
* delete item
* delete multiple items (by status)

## Frontend ##

The frontend is not coded yet

* shows the list of items
* every item can be clicked which toggles "strike through" of the item
* striked through items are of Status "CHECKED", the others have status "OPEN"
* of course items can be added, modified or deleted
* all striked through items can be deleted
* The list shows all items or only the ones in status OPEN
* every item change synchronizes to the backend asynchronously

## TODO ##

* Tests
* Input Validation
* Status constants map
* JSON Bind of body in struct does not work work
* implement delete many
* implement frontend

## License ##

GPLv3, see LICENSE file in this repo.