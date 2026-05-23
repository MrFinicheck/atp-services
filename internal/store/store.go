package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"atp-services/internal/models"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Store struct {
	db      *leveldb.DB
	dirLock *os.File
}

func (s *Store) Close() error {
	var err error
	if s.db != nil {
		err = s.db.Close()
		s.db = nil
	}
	releaseDirLock(s.dirLock)
	s.dirLock = nil
	return err
}

func (s *Store) put(prefix, id string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.db.Put([]byte(prefix+id), b, nil)
}

func (s *Store) get(prefix, id string, v any) error {
	b, err := s.db.Get([]byte(prefix+id), nil)
	if err == leveldb.ErrNotFound {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (s *Store) delete(prefix, id string) error {
	return s.db.Delete([]byte(prefix+id), nil)
}

func (s *Store) list(prefix string, dest func(key string, raw []byte) error) error {
	it := s.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer it.Release()
	for it.Next() {
		key := string(it.Key())
		if strings.Contains(key, ":") && strings.Count(key, ":") > 1 {
			continue
		}
		id := strings.TrimPrefix(key, prefix)
		if id == "" {
			continue
		}
		if err := dest(id, it.Value()); err != nil {
			return err
		}
	}
	return it.Error()
}

func newID() string {
	return uuid.New().String()
}

// Users

func (s *Store) SaveUser(u *models.User) error {
	if u.ID == "" {
		u.ID = newID()
	}
	if err := s.put("user:", u.ID, u); err != nil {
		return err
	}
	return s.db.Put([]byte("user:login:"+u.Login), []byte(u.ID), nil)
}

func (s *Store) FindUserByLogin(login string) (*models.User, error) {
	id, err := s.db.Get([]byte("user:login:"+login), nil)
	if err == leveldb.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var u models.User
	if err := s.get("user:", string(id), &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) FindUserByID(id string) (*models.User, error) {
	var u models.User
	if err := s.get("user:", id, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) ListUsers() ([]models.User, error) {
	var out []models.User
	err := s.list("user:", func(id string, raw []byte) error {
		if strings.HasPrefix(id, "login:") {
			return nil
		}
		var u models.User
		if err := json.Unmarshal(raw, &u); err != nil {
			return err
		}
		out = append(out, u)
		return nil
	})
	return out, err
}

func (s *Store) DeleteUser(id string) error {
	u, err := s.FindUserByID(id)
	if err != nil {
		return err
	}
	if err := s.delete("user:", id); err != nil {
		return err
	}
	return s.db.Delete([]byte("user:login:"+u.Login), nil)
}

// Sessions

func (s *Store) SaveSession(sess *models.Session) error {
	return s.put("session:", sess.Token, sess)
}

func (s *Store) FindSession(token string) (*models.Session, error) {
	var sess models.Session
	if err := s.get("session:", token, &sess); err != nil {
		return nil, err
	}
	if time.Now().After(sess.ExpiresAt) {
		_ = s.delete("session:", token)
		return nil, ErrNotFound
	}
	return &sess, nil
}

func (s *Store) DeleteSession(token string) error {
	return s.delete("session:", token)
}

// Clients

func (s *Store) SaveClient(c *models.Client) error {
	if c.ID == "" {
		c.ID = newID()
	}
	return s.put("client:", c.ID, c)
}

func (s *Store) ListClients() ([]models.Client, error) {
	var out []models.Client
	err := s.list("client:", func(_ string, raw []byte) error {
		var c models.Client
		if err := json.Unmarshal(raw, &c); err != nil {
			return err
		}
		out = append(out, c)
		return nil
	})
	return out, err
}

// Vehicles

func (s *Store) SaveVehicle(v *models.Vehicle) error {
	if v.ID == "" {
		v.ID = newID()
	}
	return s.put("vehicle:", v.ID, v)
}

func (s *Store) FindVehicle(id string) (*models.Vehicle, error) {
	var v models.Vehicle
	if err := s.get("vehicle:", id, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *Store) ListVehicles() ([]models.Vehicle, error) {
	var out []models.Vehicle
	err := s.list("vehicle:", func(_ string, raw []byte) error {
		var v models.Vehicle
		if err := json.Unmarshal(raw, &v); err != nil {
			return err
		}
		out = append(out, v)
		return nil
	})
	return out, err
}

// Tariffs

func (s *Store) SaveTariff(t *models.Tariff) error {
	if t.ID == "" {
		t.ID = newID()
	}
	return s.put("tariff:", t.ID, t)
}

func (s *Store) FindTariff(id string) (*models.Tariff, error) {
	var t models.Tariff
	if err := s.get("tariff:", id, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Store) ListTariffs() ([]models.Tariff, error) {
	var out []models.Tariff
	err := s.list("tariff:", func(_ string, raw []byte) error {
		var t models.Tariff
		if err := json.Unmarshal(raw, &t); err != nil {
			return err
		}
		out = append(out, t)
		return nil
	})
	return out, err
}

// Orders

func (s *Store) SaveOrder(o *models.Order) error {
	if o.ID == "" {
		o.ID = newID()
		o.CreatedAt = time.Now()
	}
	return s.put("order:", o.ID, o)
}

func (s *Store) FindOrder(id string) (*models.Order, error) {
	var o models.Order
	if err := s.get("order:", id, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *Store) ListOrders() ([]models.Order, error) {
	var out []models.Order
	err := s.list("order:", func(_ string, raw []byte) error {
		var o models.Order
		if err := json.Unmarshal(raw, &o); err != nil {
			return err
		}
		out = append(out, o)
		return nil
	})
	return out, err
}

// Waybills

func (s *Store) SaveWaybill(w *models.Waybill) error {
	if w.ID == "" {
		w.ID = newID()
		w.CreatedAt = time.Now()
	}
	return s.put("waybill:", w.ID, w)
}

func (s *Store) ListWaybills() ([]models.Waybill, error) {
	var out []models.Waybill
	err := s.list("waybill:", func(_ string, raw []byte) error {
		var w models.Waybill
		if err := json.Unmarshal(raw, &w); err != nil {
			return err
		}
		out = append(out, w)
		return nil
	})
	return out, err
}

func (s *Store) FindOpenWaybill(driverID, date string) (*models.Waybill, error) {
	wb, err := s.FindWaybillByDriverAndDate(driverID, date)
	if err != nil {
		return nil, err
	}
	if wb.Closed {
		return nil, ErrNotFound
	}
	return wb, nil
}

func (s *Store) FindWaybillByDriverAndDate(driverID, date string) (*models.Waybill, error) {
	all, err := s.ListWaybills()
	if err != nil {
		return nil, err
	}
	var open *models.Waybill
	for i := range all {
		w := all[i]
		if w.DriverID != driverID || w.Date != date {
			continue
		}
		if w.Closed {
			return &w, nil
		}
		if open == nil {
			open = &w
		}
	}
	if open != nil {
		return open, nil
	}
	return nil, ErrNotFound
}

// Audit

func (s *Store) AddAudit(userID, action, entityType, entityID, details string) error {
	e := models.AuditEntry{
		ID:         newID(),
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    details,
		CreatedAt:  time.Now(),
	}
	return s.put("audit:", e.ID, &e)
}

func (s *Store) ListAudit(limit int) ([]models.AuditEntry, error) {
	all, err := s.listAudit()
	if err != nil {
		return nil, err
	}
	if limit > 0 && len(all) > limit {
		return all[len(all)-limit:], nil
	}
	return all, nil
}

func (s *Store) listAudit() ([]models.AuditEntry, error) {
	var out []models.AuditEntry
	err := s.list("audit:", func(_ string, raw []byte) error {
		var e models.AuditEntry
		if err := json.Unmarshal(raw, &e); err != nil {
			return err
		}
		out = append(out, e)
		return nil
	})
	return out, err
}

func (s *Store) IsSeeded() (bool, error) {
	_, err := s.db.Get([]byte("meta:seeded"), nil)
	if err == leveldb.ErrNotFound {
		return false, nil
	}
	return err == nil, err
}

func (s *Store) MarkSeeded() error {
	return s.db.Put([]byte("meta:seeded"), []byte("1"), nil)
}

func (s *Store) DataDir() string {
	return ""
}

func (s *Store) String() string {
	return fmt.Sprintf("leveldb-store")
}
