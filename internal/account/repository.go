package account

import (
	"go-tasks/internal/pkg/database"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	database.Repository
	Get(p httpx.PaginationParams) ([]Account, int64, error)
	GetById(id uuid.UUID) (*Account, bool, error)
	FindForSignIn(account *AuthAccount) (bool, error)
	Delete(id uuid.UUID) (bool, error)
	Update(account *Account) error
	Create(account *Account) error
	UpdatePassword(id uuid.UUID, password string) error
	CreateTx(tx database.Transaction, account *Account) error
	UpdatePasswordTx(tx database.Transaction, id uuid.UUID, password string) error
	UpdateEmail(id uuid.UUID, email string) error
	UpdatePhone(id uuid.UUID, phone string) error
}

type gormRepository struct {
	db database.Connection
}

func NewRepository(conn database.Connection) Repository {
	return &gormRepository{db: conn}
}

func (r *gormRepository) GetConnection() database.Connection {
	return r.db
}

func (r *gormRepository) Get(p httpx.PaginationParams) ([]Account, int64, error) {
	dbc, _ := database.AsGormConn(r.db)
	return database.Paginate[Account](dbc, &Account{}, p, func(q *gorm.DB) *gorm.DB {
		if p.Filter != "" {
			return q.Where("username ILIKE ?", "%"+p.Filter+"%")
		}
		return q
	})
}

func (r *gormRepository) GetById(id uuid.UUID) (*Account, bool, error) {
	var account Account
	dbc, _ := database.AsGormConn(r.db)
	res := dbc.First(&account, "id = ?", id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, res.Error
	}
	return &account, true, nil
}

func (r *gormRepository) Delete(id uuid.UUID) (bool, error) {
	dbc, _ := database.AsGormConn(r.db)
	res := dbc.Delete(&Account{}, "id = ?", id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *gormRepository) FindForSignIn(account *AuthAccount) (bool, error) {
	dbc, _ := database.AsGormConn(r.db)
	res := dbc.
		Where("username = ? OR email = ? OR phone = ?", account.Username, account.Email, account.Phone).
		Preload("Roles").
		First(account)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, res.Error
	}
	return true, nil
}

func (r *gormRepository) Update(account *Account) error {
	dbc, _ := database.AsGormConn(r.db)
	return dbc.Model(&Account{}).Where("id = ?", account.ID).Updates(account).Error
}

func (r *gormRepository) Create(account *Account) error {
	dbc, _ := database.AsGormConn(r.db)
	return dbc.Create(account).Error
}

func (r *gormRepository) UpdatePassword(id uuid.UUID, password string) error {
	dbc, _ := database.AsGormConn(r.db)
	return dbc.Model(&Account{}).Where("id = ?", id).Update("password", password).Error
}

func (r *gormRepository) CreateTx(tx database.Transaction, account *Account) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Create(account).Error
}

func (r *gormRepository) UpdatePasswordTx(tx database.Transaction, id uuid.UUID, password string) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Model(&Account{}).Where("id = ?", id).Update("password", password).Error
}

func (r *gormRepository) UpdateEmail(id uuid.UUID, email string) error {
	dbc, _ := database.AsGormConn(r.db)
	return dbc.Model(&Account{}).Where("id = ?", id).Update("email", email).Error
}

func (r *gormRepository) UpdatePhone(id uuid.UUID, phone string) error {
	dbc, _ := database.AsGormConn(r.db)
	return dbc.Model(&Account{}).Where("id = ?", id).Update("phone", phone).Error
}
