package account

import (
	"go-tasks/internal/pkg/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
    FindForSignIn(Username string, Email string, Phone string) (*AuthAccount, bool, error)
    Create(account *Account) error
    CreateWithPassword(account *Account, password string) error
    CreateWithPasswordTx(tx database.Transaction, account *Account, password string) error
    UpdatePassword(id uuid.UUID, password string) error
    UpdatePasswordTx(tx database.Transaction, id uuid.UUID, password string) error
    VerifyPassword(password, hashedPassword string) error
    UpdateEmail(id uuid.UUID, email string) error
    UpdatePhone(id uuid.UUID, phone string) error
    GetById(id uuid.UUID) (*Account, bool, error)
    Delete(id uuid.UUID) (bool, error)
    EnsureAdmin(username, password, email, name, phone string) error
}

type service struct {
	repository Repository
}

func NewService(Repository Repository) Service {
	return &service{Repository}
}

func (s *service) EnsureAdmin(username, password, email, name, phone string) error {
	_, found, err := s.FindForSignIn(username, email, phone)
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	account := &Account{
		Username: username,
		Name:     name,
		Email:    email,
		Phone:    phone,
		Role:     RoleAdmin,
	}
	return s.CreateWithPassword(account, password)
}

func (s *service) FindForSignIn(Username string, Email string, Phone string) (*AuthAccount, bool, error) {
	account := &AuthAccount{Account: Account{Username: Username, Email: Email, Phone: Phone}, Roles: make([]AccountCollaborator, 0)}
	found, err := s.repository.FindForSignIn(account)
	if err != nil {
		return nil, false, err
	}
	if !found {
		return nil, false, nil
	}
	return account, true, nil
}

func (s *service) Create(account *Account) error {
	return s.repository.Create(account)
}

func (s *service) CreateWithPassword(account *Account, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Password = string(hashedPassword)
	return s.repository.Create(account)
}

func (s *service) CreateWithPasswordTx(tx database.Transaction, account *Account, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Password = string(hashedPassword)
	return s.repository.CreateTx(tx, account)
}

func (s *service) UpdatePassword(id uuid.UUID, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repository.UpdatePassword(id, string(hashedPassword))
}

func (s *service) UpdatePasswordTx(tx database.Transaction, id uuid.UUID, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repository.UpdatePasswordTx(tx, id, string(hashedPassword))
}

func (s *service) VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *service) UpdateEmail(id uuid.UUID, email string) error {
	return s.repository.UpdateEmail(id, email)
}

func (s *service) UpdatePhone(id uuid.UUID, phone string) error {
    return s.repository.UpdatePhone(id, phone)
}

func (s *service) GetById(id uuid.UUID) (*Account, bool, error) {
    return s.repository.GetById(id)
}

func (s *service) Delete(id uuid.UUID) (bool, error) {
    return s.repository.Delete(id)
}
