package assignment

import (
	"go-tasks/internal/pkg/database"
	"go-tasks/internal/pkg/httpx"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	database.Repository
	CheckByIdAndTaskId(id uuid.UUID, taskId uuid.UUID) (bool, error)
	Create(assignment *Assignment) error
	CreateTx(tx database.Transaction, assignment *Assignment) error
	Deactivate(id uuid.UUID) (bool, error)
	DeactivateByTaskId(taskId uuid.UUID) (int64, error)
	DeactivateByTaskIdTx(tx database.Transaction, taskId uuid.UUID) (int64, error)
	Delete(id uuid.UUID) (bool, error)
	DeleteByAssignerId(assignerId uuid.UUID) (int64, error)
	DeleteByAssignerIdTx(tx database.Transaction, assignerId uuid.UUID) (int64, error)
	GetById(id uuid.UUID) (*Assignment, bool, error)
	GetByTaskAndCollaboratorId(taskId, collaboratorId uuid.UUID) (*TaskAssignment, bool, error)
	GetByTaskId(taskId uuid.UUID, p httpx.PaginationParams) ([]TaskAssignment, int64, error)
	GetInfoById(id uuid.UUID) (*AssignmentInfo, bool, error)
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

func (r *gormRepository) GetByTaskId(taskId uuid.UUID, p httpx.PaginationParams) ([]TaskAssignment, int64, error) {
	dbc, _ := database.AsGormConn(r.db)
	return database.Paginate[TaskAssignment](dbc, &TaskAssignment{}, p, func(q *gorm.DB) *gorm.DB {
		return q.
			Where("task_id = ?", taskId).
			Preload("Collaborator.Account")
	})
}

func (r *gormRepository) GetByTaskAndCollaboratorId(taskId, collaboratorId uuid.UUID) (*TaskAssignment, bool, error) {
	var assignment TaskAssignment
	db, _ := database.AsGormConn(r.db)
	res := db.
		Where("task_id = ? AND collaborator_id = ?", taskId, collaboratorId).
		Preload("Collaborator.Account").
		Preload("Task.Lifecycle").
		First(&assignment)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, res.Error
	}
	return &assignment, true, nil
}

func (r *gormRepository) GetInfoById(id uuid.UUID) (*AssignmentInfo, bool, error) {
	var assignment AssignmentInfo
	db, _ := database.AsGormConn(r.db)
	res := db.
		First(&assignment, "id = ?", id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, res.Error
	}
	return &assignment, true, nil
}

func (r *gormRepository) CheckByIdAndTaskId(id uuid.UUID, taskId uuid.UUID) (bool, error) {
	var exists bool
	db, _ := database.AsGormConn(r.db)
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM assignments WHERE id = ? AND task_id = ?)", id, taskId).Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *gormRepository) GetById(id uuid.UUID) (*Assignment, bool, error) {
	var assignment Assignment
	db, _ := database.AsGormConn(r.db)
	res := db.
		Preload("Collaborator.Account").
		Preload("Task.Lifecycle").
		First(&assignment, "id = ?", id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, res.Error
	}
	return &assignment, true, nil
}

func (r *gormRepository) Delete(id uuid.UUID) (bool, error) {
	db, _ := database.AsGormConn(r.db)
	res := db.Delete(&Assignment{}, "id = ?", id)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *gormRepository) Deactivate(id uuid.UUID) (bool, error) {
	db, _ := database.AsGormConn(r.db)
	res := db.Model(&Assignment{}).Where("id = ?", id).Update("active", false)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *gormRepository) Create(assignment *Assignment) error {
	db, _ := database.AsGormConn(r.db)
	return db.Create(assignment).Error
}

func (r *gormRepository) CreateTx(tx database.Transaction, assignment *Assignment) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Create(assignment).Error
}

func (r *gormRepository) DeleteByAssignerId(assignerId uuid.UUID) (int64, error) {
	db, _ := database.AsGormConn(r.db)
	res := db.Where("assigner_id = ?", assignerId).Delete(&Assignment{})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (r *gormRepository) DeleteByAssignerIdTx(tx database.Transaction, assignerId uuid.UUID) (int64, error) {
	txdb, _ := database.AsGormTx(tx)
	res := txdb.Where("assigner_id = ?", assignerId).Delete(&Assignment{})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (r *gormRepository) DeactivateByTaskId(taskId uuid.UUID) (int64, error) {
	db, _ := database.AsGormConn(r.db)
	res := db.Model(&Assignment{}).Where("task_id = ? AND active = ?", taskId, true).Update("active", false)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (r *gormRepository) DeactivateByTaskIdTx(tx database.Transaction, taskId uuid.UUID) (int64, error) {
	txdb, _ := database.AsGormTx(tx)
	res := txdb.Model(&Assignment{}).Where("task_id = ? AND active = ?", taskId, true).Update("active", false)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
