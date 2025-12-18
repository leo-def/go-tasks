package lifecycle

import (
	"go-tasks/internal/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	database.Repository
	Create(lifecycleInfo *LifecycleInfo) error
	CreateTx(tx database.Transaction, lifecycleInfo *LifecycleInfo) error
	Update(lifecycleInfo *LifecycleInfo) error
	UpdateTx(tx database.Transaction, lifecycleInfo *LifecycleInfo) error
	Delete(uuid uuid.UUID) error
	DeleteTx(tx database.Transaction, uuid uuid.UUID) error
	UpdateStatus(id uuid.UUID, info *LifecycleInfo, status LifecycleStatus) error
	UpdateStatusTx(tx database.Transaction, id uuid.UUID, info *LifecycleInfo, status LifecycleStatus) error
	GetInfoById(id uuid.UUID) (LifecycleInfo, bool, error)
	GetById(id uuid.UUID) (Lifecycle, bool, error)
	GetByIdWithChildrenStatusIn(id uuid.UUID, status []LifecycleStatus) (Lifecycle, bool, error)
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

func (r *gormRepository) Create(lifecycleInfo *LifecycleInfo) error {
	db, _ := database.AsGormConn(r.db)
	return db.Create(lifecycleInfo).Error
}
func (r *gormRepository) CreateTx(tx database.Transaction, lifecycleInfo *LifecycleInfo) error {
	db, _ := database.AsGormTx(tx)
	return db.Create(lifecycleInfo).Error
}

func (r *gormRepository) Update(lifecycleInfo *LifecycleInfo) error {
	db, _ := database.AsGormConn(r.db)
	return db.Model(&LifecycleInfo{}).Where("id = ?", lifecycleInfo.ID).Updates(lifecycleInfo).Error
}

func (r *gormRepository) UpdateTx(tx database.Transaction, lifecycleInfo *LifecycleInfo) error {
	db, _ := database.AsGormTx(tx)
	return db.Model(&LifecycleInfo{}).Where("id = ?", lifecycleInfo.ID).Updates(lifecycleInfo).Error
}

func (r *gormRepository) Delete(uuid uuid.UUID) error {
	db, _ := database.AsGormConn(r.db)
	if err := db.Delete(&LifecycleStatusUpdate{}, "lifecycle_id = ?", uuid).Error; err != nil {
		return err
	}
	return db.Delete(&LifecycleInfo{}, uuid).Error
}

func (r *gormRepository) DeleteTx(tx database.Transaction, uuid uuid.UUID) error {
	db, _ := database.AsGormTx(tx)
	if err := db.Delete(&LifecycleStatusUpdate{}, "lifecycle_id = ?", uuid).Error; err != nil {
		return err
	}
	return db.Delete(&LifecycleInfo{}, uuid).Error
}

func (r *gormRepository) UpdateStatus(id uuid.UUID, info *LifecycleInfo, status LifecycleStatus) error {
	db, _ := database.AsGormConn(r.db)
	var lifecycleStatusUpdate = LifecycleStatusUpdate{
		LifecycleID:  info.ID,
		StatusBefore: info.Status,
		StatusAfter:  status,
	}
	if err := db.Create(&lifecycleStatusUpdate).Error; err != nil {
		return err
	}
	if err := db.Model(&info).Where("id = ?", id.String()).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
func (r *gormRepository) UpdateStatusTx(tx database.Transaction, id uuid.UUID, info *LifecycleInfo, status LifecycleStatus) error {
	db, _ := database.AsGormTx(tx)
	var lifecycleStatusUpdate = LifecycleStatusUpdate{
		LifecycleID:  info.ID,
		StatusBefore: info.Status,
		StatusAfter:  status,
	}
	if err := db.Create(&lifecycleStatusUpdate).Error; err != nil {
		return err
	}

	if err := db.Model(&info).Where("id = ?", id.String()).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func (r *gormRepository) GetInfoById(id uuid.UUID) (LifecycleInfo, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var lifecycle LifecycleInfo
	if err := db.Where("id = ?", id.String()).First(&lifecycle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return LifecycleInfo{}, false, ErrLifecycleNotFound
		}
		return LifecycleInfo{}, false, err
	}
	return lifecycle, true, nil
}

func (r *gormRepository) GetById(id uuid.UUID) (Lifecycle, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var lifecycle Lifecycle
	if err := db.Where("id = ?", id.String()).Preload("Parent").Preload("Children").First(&lifecycle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Lifecycle{}, false, ErrLifecycleNotFound
		}
		return Lifecycle{}, false, err
	}
	return lifecycle, true, nil
}

func (r *gormRepository) GetByIdWithChildrenStatusIn(id uuid.UUID, status []LifecycleStatus) (Lifecycle, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var lifecycle Lifecycle
	if err := db.Where("id = ?", id.String()).
		Joins("LEFT JOIN lifecycles AS children ON children.parent_id = lifecycles.id AND children.status IN ?", status).
		Preload("Children").
		First(&lifecycle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return Lifecycle{}, false, ErrLifecycleNotFound
		}
		return Lifecycle{}, false, err
	}
	return lifecycle, true, nil
}
