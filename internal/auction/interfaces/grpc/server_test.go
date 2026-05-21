package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
)


func TestToGRPCError_NotFound(t *testing.T) {
	err := toGRPCError(errors.NotFound("resource not found"))
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.NotFound {
		t.Errorf("expected NotFound, got %v", st.Code())
	}
}

func TestToGRPCError_Forbidden(t *testing.T) {
	err := toGRPCError(errors.Forbidden("access denied"))
	st, _ := status.FromError(err)
	if st.Code() != codes.PermissionDenied {
		t.Errorf("expected PermissionDenied, got %v", st.Code())
	}
}

func TestToGRPCError_Conflict(t *testing.T) {
	err := toGRPCError(errors.Conflict("already exists"))
	st, _ := status.FromError(err)
	if st.Code() != codes.Aborted {
		t.Errorf("expected Aborted, got %v", st.Code())
	}
}

func TestToGRPCError_BadRequest(t *testing.T) {
	err := toGRPCError(errors.BadRequest("invalid input"))
	st, _ := status.FromError(err)
	if st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument, got %v", st.Code())
	}
}

func TestToGRPCError_Unauthorized(t *testing.T) {
	err := toGRPCError(errors.Unauthorized("not authenticated"))
	st, _ := status.FromError(err)
	if st.Code() != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated, got %v", st.Code())
	}
}

func TestToGRPCError_Internal(t *testing.T) {
	err := toGRPCError(errors.Internal("db error"))
	st, _ := status.FromError(err)
	if st.Code() != codes.Internal {
		t.Errorf("expected Internal, got %v", st.Code())
	}
}

type mockQueryUseCase struct {
	result interface{}
	err    error
}

func (m *mockQueryUseCase) GetByID(ctx context.Context, qry interface{}) (interface{}, error) {
	return m.result, m.err
}

func TestRecoveryInterceptor_NoPanic(t *testing.T) {
	interceptor := recoveryInterceptor()
	handler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/auction.v1.AuctionService/GetAuction"}
	resp, err := interceptor(context.Background(), nil, info, handler)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != "ok" {
		t.Errorf("expected resp 'ok', got %v", resp)
	}
}

func TestRecoveryInterceptor_Panic(t *testing.T) {
	interceptor := recoveryInterceptor()
	handler := func(ctx context.Context, req any) (any, error) {
		panic("unexpected panic")
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/auction.v1.AuctionService/GetAuction"}
	_, err := interceptor(context.Background(), nil, info, handler)
	if err == nil {
		t.Fatal("expected error after panic, got nil")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("expected gRPC status error")
	}
	if st.Code() != codes.Internal {
		t.Errorf("expected Internal after panic, got %v", st.Code())
	}
}
