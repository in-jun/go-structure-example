package command

import (
	"context"
	"testing"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func TestUpdatePasswordHandler_Success(t *testing.T) {
	h := NewUpdatePasswordHandler(&mockUserRepo{user: makeUser()}, &mockHasher{})

	err := h.Handle(context.Background(), UpdatePassword{
		UserID:          1,
		CurrentPassword: "oldpass",
		NewPassword:     "newpassword123",
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestUpdatePasswordHandler_WrongCurrentPassword(t *testing.T) {
	h := NewUpdatePasswordHandler(&mockUserRepo{user: makeUser()}, &mockHasher{})

	err := h.Handle(context.Background(), UpdatePassword{
		UserID:          1,
		CurrentPassword: "wrongpassword",
		NewPassword:     "newpassword123",
	})
	if err == nil {
		t.Fatal("expected error for wrong current password, got nil")
	}
}

func TestUpdatePasswordHandler_SamePassword(t *testing.T) {
	h := NewUpdatePasswordHandler(&mockUserRepo{user: makeUser()}, &mockHasher{})

	err := h.Handle(context.Background(), UpdatePassword{
		UserID:          1,
		CurrentPassword: "oldpass",
		NewPassword:     "oldpass",
	})
	if err == nil {
		t.Fatal("expected error for same password, got nil")
	}
}

func TestUpdatePasswordHandler_NotFound(t *testing.T) {
	h := NewUpdatePasswordHandler(&mockUserRepo{err: errors.NotFound("user not found")}, &mockHasher{})

	err := h.Handle(context.Background(), UpdatePassword{
		UserID:          99,
		CurrentPassword: "oldpass",
		NewPassword:     "newpassword123",
	})
	if err == nil {
		t.Fatal("expected error for user not found, got nil")
	}
}

func TestUpdatePasswordHandler_EmptyNewPassword(t *testing.T) {
	h := NewUpdatePasswordHandler(&mockUserRepo{user: makeUser()}, &mockHasher{})

	err := h.Handle(context.Background(), UpdatePassword{
		UserID:          1,
		CurrentPassword: "oldpass",
		NewPassword:     "",
	})
	if err == nil {
		t.Fatal("expected error for empty new password, got nil")
	}
}
