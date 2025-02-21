package utils

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func SetupDBConnection(dsn string) (*bun.DB, error) {
	
	maxOpenConnectionsStr := os.Getenv("DB_MAX_OPEN_CONNECTIONS")
	maxOpenConnections, err := strconv.Atoi(maxOpenConnectionsStr)

	if err != nil || maxOpenConnectionsStr == "" {
		maxOpenConnections = 10
	}

	pgDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	pgDB.SetMaxOpenConns(maxOpenConnections)

	db := bun.NewDB(pgDB, pgdialect.New())

	dbConnectionError := db.Ping()
	if dbConnectionError != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", dbConnectionError)
	}

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db, nil
}
