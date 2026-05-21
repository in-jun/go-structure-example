package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/in-jun/go-structure-example/internal/auction/application/command"
	"github.com/in-jun/go-structure-example/internal/auction/application/query"
	"github.com/in-jun/go-structure-example/internal/auction/domain/entity"
	sharedQuery "github.com/in-jun/go-structure-example/internal/shared/query"
	domainEvent "github.com/in-jun/go-structure-example/internal/auction/domain/event"
	domainService "github.com/in-jun/go-structure-example/internal/auction/domain/service"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type mockAuctionRepo struct {
	auction  *entity.Auction
	auctions []*entity.Auction
	total    int64
	err      error
}

func (m *mockAuctionRepo) Save(_ context.Context, _ *entity.Auction) error { return m.err }
func (m *mockAuctionRepo) FindByID(_ context.Context, _ string, _ ...sharedQuery.Option) (*entity.Auction, error) {
	return m.auction, m.err
}
func (m *mockAuctionRepo) FindAll(_ context.Context, _, _ int) ([]*entity.Auction, int64, error) {
	return m.auctions, m.total, m.err
}
func (m *mockAuctionRepo) Update(_ context.Context, _ *entity.Auction) error { return m.err }

type mockPublisher struct{}

func (m *mockPublisher) Publish(_ context.Context, _ ...domainEvent.Event) error { return nil }

type mockEventReader struct{}

func (m *mockEventReader) FindByAuctionID(_ context.Context, _ string) ([]domainEvent.StoredEvent, error) {
	return nil, nil
}

type mockTransactor struct{}

func (m *mockTransactor) WithinTransaction(_ context.Context, fn func(ctx context.Context) error, _ ...transaction.TxOption) error {
	return fn(context.Background())
}

func newTestService(repo *mockAuctionRepo) *service {
	scheduler := &domainService.AuctionScheduler{}
	publisher := &mockPublisher{}
	reader := &mockEventReader{}
	transactor := &mockTransactor{}

	return NewService(
		command.NewCreateHandler(repo, publisher, scheduler, transactor),
		command.NewOpenHandler(repo, publisher, transactor),
		command.NewCloseHandler(repo, publisher, transactor),
		command.NewSettleHandler(repo, publisher, transactor),
		command.NewCancelHandler(repo, publisher, transactor),
		query.NewGetHandler(repo),
		query.NewListHandler(repo),
		query.NewEventHistoryHandler(reader),
	)
}

