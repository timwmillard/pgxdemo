package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Supported enviroment variables
//
// DATABASE_URL
// PGHOST
// PGPORT
// PGDATABASE
// PGUSER
// PGPASSWORD
// PGPASSFILE
// PGSERVICE
// PGSERVICEFILE
// PGSSLMODE
// PGSSLCERT
// PGSSLKEY
// PGSSLROOTCERT
// PGSSLPASSWORD
// PGAPPNAME
// PGCONNECT_TIMEOUT
// PGTARGETSESSIONATTRS
func DBConnect(ctx context.Context, connString string) (*pgx.Conn, error) {
	if connString == "" {
		connString = os.Getenv("DATABASE_URL")
	}

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		slog.Error("Database: connection", "error", err)
		return nil, err
	}

	slog.Info("Database: successful connection",
		"host", conn.Config().Config.Host,
		"port", conn.Config().Config.Port,
		"user", conn.Config().Config.User,
		"database", conn.Config().Config.Database,
	)

	err = conn.Ping(ctx)
	if err != nil {
		slog.Error("Database: ping", "error", err)
		os.Exit(1)
	}
	slog.Info("Database: ping success")

	return conn, err
}

func DBClose(ctx context.Context, conn *pgx.Conn) {

	slog.Info("Database: connection closing")
	err := conn.Close(ctx)
	if err != nil {
		slog.Error("Database: connection close", "error", err)
		return
	}
	slog.Info("Database: connection closed")
}

type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Execer interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// DB is the interface access the database. It is satisfied by *pgx.Conn, pgx.Tx, *pgxpool.Pool, etc.
type DB interface {
	// Begin starts a new pgx.Tx. It may be a true transaction or a pseudo nested transaction implemented by savepoints.
	Begin(ctx context.Context) (pgx.Tx, error)

	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) (br pgx.BatchResults)
}
