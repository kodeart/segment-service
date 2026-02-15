package database

import (
    "context"
    "os"
    "testing"

    "github.com/jmoiron/sqlx"
    "github.com/stretchr/testify/require"

    "segment-service/infra/config"
    "segment-service/internal/repository"
    "segment-service/internal/service"
)

type TestEnv struct {
    Config         *config.ServiceConfig
    DB             *sqlx.DB
    SegmentService *service.SegmentsService
}

// NewTestEnv a helper function that initializes a
// testing environment with a configured database
// and services, and returns the environment instance.
func NewTestEnv(t *testing.T) *TestEnv {
    os.Setenv("ENV", "testing")
    cfg, err := config.Load()
    require.NoError(t, err)

    db, err := NewPostgresDatabase(context.Background(), cfg.GetPostgresDsn())
    require.NoError(t, err)

    segRepo := repository.NewSegmentsRepository(db)
    segService := service.NewSegmentService(segRepo)

    return &TestEnv{
        Config:         cfg,
        DB:             db,
        SegmentService: segService,
    }
}

// Teardown deletes the segments table.
func (te *TestEnv) Teardown(t *testing.T) {
    _, err := te.DB.Exec("TRUNCATE TABLE segments RESTART IDENTITY;")
    require.NoError(t, err)
}