func TestAuctionService_Create(t *testing.T) {
	userID := uuid.New().String()
	svc := newTestService(&mockAuctionRepo{})

	result, err := svc.Create(context.Background(), command.Create{
		UserID:      userID,
		Title:       "Test Auction",
		Description: "Description",
		StartPrice:  1000,
		EndTime:     time.Now().Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.Title != "Test Auction" {
		t.Errorf("Title = %q, want %q", result.Title, "Test Auction")
	}
	if result.Status != entity.StatusDraft {
		t.Errorf("Status = %q, want %q", result.Status, entity.StatusDraft)
	}
}

func TestAuctionService_Create_InvalidInput(t *testing.T) {
	svc := newTestService(&mockAuctionRepo{})

	_, err := svc.Create(context.Background(), command.Create{
		UserID: "bad-id",
		Title:  "Test",
	})
	if err == nil {
		t.Error("expected error for invalid user ID")
	}
}

func TestAuctionService_Open(t *testing.T) {
	userID := uuid.New().String()
	auction, _ := entity.NewAuction(userID, "Test", "", 100, time.Now().Add(2*time.Hour))
	svc := newTestService(&mockAuctionRepo{auction: auction})

	err := svc.Open(context.Background(), command.Open{
		UserID:    userID,
		AuctionID: auction.ID(),
	})
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
}

func TestAuctionService_Open_NotOwner(t *testing.T) {
	ownerID := uuid.New().String()
	otherID := uuid.New().String()
	auction, _ := entity.NewAuction(ownerID, "Test", "", 100, time.Now().Add(2*time.Hour))
	svc := newTestService(&mockAuctionRepo{auction: auction})

	err := svc.Open(context.Background(), command.Open{
		UserID:    otherID,
		AuctionID: auction.ID(),
	})
	if err == nil {
		t.Error("expected error for non-owner")
	}
}

func TestAuctionService_Close(t *testing.T) {
	userID := uuid.New().String()
	auction, _ := entity.NewAuction(userID, "Test", "", 100, time.Now().Add(2*time.Hour))
	if err := auction.Open(); err != nil {
		t.Fatal(err)
	}
	auction.ClearEvents()
	svc := newTestService(&mockAuctionRepo{auction: auction})

	err := svc.Close(context.Background(), command.Close{
		UserID:    userID,
		AuctionID: auction.ID(),
	})
	if err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestAuctionService_GetByID(t *testing.T) {
	userID := uuid.New().String()
	now := time.Now()
	auction := entity.ReconstructAuction(uuid.New().String(), userID, "Test", "Desc", 500, entity.StatusOpen, now.Add(24*time.Hour), now, now)
	svc := newTestService(&mockAuctionRepo{auction: auction})

	result, err := svc.GetByID(context.Background(), query.Get{
		AuctionID: auction.ID(),
	})
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if result.Title != "Test" {
		t.Errorf("Title = %q, want %q", result.Title, "Test")
	}
	if result.SellerID != userID {
		t.Errorf("SellerID = %q, want %q", result.SellerID, userID)
	}
}

func TestAuctionService_GetList(t *testing.T) {
	userID := uuid.New().String()
	now := time.Now()
	a1 := entity.ReconstructAuction(uuid.New().String(), userID, "Auction 1", "", 100, entity.StatusOpen, now.Add(24*time.Hour), now, now)
	a2 := entity.ReconstructAuction(uuid.New().String(), userID, "Auction 2", "", 200, entity.StatusDraft, now.Add(48*time.Hour), now, now)
	repo := &mockAuctionRepo{auctions: []*entity.Auction{a1, a2}, total: 2}
	svc := newTestService(repo)

	result, err := svc.GetList(context.Background(), query.List{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("GetList() error = %v", err)
	}
	if result.Total != 2 {
		t.Errorf("Total = %d, want 2", result.Total)
	}
	if len(result.Auctions) != 2 {
		t.Errorf("len(Auctions) = %d, want 2", len(result.Auctions))
	}
}

func TestAuctionService_Settle(t *testing.T) {
	userID := uuid.New().String()
	auction, _ := entity.NewAuction(userID, "Test", "", 100, time.Now().Add(2*time.Hour))
	if err := auction.Open(); err != nil {
		t.Fatal(err)
	}
	if err := auction.Close(); err != nil {
		t.Fatal(err)
	}
	auction.ClearEvents()
	svc := newTestService(&mockAuctionRepo{auction: auction})

	err := svc.Settle(context.Background(), command.Settle{AuctionID: auction.ID()})
	if err != nil {
		t.Fatalf("Settle() error = %v", err)
	}
}

func TestAuctionService_Cancel(t *testing.T) {
	userID := uuid.New().String()
	auction, _ := entity.NewAuction(userID, "Test", "", 100, time.Now().Add(2*time.Hour))
	svc := newTestService(&mockAuctionRepo{auction: auction})

	err := svc.Cancel(context.Background(), command.Cancel{AuctionID: auction.ID()})
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
}

func TestAuctionService_Cancel_ByOwner(t *testing.T) {
	ownerID := uuid.New().String()
	auction, _ := entity.NewAuction(ownerID, "Test", "", 100, time.Now().Add(2*time.Hour))
	svc := newTestService(&mockAuctionRepo{auction: auction})

	err := svc.Cancel(context.Background(), command.Cancel{
		UserID:    ownerID,
		AuctionID: auction.ID(),
	})
	if err != nil {
		t.Fatalf("Cancel() by owner error = %v", err)
	}
}

func TestAuctionService_Cancel_NotOwner(t *testing.T) {
	ownerID := uuid.New().String()
	otherID := uuid.New().String()
	auction, _ := entity.NewAuction(ownerID, "Test", "", 100, time.Now().Add(2*time.Hour))
	svc := newTestService(&mockAuctionRepo{auction: auction})

	err := svc.Cancel(context.Background(), command.Cancel{
		UserID:    otherID,
		AuctionID: auction.ID(),
	})
	if err == nil {
		t.Error("expected error for non-owner cancellation")
	}
}

func TestAuctionService_GetEvents(t *testing.T) {
	svc := newTestService(&mockAuctionRepo{})

	result, err := svc.GetEvents(context.Background(), query.EventHistory{AuctionID: uuid.New().String()})
	if err != nil {
		t.Fatalf("GetEvents() error = %v", err)
	}
	if len(result.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(result.Events))
	}
}

func TestAuctionService_GetByID_NotFound(t *testing.T) {
	svc := newTestService(&mockAuctionRepo{auction: nil})

	_, err := svc.GetByID(context.Background(), query.Get{AuctionID: uuid.New().String()})
	if err == nil {
		t.Error("expected error for not found auction")
	}
}

func TestAuctionService_Open_NotFound(t *testing.T) {
	svc := newTestService(&mockAuctionRepo{auction: nil})

	err := svc.Open(context.Background(), command.Open{
		UserID:    uuid.New().String(),
		AuctionID: uuid.New().String(),
	})
	if err == nil {
		t.Error("expected error for not found auction")
	}
}

func TestAuctionService_Close_NotFound(t *testing.T) {
	svc := newTestService(&mockAuctionRepo{auction: nil})

	err := svc.Close(context.Background(), command.Close{
		UserID:    uuid.New().String(),
		AuctionID: uuid.New().String(),
	})
	if err == nil {
		t.Error("expected error for not found auction")
	}
}

func TestAuctionService_Create_EndTimeTooShort(t *testing.T) {
	svc := newTestService(&mockAuctionRepo{})

	_, err := svc.Create(context.Background(), command.Create{
		UserID:     uuid.New().String(),
		Title:      "Test",
		StartPrice: 1000,
		EndTime:    time.Now().Add(30 * time.Minute),
	})
	if err == nil {
		t.Error("expected error for end time too short")
	}
}
