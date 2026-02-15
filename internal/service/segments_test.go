package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "segment-service/infra/database"
    "segment-service/internal/entities"
)

func TestSegmentCreateAndGet(t *testing.T) {
    env := database.NewTestEnv(t)
    t.Cleanup(func() { env.Teardown(t) })

    ctx := context.Background()
    in := &entities.Segment{
        Name: "test-segment",
    }

    created, err := env.SegmentService.Create(ctx, in)
    require.NoError(t, err)
    assert.Same(t, in, created)
    assert.NotZero(t, created.ID)
    assert.Nil(t, created.TTLSeconds)
    assert.NotZero(t, created.CreateAt)
    assert.NotZero(t, created.UpdatedAt)
    assert.Nil(t, created.DeletedAt)

    got, err := env.SegmentService.Get(ctx, created.ID)
    require.NoError(t, err)
    assert.NotSame(t, got, in)
}

func TestSegmentListWithSoftDelete(t *testing.T) {
    env := database.NewTestEnv(t)
    t.Cleanup(func() { env.Teardown(t) })

    ctx := context.Background()
    a, err := env.SegmentService.Create(ctx, &entities.Segment{Name: "a"})
    require.NoError(t, err)
    b, err := env.SegmentService.Create(ctx, &entities.Segment{Name: "b"})
    require.NoError(t, err)

    list, err := env.SegmentService.List(ctx)
    require.NoError(t, err)
    assert.Equal(t, 2, len(list.Segments))

    // Delete a, only b should be present
    err = env.SegmentService.Delete(ctx, a.ID)
    require.NoError(t, err)

    list, err = env.SegmentService.List(ctx)
    require.NoError(t, err)
    assert.Equal(t, 1, len(list.Segments))
    assert.Equal(t, list.Segments[0].ID, b.ID)
}

func TestSegmentDeleteNotFound(t *testing.T) {
    env := database.NewTestEnv(t)
    t.Cleanup(func() { env.Teardown(t) })

    ctx := context.Background()
    err := env.SegmentService.Delete(ctx, 99999)
    assert.ErrorIs(t, err, entities.ErrSegmentNotFound)
}

func TestSegmentUpdate(t *testing.T) {
    env := database.NewTestEnv(t)
    t.Cleanup(func() { env.Teardown(t) })

    ctx := context.Background()
    created, err := env.SegmentService.Create(ctx, &entities.Segment{Name: "before"})
    require.NoError(t, err)

    ttl := 60
    updated, err := env.SegmentService.Update(ctx, &entities.Segment{
        ID:         created.ID,
        Name:       "after",
        TTLSeconds: &ttl,
    })
    require.NoError(t, err)
    assert.Equal(t, "after", updated.Name)
    assert.Equal(t, 60, *updated.TTLSeconds)

    got, err := env.SegmentService.Get(ctx, created.ID)
    require.NoError(t, err)
    assert.Equal(t, "after", got.Name)
    assert.Equal(t, 60, *got.TTLSeconds)
}
