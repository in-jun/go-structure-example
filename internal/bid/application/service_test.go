package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/in-jun/go-structure-example/internal/bid/application/command"
	"github.com/in-jun/go-structure-example/internal/bid/application/query"
	"github.com/in-jun/go-structure-example/internal/bid/domain"
	"github.com/in-jun/go-structure-example/internal/bid/domain/entity"
	sharedQuery "github.com/in-jun/go-structure-example/internal/shared/query"
	domainEvent "github.com/in-jun/go-structure-example/internal/bid/domain/event"
	domainService "github.com/in-jun/go-structure-example/internal/bid/domain/service"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type mockBidRepo struct {
	bid  *entity.Bid
	bids []*entity.Bid
	total int64
	err  error
}

func (m *mockBidRepo) Save(_ context.Context, _ *entity.Bid) error { return m.err }
func (m *mockBidRepo) FindHighestByAuctionID(_ context.Context, _ string, _ ...sharedQuery.Option) (*entity.Bid, error) {
	return m.bid, m.err
}
func (m *mockBidRepo) FindByAuctionID(_ context.Context, _ string, _, _ int) ([]*entity.Bid, int64, error) {
	return m.bids, m.total, m.err
}

type mockAuctionClient struct {
	info *domain.AuctionInfo
	err  error
}

func (m *mockAuctionClient) GetAuction(_ context.Context, _ string) (*domain.AuctionInfo, error) {
	return m.info, m.err
}

type mockPublisher struct{}

func (m *mockPublisher) Publish(_ context.Context, _ ...domainEvent.Event) error { return nil }

type mockTransactor struct{}

func (m *mockTransactor) WithinTransaction(_ context.Context, fn func(ctx context.Context) error, _ ...transaction.TxOption) error {
	return fn(context.Background())
}

type mockEventReader struct{}

func (m *mockEventReader) FindByAuctionID(_ context.Context, _ string) ([]domainEvent.StoredEvent, error) {
	return nil, nil
}

func newTestService(repo *mockBidRepo, client *mockAuctionClient) *service {
	return NewService(
		command.NewPlaceBidHandler(repo, client, &domainService.BidPolicy{}, &mockPublisher{}, &mockTransactor{}),
		command.NewDetermineWinnerHandler(repo, &mockPublisher{}, &mockTransactor{}),
		query.NewGetHighestHandler(repo),
		query.NewListBidsHandler(repo),
		query.NewEventHistoryHandler(&mockEventReader{}),
	)
}

func TestBidService_PlaceBid(t *testing.T) {
	auctionID := uuid.New().String()
	sellerID := uuid.New().String()
	bidderID := uuid.New().String()

	client := &mockAuctionClient{
		info: &domain.AuctionInfo{
			ID: auctionID, SellerID: sellerID, StartPrice: 1000, Status: "open",
		},
	}
	svc := newTestService(&mockBidRepo{}, client)

	result, err := svc.PlaceBid(context.Background(), command.PlaceBid{
		UserID:    bidderID,
		AuctionID: auctionID,
		Amount:    1000,
	})
	if err != nil {
		t.Fatalf("PlaceBid() error = %v", err)
	}
	if result.Amount != 1000 {
		t.Errorf("Amount = %d, want 1000", result.Amount)
	}
}

func TestBidService_PlaceBid_SelfBid(t *testing.T) {
	sellerID := uuid.New().String()
	auctionID := uuid.New().String()

	client := &mockAuctionClient{
		info: &domain.AuctionInfo{
			ID: auctionID, SellerID: sellerID, StartPrice: 1000, Status: "open",
		},
	}
	svc := newTestService(&mockBidRepo{}, client)

	_, err := svc.PlaceBid(context.Background(), command.PlaceBid{
		UserID:    sellerID,
		AuctionID: auctionID,
		Amount:    1000,
	})
	if err == nil {
		t.Error("expected error for self bid")
	}
}

func TestBidService_PlaceBid_NotOpen(t *testing.T) {
	auctionID := uuid.New().String()
	bidderID := uuid.New().String()

	client := &mockAuctionClient{
		info: &domain.AuctionInfo{
			ID: auctionID, SellerID: uuid.New().String(), StartPrice: 1000, Status: "closed",
		},
	}
	svc := newTestService(&mockBidRepo{}, client)

	_, err := svc.PlaceBid(context.Background(), command.PlaceBid{
		UserID:    bidderID,
		AuctionID: auctionID,
		Amount:    1000,
	})
	if err == nil {
		t.Error("expected error for closed auction")
	}
}

