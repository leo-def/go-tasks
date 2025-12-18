package company

import (
	"go-tasks/internal/pkg/database"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	database.Repository
	Create(company *Company) error
	CreateTx(tx database.Transaction, company *Company) error
	Delete(id uuid.UUID) (bool, error)
	Get(p httpx.PaginationParams) ([]Company, int64, error)
	GetById(id uuid.UUID) (*Company, bool, error)
	Update(company *Company) error
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

func (r *gormRepository) Create(company *Company) error {
	db, _ := database.AsGormConn(r.db)
	return db.Create(company).Error
}

func (r *gormRepository) CreateTx(tx database.Transaction, company *Company) error {
	db, _ := database.AsGormTx(tx)
	return db.Create(company).Error
}

func (r *gormRepository) Delete(id uuid.UUID) (bool, error) {
	db, _ := database.AsGormConn(r.db)
	res := db.Delete(&Company{}, "id = ?", id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *gormRepository) Get(p httpx.PaginationParams) ([]Company, int64, error) {
	db, _ := database.AsGormConn(r.db)
	return database.Paginate[Company](db, &Company{}, p, func(q *gorm.DB) *gorm.DB {
		if p.Filter != "" {
			q = q.Where("companies.title ILIKE ?", "%"+p.Filter+"%")
		}
		return q.Preload("Owner").Preload("Owner.Account")
	})
}

func (r *gormRepository) GetById(id uuid.UUID) (*Company, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var company Company
	res := db.Preload("Owner").Preload("Owner.Account").First(&company, "id = ?", id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, res.Error
	}
	return &company, true, nil
}

func (r *gormRepository) Update(company *Company) error {
	db, _ := database.AsGormConn(r.db)
	return db.Model(&Company{}).Where("id = ?", company.ID).Updates(company).Error
}
