package activity

import (
	"go-tasks/internal/pkg/database"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	database.Repository
	GetByOwnerId(p httpx.PaginationParams, ownerID uuid.UUID) ([]OwnActivity, int64, error)
	GetByCompanyId(p httpx.PaginationParams, companyID uuid.UUID) ([]CompanyActivity, int64, error)
	Get(p httpx.PaginationParams) ([]Activity, int64, error)
	GetById(id uuid.UUID) (*Activity, bool, error)
	CheckByIdAndCompanyId(id uuid.UUID, companyID uuid.UUID) (bool, error)
	CheckByIdAndOwnerId(id uuid.UUID, ownerID uuid.UUID) (bool, error)
	GetByIdAndOwnerId(id uuid.UUID, ownerID uuid.UUID) (*Activity, bool, error)
	GetCompanyActivityById(id, companyId uuid.UUID) (*CompanyActivity, bool, error)
	GetOwnActivityById(id, ownId uuid.UUID) (*OwnActivity, bool, error)
	GetActivityInfoById(id uuid.UUID) (*ActivityInfo, bool, error)
	Delete(id uuid.UUID) (bool, error)
	DeleteTx(tx database.Transaction, id uuid.UUID) (bool, error)
	UpdateInfo(info *ActivityInfo) error
	UpdateInfoTx(tx database.Transaction, info *ActivityInfo) error
	CreateInfo(info *ActivityInfo) error
	CreateInfoTx(tx database.Transaction, info *ActivityInfo) error
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

func (r *gormRepository) GetByOwnerId(p httpx.PaginationParams, ownerID uuid.UUID) ([]OwnActivity, int64, error) {
	dbc, _ := database.AsGormConn(r.db)
	return database.Paginate[OwnActivity](dbc, &OwnActivity{}, p, func(q *gorm.DB) *gorm.DB {
		return q.
			Where("owner_id = ?", ownerID).
			Preload("Lifecycle").
			Preload("Lifecycle.Updates").
			Preload("CreatedBy.Account")
	})
}

func (r *gormRepository) GetByCompanyId(p httpx.PaginationParams, companyID uuid.UUID) ([]CompanyActivity, int64, error) {
	dbc2, _ := database.AsGormConn(r.db)
	return database.Paginate[CompanyActivity](dbc2, &CompanyActivity{}, p, func(q *gorm.DB) *gorm.DB {
		return q.
			Where("company_id = ?", companyID).
			Preload("Lifecycle").
			Preload("Lifecycle.Updates").
			Preload("Owner.Account").
			Preload("CreatedBy.Account")
	})
}

func (r *gormRepository) Get(p httpx.PaginationParams) ([]Activity, int64, error) {
	dbc3, _ := database.AsGormConn(r.db)
	return database.Paginate[Activity](dbc3, &Activity{}, p, func(q *gorm.DB) *gorm.DB {
		return q.
			Preload("Company").
			Preload("Lifecycle").
			Preload("Lifecycle.Updates").
			Preload("Owner.Account").
			Preload("CreatedBy.Account")
	})
}

func (r *gormRepository) GetById(id uuid.UUID) (*Activity, bool, error) {
	var activity Activity
	db, _ := database.AsGormConn(r.db)
	err := db.
		Preload("Company").
		Preload("Lifecycle").
		Preload("Lifecycle.Updates").
		Preload("Owner.Account").
		Preload("CreatedBy.Account").
		First(&activity, "activities.id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &activity, true, nil
}

func (r *gormRepository) CheckByIdAndCompanyId(id uuid.UUID, companyID uuid.UUID) (bool, error) {
	var exists bool
	db, _ := database.AsGormConn(r.db)
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM activities WHERE id = ? AND company_id = ?)", id, companyID).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *gormRepository) CheckByIdAndOwnerId(id uuid.UUID, ownerID uuid.UUID) (bool, error) {
	var exists bool
	db, _ := database.AsGormConn(r.db)
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM activities WHERE id = ? AND owner_id = ?)", id, ownerID).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *gormRepository) GetByIdAndOwnerId(id uuid.UUID, ownerID uuid.UUID) (*Activity, bool, error) {
	var activity Activity
	db, _ := database.AsGormConn(r.db)
	err := db.
		Preload("Company").
		Preload("Lifecycle").
		Preload("Lifecycle.Updates").
		Preload("Owner.Account").
		Preload("CreatedBy.Account").
		First(&activity, "activities.id = ? AND owner_id = ?", id, ownerID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &activity, true, nil
}

func (r *gormRepository) GetCompanyActivityById(id, companyId uuid.UUID) (*CompanyActivity, bool, error) {
	var activity CompanyActivity
	db, _ := database.AsGormConn(r.db)
	err := db.
		Preload("Lifecycle").
		Preload("Lifecycle.Updates").
		Preload("Owner.Account").
		Preload("CreatedBy.Account").
		First(&activity, "activities.id = ? AND activities.company_id = ?", id, companyId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &activity, true, nil
}

func (r *gormRepository) GetOwnActivityById(id, ownerID uuid.UUID) (*OwnActivity, bool, error) {
	var activity OwnActivity
	db, _ := database.AsGormConn(r.db)
	err := db.
		Preload("Lifecycle").
		Preload("Lifecycle.Updates").
		Preload("CreatedBy.Account").
		First(&activity, "activities.id = ? AND activities.owner_id = ?", id, ownerID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &activity, true, nil
}

func (r *gormRepository) GetActivityInfoById(id uuid.UUID) (*ActivityInfo, bool, error) {
	var activity ActivityInfo
	db, _ := database.AsGormConn(r.db)
	err := db.Joins("INNER JOIN lifecycles l ON l.id = activities.lifecycle_id").First(&activity, "activities.id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &activity, true, nil
}

func (r *gormRepository) Delete(id uuid.UUID) (bool, error) {
	db, _ := database.AsGormConn(r.db)
	if err := db.Delete(&Activity{}, "activities.id = ?", id).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *gormRepository) DeleteTx(tx database.Transaction, id uuid.UUID) (bool, error) {
	txdb, _ := database.AsGormTx(tx)
	if err := txdb.Delete(&Activity{}, "activities.id = ?", id).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *gormRepository) UpdateInfo(info *ActivityInfo) error {
	db, _ := database.AsGormConn(r.db)
	return db.Model(&ActivityInfo{}).Where("id = ?", info.ID).Select("title", "owner_id").Updates(info).Error
}

func (r *gormRepository) UpdateInfoTx(tx database.Transaction, info *ActivityInfo) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Model(&ActivityInfo{}).Where("id = ?", info.ID).Select("title", "owner_id").Updates(info).Error
}

func (r *gormRepository) CreateInfo(info *ActivityInfo) error {
	db, _ := database.AsGormConn(r.db)
	return db.Create(info).Error
}

func (r *gormRepository) CreateInfoTx(tx database.Transaction, info *ActivityInfo) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Create(info).Error
}
