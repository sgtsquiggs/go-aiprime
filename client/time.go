package client

import (
	"context"
	"net/http"

	"github.com/sgtsquiggs/go-aiprime/models"
)

const (
	timeBasePath = "/api/time"
)

type TimeService interface {
	Update(context.Context, *TimeRequest) (*Response, error)
}

type TimeServiceOp struct {
	client *Client
}

type TimeRequest models.Time

var _ TimeService = &TimeServiceOp{}

// Update time
func (s *TimeServiceOp) Update(ctx context.Context, tr *TimeRequest) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, timeBasePath, tr)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
