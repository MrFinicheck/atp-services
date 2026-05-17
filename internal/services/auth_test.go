package services

import (
	"path/filepath"
	"testing"

	"atp-services/internal/models"
	"atp-services/internal/store"
)

func TestLoginDemoUser(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	auth := NewAuthService(st)
	if err := SeedDemoData(st, auth); err != nil {
		t.Fatal(err)
	}

	resp, err := auth.Login(models.LoginRequest{Login: "dispatcher", Password: "disp123"})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("empty token")
	}
	if resp.User.Login != "dispatcher" {
		t.Fatalf("unexpected user: %s", resp.User.Login)
	}
}

func TestFindUserByLogin(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested")
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	auth := NewAuthService(st)
	_ = SeedDemoData(st, auth)
	u, err := st.FindUserByLogin("admin")
	if err != nil {
		t.Fatal(err)
	}
	if u.Login != "admin" {
		t.Fatal(u.Login)
	}
}
