package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"atp-services/internal/models"
	"atp-services/internal/store"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	store *store.Store
}

func NewAuthService(s *store.Store) *AuthService {
	return &AuthService{store: s}
}

func (a *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	u, err := a.store.FindUserByLogin(req.Login)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !u.Active {
		return nil, errors.New("user is inactive")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	token, err := randomToken()
	if err != nil {
		return nil, err
	}
	sess := &models.Session{
		Token:     token,
		UserID:    u.ID,
		Role:      u.Role,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := a.store.SaveSession(sess); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return &models.LoginResponse{Token: token, User: *u}, nil
}

func (a *AuthService) Logout(token string) error {
	return a.store.DeleteSession(token)
}

func (a *AuthService) Validate(token string) (*models.User, error) {
	sess, err := a.store.FindSession(token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	u, err := a.store.FindUserByID(sess.UserID)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	u.PasswordHash = ""
	return u, nil
}

func (a *AuthService) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
