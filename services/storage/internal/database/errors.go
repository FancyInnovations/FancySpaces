package database

import "errors"

var (
	ErrDatabaseNotFound      = errors.New("database not found")
	ErrDatabaseAlreadyExists = errors.New("a database with this name already exists")

	ErrCollectionNotFound      = errors.New("collection not found")
	ErrCollectionAlreadyExists = errors.New("a collection with this name already exists")
)
