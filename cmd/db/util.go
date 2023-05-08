package cmd

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Connect to a SQLite database and return a GORM database object
func ConnectToSQLiteDatabase(databaseFilePath string) (*gorm.DB, error) {
	// Open a connection to the database
	db, err := gorm.Open(sqlite.Open(databaseFilePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Set up deferred closure of the database connection
	dbSQL, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer dbSQL.Close()

	return db, nil
}

func ConnectToSQLiteDB(dbFilePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func OpenDB(dbFilePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

func OpenDB1(dbFilePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Defer closing the database connection
	defer func() {
		if err := sqlDB.Close(); err != nil {
			// Handle the error if needed
		}
	}()

	return db, nil
}
