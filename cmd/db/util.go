package cmd

import (
	"context"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	log "github.com/taylormonacelli/deliverhalf/cmd/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect to a SQLite database and return a GORM database object
func ConnectToSQLiteDatabase(databaseFilePath string) (*gorm.DB, error) {
	// Create the custom logger that forwards GORM logs to the global logger
	gormLogger := CustomLogger{logger: log.Logger}

	// Open a connection to the database
	db, err := gorm.Open(sqlite.Open(databaseFilePath), &gorm.Config{
		Logger: gormLogger,
	})
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

// CustomLogger is a custom logger that forwards GORM logs to logrus for Trace level
type CustomLogger struct {
	logger *logrus.Logger
}

// LogMode sets the log mode for the logger
func (l CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info logs informational messages
func (l CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Tracef(msg, data...)
}

// Warn logs warning messages
func (l CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warnf(msg, data...)
}

// Error logs error messages
func (l CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Errorf(msg, data...)
}

// Trace logs SQL statements at the Trace level
func (l CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logger.IsLevelEnabled(logrus.TraceLevel) {
		sql, rows := fc()
		l.logger.Tracef("%s[%.2fms] %s, row count %d\n",
			err,
			float64(time.Since(begin).Milliseconds()),
			sql,
			rows)
	}
}

func Maintenance() {
	// open a SQLite database
	log.Logger.Debugf("cleaning db")

	// get a list of all tables in the database
	var tables []string
	result := Db.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables)
	if result.Error != nil {
		log.Logger.Fatalf("error selecting from sqlite_master: %s", result.Error)
	}

	// loop over each table
	for _, table := range tables {
		// determine the cutoff date for records to delete
		since := time.Now().AddDate(0, -6, 0)

		// delete records older than the cutoff date
		result := Db.Table(table).Where("created_at < ?", since).Delete(nil)
		if result.Error != nil {
			log.Logger.Warnf("could not delete records from table: %s", table)
		}

		// print out the number of deleted records for this table
		log.Logger.Tracef("Deleted %d records from table %s", result.RowsAffected, table)
	}
}
