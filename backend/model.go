package main

import (
	"fmt"
	"strings"
	"time"
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
	UId     string `json:"uid"`
	Title   string `json:"title"`
	Status  string `json:"status"`
	Orderno int    `json:"orderno"`
	Shop    *Shop  `json:"shop,omitempty"`
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
	Version int64  `json:"version"`
	Items   []Item `json:"items"`
}

// Valid tells you whether all items in the collection are valid
func (i *ItemCollection) Valid() bool {
	for _, item := range i.Items {
		if ok, _ := item.Valid(); !ok {
			return false
		}
	}
	return i.Version <= time.Now().Unix()
}

// Shop is the entity of a shop.
type Shop struct {
	UId     string `json:"uid"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Orderno int    `json:"orderno"`
}

// ShopCollection is a list of Shops
type ShopCollection struct {
	Version int64  `json:"version"`
	Shops   []Shop `json:"items"`
}

// validation
func (s *ShopCollection) Valid() bool {
	return s.Version <= time.Now().Unix()
}

type Versions struct {
	ItemVersion int64
	ShopVersion int64
}
