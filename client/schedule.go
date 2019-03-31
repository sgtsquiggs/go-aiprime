package client

import (
	"context"
	"net/http"
	"path"

	"github.com/sgtsquiggs/go-aiprime/models"
)

const (
	scheduleBasePath    = "/api/schedule"
	scheduleDefaultName = "default"
)

type ScheduleService interface {
	Update(context.Context, *ScheduleRequest) (*Response, error)
	UpdateLunar(context.Context, *LunarScheduleRequest) (*Response, error)
}

type ScheduleServiceOp struct {
	client *Client
}

type ScheduleRequest models.Schedule

var _ ScheduleService = &ScheduleServiceOp{}

// Update or update a Schedule
func (s *ScheduleServiceOp) Update(ctx context.Context, sr *ScheduleRequest) (*Response, error) {
	if sr.Name == "" {
		sr.Name = scheduleDefaultName
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, scheduleBasePath, sr)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

type LunarScheduleRequest models.LunarSchedule

func (s *ScheduleServiceOp) UpdateLunar(ctx context.Context, lr *LunarScheduleRequest) (*Response, error) {
	basePath := path.Join(scheduleBasePath, "lunar")
	req, err := s.client.NewRequest(ctx, http.MethodPut, basePath, lr)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
