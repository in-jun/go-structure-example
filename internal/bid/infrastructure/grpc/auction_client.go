package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"github.com/sony/gobreaker/v2"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	auctionv1 "github.com/in-jun/go-structure-example/proto/auction/v1"
)

type AuctionClient struct {
	client auctionv1.AuctionServiceClient
	conn   *grpc.ClientConn
	cb     *gobreaker.CircuitBreaker[*domain.AuctionInfo]
}

func NewAuctionClient(addr string) (*AuctionClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	cb := gobreaker.NewCircuitBreaker[*domain.AuctionInfo](gobreaker.Settings{
		Name:        "auction-grpc",
		MaxRequests: 1,
		Interval:    60 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
	})

	return &AuctionClient{
		client: auctionv1.NewAuctionServiceClient(conn),
		conn:   conn,
		cb:     cb,
	}, nil
}

func (c *AuctionClient) GetAuction(ctx context.Context, auctionID string) (*domain.AuctionInfo, error) {
	result, err := c.cb.Execute(func() (*domain.AuctionInfo, error) {
		resp, err := c.client.GetAuction(ctx, &auctionv1.GetAuctionRequest{
			AuctionId: auctionID,
		})
		if err != nil {
			return nil, fromGRPCError(err)
		}
		return &domain.AuctionInfo{
			ID:         resp.Id,
			SellerID:   resp.SellerId,
			StartPrice: resp.StartPrice,
			Status:     resp.Status,
		}, nil
	})
	if err != nil {
		if err == gobreaker.ErrOpenState || err == gobreaker.ErrTooManyRequests {
			return nil, errors.Internal("Auction service temporarily unavailable")
		}
		return nil, err
	}
	return result, nil
}

func (c *AuctionClient) Close() error {
	return c.conn.Close()
}

func fromGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return errors.Internal("Auction service communication error")
	}
	switch st.Code() {
	case codes.NotFound:
		return errors.NotFound(st.Message())
	case codes.InvalidArgument:
		return errors.BadRequest(st.Message())
	default:
		return errors.Internal(st.Message())
	}
}
