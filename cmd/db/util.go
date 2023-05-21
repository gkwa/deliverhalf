package cmd

import (
	"time"

	"github.com/glebarez/sqlite"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
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

func Maintenance() {
	// open a SQLite database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Logger.Fatalf("can't connect to db: %s", err)
	}
	log.Logger.Debugf("cleaning db")

	// get a list of all tables in the database
	var tables []string
	result := db.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables)
	if result.Error != nil {
		log.Logger.Fatalf("error selecting from sqlite_master: %s", err)
	}

	// loop over each table
	for _, table := range tables {
		// determine the cutoff date for records to delete
		since := time.Now().AddDate(0, -6, 0)

		// delete records older than the cutoff date
		result := db.Table(table).Where("created_at < ?", since).Delete(nil)
		if result.Error != nil {
			log.Logger.Warnf("could not delete records from table: %s", table)
		}

		// print out the number of deleted records for this table
		log.Logger.Tracef("Deleted %d records from table %s", result.RowsAffected, table)
	}
}
