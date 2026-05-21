package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/sony/gobreaker/v2"
)

type auctionClient struct {
	baseURL    string
	httpClient *http.Client
	cb         *gobreaker.CircuitBreaker[*domain.AuctionInfo]
}

func NewAuctionClient(baseURL string) domain.AuctionClient {
	cb := gobreaker.NewCircuitBreaker[*domain.AuctionInfo](gobreaker.Settings{
		Name:        "auction-service",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
	})

	return &auctionClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cb:         cb,
	}
}

type auctionResponse struct {
	ID         string `json:"id"`
	SellerID   string `json:"seller_id"`
	StartPrice int64  `json:"start_price"`
	Status     string `json:"status"`
}

func (c *auctionClient) GetAuction(ctx context.Context, auctionID string) (*domain.AuctionInfo, error) {
	result, err := c.cb.Execute(func() (*domain.AuctionInfo, error) {
		return c.getAuctionWithRetry(ctx, auctionID)
	})
	if err != nil {
		if err == gobreaker.ErrOpenState || err == gobreaker.ErrTooManyRequests {
			return nil, errors.Internal("Auction service temporarily unavailable")
		}
		return nil, err
	}
	return result, nil
}

func (c *auctionClient) getAuctionWithRetry(ctx context.Context, auctionID string) (*domain.AuctionInfo, error) {
	var lastErr error
	for attempt := range 3 {
		if attempt > 0 {
			delay := time.Duration(100<<(attempt-1)) * time.Millisecond
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		result, err := c.doRequest(ctx, auctionID)
		if err == nil {
			return result, nil
		}
		lastErr = err

		if isClientError(err) {
			return nil, err
		}
	}
	return nil, lastErr
}

func (c *auctionClient) doRequest(ctx context.Context, auctionID string) (*domain.AuctionInfo, error) {
	url := fmt.Sprintf("%s/api/v1/auctions/%s", c.baseURL, auctionID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Internal("Failed to create request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Internal("Failed to call auction service")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.NotFound("Auction not found")
	}
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return nil, errors.BadRequest("Auction service client error")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Internal("Auction service returned unexpected status")
	}

	var ar auctionResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return nil, errors.Internal("Failed to decode auction response")
	}

	return &domain.AuctionInfo{
		ID: ar.ID, SellerID: ar.SellerID,
		StartPrice: ar.StartPrice, Status: ar.Status,
	}, nil
}

func isClientError(err error) bool {
	if ce, ok := err.(errors.CustomError); ok {
		return ce.Status >= 400 && ce.Status < 500
	}
	return false
}
