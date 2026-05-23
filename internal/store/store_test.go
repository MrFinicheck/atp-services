package store

import (
	"testing"
	"time"

	"atp-services/internal/models"
)

func TestSaveAndGetClient(t *testing.T) {
	st, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	c := &models.Client{Name: "ООО Тест", Phone: "+7000", DebtLimit: 50000}
	if err := st.SaveClient(c); err != nil {
		t.Fatal(err)
	}
	if c.ID == "" {
		t.Fatal("SaveClient must assign ID")
	}

	list, err := st.ListClients()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) < 1 {
		t.Fatal("ListClients empty")
	}
	var found bool
	for _, x := range list {
		if x.ID == c.ID && x.Name == c.Name {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("saved client not in list")
	}
}

func TestSaveOrderPersistsJSON(t *testing.T) {
	st, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	o := &models.Order{
		ClientID: "c1", VehicleID: "v1", DriverID: "d1",
		FromAddr: "ул. А", ToAddr: "ул. Б",
		DistanceKm: 12, Price: 1500, Status: models.OrderAssigned,
	}
	if err := st.SaveOrder(o); err != nil {
		t.Fatal(err)
	}
	got, err := st.FindOrder(o.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.FromAddr != "ул. А" || got.Price != 1500 {
		t.Fatalf("unexpected order: %+v", got)
	}
}

func TestSessionExpiry(t *testing.T) {
	st, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	s := &models.Session{
		Token: "tok123", UserID: "u1", Role: models.RoleAdmin,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if err := st.SaveSession(s); err != nil {
		t.Fatal(err)
	}
	found, err := st.FindSession("tok123")
	if err != nil {
		t.Fatal(err)
	}
	if found.UserID != "u1" {
		t.Fatal(found.UserID)
	}
}
