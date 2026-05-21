package grpc

import (
	"context"
	stderrors "errors"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/in-jun/go-structure-example/internal/auction/application"
	"github.com/in-jun/go-structure-example/internal/auction/application/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	auctionv1 "github.com/in-jun/go-structure-example/proto/auction/v1"
)

var _ auctionv1.AuctionServiceServer = (*server)(nil)

type server struct {
	auctionv1.UnimplementedAuctionServiceServer
	queries application.QueryUseCase
}

func (s *server) GetAuction(ctx context.Context, req *auctionv1.GetAuctionRequest) (*auctionv1.GetAuctionResponse, error) {
	result, err := s.queries.GetByID(ctx, query.Get{AuctionID: req.AuctionId})
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &auctionv1.GetAuctionResponse{
		Id:         result.ID,
		SellerId:   result.SellerID,
		StartPrice: result.StartPrice,
		Status:     result.Status,
	}, nil
}

func toGRPCError(err error) error {
	switch {
	case stderrors.Is(err, errors.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case stderrors.Is(err, errors.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	case stderrors.Is(err, errors.ErrConflict):
		return status.Error(codes.Aborted, err.Error())
	case stderrors.Is(err, errors.ErrBadRequest):
		return status.Error(codes.InvalidArgument, err.Error())
	case stderrors.Is(err, errors.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		slog.Error("unhandled gRPC error", "error", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}
}

func StartGRPCServer(port string, queries application.QueryUseCase) (func(), error) {
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			recoveryInterceptor(),
			loggingInterceptor(),
		),
	)

	auctionv1.RegisterAuctionServiceServer(grpcServer, &server{queries: queries})

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("auction.v1.AuctionService", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	go func() {
		slog.Info("gRPC server starting", "port", port)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed", "error", err)
		}
	}()

	return func() {
		healthServer.SetServingStatus("auction.v1.AuctionService", healthpb.HealthCheckResponse_NOT_SERVING)
		grpcServer.GracefulStop()
		slog.Info("gRPC server stopped")
	}, nil
}
