package database

import (
    "context"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func NewPostgresDatabase(ctx context.Context, dsn string) (*sqlx.DB, error) {
    initCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    db, err := sqlx.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    err = db.PingContext(initCtx)
    if err != nil {
        return nil, err
    }
    return db, nil
}
