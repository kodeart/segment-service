package service

import (
    "context"

    "segment-service/internal/entities"
    "segment-service/internal/repository"
)

type SegmentsService struct {
    segmentsRepo repository.SegmentsRepository
}

func NewSegmentService(r repository.SegmentsRepository) *SegmentsService {
    return &SegmentsService{segmentsRepo: r}
}

func (s *SegmentsService) List(ctx context.Context) (*entities.SegmentsList, error) {
    segments, err := s.segmentsRepo.List(ctx)
    if err != nil {
        return nil, err
    }
    return &entities.SegmentsList{Segments: segments}, nil
}

func (s *SegmentsService) Get(ctx context.Context, id int) (*entities.Segment, error) {
    return s.segmentsRepo.Get(ctx, id)
}

func (s *SegmentsService) Create(ctx context.Context, segment *entities.Segment) (*entities.Segment, error) {
    return s.segmentsRepo.Create(ctx, segment)
}

func (s *SegmentsService) Update(ctx context.Context, segment *entities.Segment) (*entities.Segment, error) {
    return s.segmentsRepo.Update(ctx, segment)
}

func (s *SegmentsService) Delete(ctx context.Context, id int) error {
    return s.segmentsRepo.Delete(ctx, id)
}
