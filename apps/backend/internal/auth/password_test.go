package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	plain := "P@ssw0rd!"
	hash, err := HashPassword(plain)
	if err != nil {
		t.Fatalf("expected hash success, got err=%v", err)
	}
	if hash == plain {
		t.Fatalf("hash should not equal plain text")
	}

	if !CheckPassword(hash, plain) {
		t.Fatalf("expected password check to pass")
	}

	if CheckPassword(hash, "wrong-password") {
		t.Fatalf("expected password check to fail on wrong password")
	}
}
