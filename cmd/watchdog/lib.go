package cmd

import (
	"time"

	db "github.com/taylormonacelli/deliverhalf/cmd/db"
)

func maintainDb(c chan<- string) {
	go func() {
		for {
			db.Maintenance()
			sleep := 24 * time.Hour
			time.Sleep(sleep)
		}
	}()
	c <- "longRunningFunc completed"
}
