package cdb

import (
	"database/sql"
	_ "github.com/alexbrainman/odbc"
)

type CdbConnection struct {
	db *sql.DB
}

func Open() (*CdbConnection, error) {

	db, err := sql.Open("odbc", "dsn=CitectAlarms;uid=;pwd=")
	if err != nil {
		return nil, err
	}

	return &CdbConnection{
		db: db,
	}, nil

}

func (this *CdbConnection) Close() error {
	return this.db.Close()
}
