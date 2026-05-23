package app

import (
	"testing"

	"atp-services/internal/models"
)

func TestCreateUserRoles(t *testing.T) {
	core, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer core.Close()

	admin, err := core.Login(models.LoginRequest{Login: "admin", Password: "admin123"})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		role models.Role
		login string
	}{
		{models.RoleDriver, "newdriver"},
		{models.RoleDispatcher, "newdisp"},
		{models.RoleAdmin, "newadmin"},
	}

	for _, c := range cases {
		u, err := core.CreateUser(admin.Token, models.User{
			Login: c.login, Role: c.role,
			FirstName: "Тест", LastName: "Пользователь", Phone: "+7000",
		}, "pass1234")
		if err != nil {
			t.Fatalf("create %s: %v", c.role, err)
		}
		if u.Role != c.role || u.Login != c.login {
			t.Fatalf("unexpected user: %+v", u)
		}
	}

	_, err = core.CreateUser(admin.Token, models.User{
		Login: "newdriver", Role: models.RoleDriver,
		FirstName: "X", LastName: "Y",
	}, "pass1234")
	if err == nil {
		t.Fatal("expected duplicate login error")
	}
}

func TestCreateUserForbiddenForDispatcher(t *testing.T) {
	core, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer core.Close()

	disp, _ := core.Login(models.LoginRequest{Login: "dispatcher", Password: "disp123"})
	_, err = core.CreateUser(disp.Token, models.User{
		Login: "x", Role: models.RoleDriver, FirstName: "A", LastName: "B",
	}, "pass1234")
	if err == nil {
		t.Fatal("dispatcher must not create users")
	}
}

func TestDeleteUser(t *testing.T) {
	core, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer core.Close()

	admin, err := core.Login(models.LoginRequest{Login: "admin", Password: "admin123"})
	if err != nil {
		t.Fatal(err)
	}

	created, err := core.CreateUser(admin.Token, models.User{
		Login: "todelete", Role: models.RoleDriver,
		FirstName: "Удал", LastName: "Яемый",
	}, "pass1234")
	if err != nil {
		t.Fatal(err)
	}

	if err := core.DeleteUser(admin.Token, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	users, err := core.ListUsers(admin.Token)
	if err != nil {
		t.Fatal(err)
	}
	for _, u := range users {
		if u.ID == created.ID {
			t.Fatal("deleted user still in list")
		}
	}
}

func TestDeleteUserCannotDeleteSelf(t *testing.T) {
	core, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer core.Close()

	admin, err := core.Login(models.LoginRequest{Login: "admin", Password: "admin123"})
	if err != nil {
		t.Fatal(err)
	}

	if err := core.DeleteUser(admin.Token, admin.User.ID); err == nil {
		t.Fatal("expected error when deleting self")
	}
}

func TestDeleteUserCannotDeleteLastActiveAdmin(t *testing.T) {
	core, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer core.Close()

	admin, err := core.Login(models.LoginRequest{Login: "admin", Password: "admin123"})
	if err != nil {
		t.Fatal(err)
	}

	second, err := core.CreateUser(admin.Token, models.User{
		Login: "admin2", Role: models.RoleAdmin,
		FirstName: "Второй", LastName: "Админ",
	}, "pass1234")
	if err != nil {
		t.Fatal(err)
	}
	if err := core.DeleteUser(admin.Token, second.ID); err != nil {
		t.Fatalf("delete second admin: %v", err)
	}

	users, err := core.ListUsers(admin.Token)
	if err != nil {
		t.Fatal(err)
	}
	admins := 0
	for _, u := range users {
		if u.Role == models.RoleAdmin && u.Active {
			admins++
		}
	}
	if admins != 1 {
		t.Fatalf("expected one admin left, got %d", admins)
	}

	// Единственный оставшийся админ — не может удалить себя (проверка «последнего» через self).
	if err := core.DeleteUser(admin.Token, admin.User.ID); err == nil {
		t.Fatal("expected error when sole admin deletes self")
	}
}