func TestBidService_ListBids(t *testing.T) {
	auctionID := uuid.New().String()
	now := time.Now()
	b1 := entity.ReconstructBid(uuid.New().String(), auctionID, uuid.New().String(), 2000, now)
	b2 := entity.ReconstructBid(uuid.New().String(), auctionID, uuid.New().String(), 1500, now)

	repo := &mockBidRepo{bids: []*entity.Bid{b1, b2}, total: 2}
	svc := newTestService(repo, &mockAuctionClient{})

	result, err := svc.ListBids(context.Background(), query.ListBids{
		AuctionID: auctionID, Page: 1, Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListBids() error = %v", err)
	}
	if result.Total != 2 {
		t.Errorf("Total = %d, want 2", result.Total)
	}
}

func TestBidService_GetHighest(t *testing.T) {
	auctionID := uuid.New().String()
	now := time.Now()
	bid := entity.ReconstructBid(uuid.New().String(), auctionID, uuid.New().String(), 5000, now)

	svc := newTestService(&mockBidRepo{bid: bid}, &mockAuctionClient{})

	result, err := svc.GetHighest(context.Background(), query.GetHighest{AuctionID: auctionID})
	if err != nil {
		t.Fatalf("GetHighest() error = %v", err)
	}
	if result.Amount != 5000 {
		t.Errorf("Amount = %d, want 5000", result.Amount)
	}
}

func TestBidService_DetermineWinner(t *testing.T) {
	auctionID := uuid.New().String()
	now := time.Now()
	bid := entity.ReconstructBid(uuid.New().String(), auctionID, uuid.New().String(), 5000, now)

	svc := newTestService(&mockBidRepo{bid: bid}, &mockAuctionClient{})

	err := svc.DetermineWinner(context.Background(), command.DetermineWinner{AuctionID: auctionID})
	if err != nil {
		t.Fatalf("DetermineWinner() error = %v", err)
	}
}

func TestBidService_DetermineWinner_NoBids(t *testing.T) {
	auctionID := uuid.New().String()
	svc := newTestService(&mockBidRepo{bid: nil}, &mockAuctionClient{})

	err := svc.DetermineWinner(context.Background(), command.DetermineWinner{AuctionID: auctionID})
	if err == nil {
		t.Error("expected error when no bids exist")
	}
}

func TestBidService_PlaceBid_BelowStartPrice(t *testing.T) {
	auctionID := uuid.New().String()
	bidderID := uuid.New().String()

	client := &mockAuctionClient{
		info: &domain.AuctionInfo{
			ID: auctionID, SellerID: uuid.New().String(), StartPrice: 2000, Status: "open",
		},
	}
	svc := newTestService(&mockBidRepo{}, client)

	_, err := svc.PlaceBid(context.Background(), command.PlaceBid{
		UserID:    bidderID,
		AuctionID: auctionID,
		Amount:    500,
	})
	if err == nil {
		t.Error("expected error for bid below start price")
	}
}

func TestBidService_PlaceBid_BidTooLow(t *testing.T) {
	auctionID := uuid.New().String()
	bidderID := uuid.New().String()
	now := time.Now()
	existingBid := entity.ReconstructBid(uuid.New().String(), auctionID, uuid.New().String(), 1000, now)

	client := &mockAuctionClient{
		info: &domain.AuctionInfo{
			ID: auctionID, SellerID: uuid.New().String(), StartPrice: 500, Status: "open",
		},
	}
	svc := newTestService(&mockBidRepo{bid: existingBid}, client)

	_, err := svc.PlaceBid(context.Background(), command.PlaceBid{
		UserID:    bidderID,
		AuctionID: auctionID,
		Amount:    1050,
	})
	if err == nil {
		t.Error("expected error for bid below minimum increment")
	}
}

func TestBidService_GetEvents(t *testing.T) {
	svc := newTestService(&mockBidRepo{}, &mockAuctionClient{})

	result, err := svc.GetEvents(context.Background(), query.EventHistory{AuctionID: uuid.New().String()})
	if err != nil {
		t.Fatalf("GetEvents() error = %v", err)
	}
	if len(result.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(result.Events))
	}
}
