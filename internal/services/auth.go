package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"atp-services/internal/models"
	"atp-services/internal/ports"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users    ports.UserRepository
	sessions ports.SessionRepository
}

func NewAuthService(uow ports.UnitOfWork) *AuthService {
	return &AuthService{users: uow, sessions: uow}
}

func (a *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	u, err := a.users.FindUserByLogin(req.Login)
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
	if err := a.sessions.SaveSession(sess); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return &models.LoginResponse{Token: token, User: *u}, nil
}

func (a *AuthService) Logout(token string) error {
	return a.sessions.DeleteSession(token)
}

func (a *AuthService) Validate(token string) (*models.User, error) {
	sess, err := a.sessions.FindSession(token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	u, err := a.users.FindUserByID(sess.UserID)
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
