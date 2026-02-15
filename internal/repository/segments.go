package repository

import (
    "context"
    "database/sql"
    "errors"

    "github.com/jmoiron/sqlx"

    "segment-service/internal/entities"
)

//
type SegmentsRepository interface {
    List(ctx context.Context) ([]entities.Segment, error)
    Get(ctx context.Context, id int) (*entities.Segment, error)
    Create(ctx context.Context, segment *entities.Segment) (*entities.Segment, error)
    Update(ctx context.Context, segment *entities.Segment) (*entities.Segment, error)
    Delete(ctx context.Context, id int) error
}

type PostgresSegments struct {
    db *sqlx.DB
}

func NewSegmentsRepository(db *sqlx.DB) *PostgresSegments {
    return &PostgresSegments{db: db}
}

func (r *PostgresSegments) List(ctx context.Context) ([]entities.Segment, error) {
    query := `
    SELECT id, name, ttl_seconds, created_at, updated_at
    FROM segments
    WHERE deleted_at IS NULL`

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    list := make([]entities.Segment, 0)
    for rows.Next() {
        var seg entities.Segment
        if err = rows.Scan(
            &seg.ID,
            &seg.Name,
            &seg.TTLSeconds,
            &seg.CreateAt,
            &seg.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        list = append(list, seg)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return list, nil
}

func (r *PostgresSegments) Get(ctx context.Context, id int) (*entities.Segment, error) {
    query := `SELECT * FROM segments WHERE id = $1`

    var segment entities.Segment
    err := r.db.GetContext(ctx, &segment, query, id)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, entities.ErrSegmentNotFound
    }
    if err != nil {
        return nil, err
    }
    return &segment, nil
}

func (r *PostgresSegments) Create(ctx context.Context, segment *entities.Segment) (*entities.Segment, error) {
    query := `
    INSERT INTO segments (name, ttl_seconds, created_at, updated_at)
    VALUES ($1, $2, now(), now()) 
    RETURNING *`

    err := r.db.GetContext(ctx, segment, query, segment.Name, segment.TTLSeconds)
    if err != nil {
        return nil, err
    }
    return segment, nil
}

func (r *PostgresSegments) Update(ctx context.Context, segment *entities.Segment) (*entities.Segment, error) {
    query := `
    UPDATE segments SET
        name = $1,
        ttl_seconds = $2,
        updated_at = now(),
        deleted_at = $3
    WHERE id = $4 
    RETURNING *`

    err := r.db.GetContext(ctx, segment, query, segment.Name, segment.TTLSeconds, segment.DeletedAt, segment.ID)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, entities.ErrSegmentNotFound
    }
    if err != nil {
        return nil, err
    }
    return segment, nil
}

func (r *PostgresSegments) Delete(ctx context.Context, id int) error {
    query := `UPDATE segments SET deleted_at = now() WHERE id = $1`

    res, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    affected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return entities.ErrSegmentNotFound
    }
    return nil
}
