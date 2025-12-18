package task

import (
	"go-tasks/internal/lifecycle"
	"go-tasks/internal/pkg/database"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
    database.Repository
    CreateInfo(info *TaskInfo) error
    CreateInfoTx(tx database.Transaction, info *TaskInfo) error
    Delete(id uuid.UUID) (bool, error)
    DeleteTx(tx database.Transaction, id uuid.UUID) (bool, error)
    Get(p httpx.PaginationParams) ([]Task, int64, error)
    GetActivityTaskById(id uuid.UUID, activityID uuid.UUID) (*ActivityTask, bool, error)
    GetByActivityId(activityId uuid.UUID, p httpx.PaginationParams) ([]ActivityTask, int64, error)
    GetByActivityIdInStatus(activityId uuid.UUID, status []lifecycle.LifecycleStatus) ([]ActivityTask, int64, error)
    GetById(id uuid.UUID) (*Task, bool, error)
    GetInfoById(id uuid.UUID) (*TaskInfo, bool, error)
    UpdateInfo(info *TaskInfo) error
    UpdateInfoTx(tx database.Transaction, info *TaskInfo) error
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

func (r *gormRepository) Get(p httpx.PaginationParams) ([]Task, int64, error) {
	dbc, _ := database.AsGormConn(r.db)
	return database.Paginate[Task](dbc, &Task{}, p, func(q *gorm.DB) *gorm.DB {
		q = q.Preload("Lifecycle").Preload("Lifecycle.Updates").Preload("Activity")
		if p.Filter != "" {
			q = q.Where("title ILIKE ?", "%"+p.Filter+"%")
		}
		return q
	})
}

func (r *gormRepository) GetByActivityId(activityId uuid.UUID, p httpx.PaginationParams) ([]ActivityTask, int64, error) {
	dbc2, _ := database.AsGormConn(r.db)
	return database.Paginate[ActivityTask](dbc2, &ActivityTask{}, p, func(q *gorm.DB) *gorm.DB {
		q = q.Where("activity_id = ?", activityId.String())
		q = q.Preload("Lifecycle").Preload("Lifecycle.Updates")
		if p.Filter != "" {
			q = q.Where("title ILIKE ?", "%"+p.Filter+"%")
		}
		return q
	})
}

func (r *gormRepository) GetByActivityIdInStatus(activityId uuid.UUID, status []lifecycle.LifecycleStatus) ([]ActivityTask, int64, error) {
	var tasks []ActivityTask
	db, _ := database.AsGormConn(r.db)
	q := db.Where("activity_id = ?", activityId.String())
	q = q.Joins("INNER JOIN lifecycles l ON l.status IN ?", status)
	q = q.Preload("Lifecycle").Preload("Lifecycle.Updates")
	if err := q.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, int64(len(tasks)), nil
}

func (r *gormRepository) GetInfoById(id uuid.UUID) (*TaskInfo, bool, error) {
	var task TaskInfo
	db, _ := database.AsGormConn(r.db)
	q := db.Where("id = ?", id)
	q = q.Preload("Lifecycle").Preload("Activity")
	result := q.First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, result.Error
	}
	return &task, true, nil
}

func (r *gormRepository) GetById(id uuid.UUID) (*Task, bool, error) {
	var task Task
	db, _ := database.AsGormConn(r.db)
	q := db.Where("id = ?", id)
	q = q.
		Preload("Lifecycle").
		Preload("Lifecycle.Updates").
		Preload("Activity").
		Preload("Activity.Lifecycle").
		Preload("Activity.Lifecycle.Updates").
		Preload("Activity.Company").
		Preload("Activity.Owner").
		Preload("Activity.Owner.Account").
		Preload("Activity.CreatedBy").
		Preload("Activity.CreatedBy.Account")
	result := q.First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, result.Error
	}
	return &task, true, nil
}

func (r *gormRepository) GetActivityTaskById(id uuid.UUID, activityID uuid.UUID) (*ActivityTask, bool, error) {
	var task ActivityTask
	db, _ := database.AsGormConn(r.db)
	q := db.Where("id = ? AND activity_id = ?", id, activityID)
	q = q.
		Preload("Lifecycle").
		Preload("Lifecycle.Updates").
		Preload("Activity").
		Preload("Activity.Lifecycle").
		Preload("Activity.Lifecycle.Updates").
		Preload("Activity.Company").
		Preload("Activity.Owner").
		Preload("Activity.Owner.Account").
		Preload("Activity.CreatedBy").
		Preload("Activity.CreatedBy.Account")
	result := q.First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, result.Error
	}
	return &task, true, nil
}

func (r *gormRepository) Delete(id uuid.UUID) (bool, error) {
	db, _ := database.AsGormConn(r.db)
	q := db.Model(&Task{})
	if err := q.Delete(&Task{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *gormRepository) DeleteTx(tx database.Transaction, id uuid.UUID) (bool, error) {
	txdb, _ := database.AsGormTx(tx)
	q := txdb.Model(&Task{})
	if err := q.Delete(&Task{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *gormRepository) UpdateInfo(info *TaskInfo) error {
	db, _ := database.AsGormConn(r.db)
	return db.Model(&TaskInfo{}).Where("id = ?", info.ID).Select("title").Updates(info).Error
}

func (r *gormRepository) UpdateInfoTx(tx database.Transaction, info *TaskInfo) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Model(&TaskInfo{}).Where("id = ?", info.ID).Select("title").Updates(info).Error
}

func (r *gormRepository) CreateInfo(info *TaskInfo) error {
	db, _ := database.AsGormConn(r.db)
	return db.Create(info).Error
}

func (r *gormRepository) CreateInfoTx(tx database.Transaction, info *TaskInfo) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Create(info).Error
}
