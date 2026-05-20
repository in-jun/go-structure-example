package command

import (
	"context"
	"testing"
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
