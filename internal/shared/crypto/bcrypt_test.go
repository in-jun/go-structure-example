package crypto

import "testing"

func TestBcryptHasher_HashAndCompare(t *testing.T) {
	h := NewBcryptHasher()
	password := "test-password"

	hashed, err := h.Hash(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hashed == "" {
		t.Fatal("expected non-empty hash")
	}
	if hashed == password {
		t.Fatal("hash should differ from plaintext")
	}

	if !h.Compare(hashed, password) {
		t.Error("expected Compare to return true for correct password")
	}
	if h.Compare(hashed, "wrong-password") {
		t.Error("expected Compare to return false for wrong password")
	}
}

func TestBcryptHasher_DifferentHashes(t *testing.T) {
	h := NewBcryptHasher()
	h1, _ := h.Hash("password")
	h2, _ := h.Hash("password")
	if h1 == h2 {
		t.Error("expected different hashes for same password (different salts)")
	}
}
