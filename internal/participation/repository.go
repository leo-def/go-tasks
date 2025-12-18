package participation

import (
	"go-tasks/internal/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	database.Repository
	GetByCollaboratorId(collaboratorId uuid.UUID) ([]CollaboratorParticipation, int64, error)
	GetByActivityId(activityId uuid.UUID) ([]ActivityParticipation, int64, error)
	GetInfoByActivityIdAndCollaboratorId(activityId uuid.UUID, collaboratorId uuid.UUID) (ParticipationInfo, bool, error)
	GetById(id uuid.UUID) (*Participation, bool, error)
	GetActivityParticipationById(id uuid.UUID, activityId uuid.UUID) (*ActivityParticipation, bool, error)
	Create(participation *ParticipationInfo) error
	CreateTx(tx database.Transaction, participation *ParticipationInfo) error
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

func (r *gormRepository) GetByCollaboratorId(collaboratorId uuid.UUID) ([]CollaboratorParticipation, int64, error) {
	var participations []CollaboratorParticipation
	var count int64
	dbc, _ := database.AsGormConn(r.db)
	q := dbc.Where("collaborator_id = ?", collaboratorId).
		Joins("JOIN activities a ON a.id = participations.activity_id").
		Preload("Activity")
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Find(&participations).Error; err != nil {
		return nil, 0, err
	}
	return participations, count, nil
}

func (r *gormRepository) GetByActivityId(activityId uuid.UUID) ([]ActivityParticipation, int64, error) {
	var participations []ActivityParticipation
	var count int64
	dbc2, _ := database.AsGormConn(r.db)
	q := dbc2.Where("activity_id = ?", activityId).
		Joins("JOIN collaborators c ON c.id = participations.collaborator_id").
		Preload("Collaborator.Account")
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Find(&participations).Error; err != nil {
		return nil, 0, err
	}
	return participations, count, nil
}

func (r *gormRepository) GetById(id uuid.UUID) (*Participation, bool, error) {
	var participation Participation
	db, _ := database.AsGormConn(r.db)
	res := db.
		Preload("Collaborator.Account").
		Preload("Activity").
		First(&participation, "id = ?", id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, res.Error
	}
	return &participation, true, nil
}

func (r *gormRepository) GetActivityParticipationById(id uuid.UUID, activityId uuid.UUID) (*ActivityParticipation, bool, error) {
	var participation ActivityParticipation
	db, _ := database.AsGormConn(r.db)
	res := db.
		Preload("Collaborator.Account").
		First(&participation, "id = ? AND activity_id = ?", id, activityId)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, res.Error
	}
	return &participation, true, nil
}

func (r *gormRepository) GetInfoByActivityIdAndCollaboratorId(activityId uuid.UUID, collaboratorId uuid.UUID) (ParticipationInfo, bool, error) {
	var participation ParticipationInfo
	db, _ := database.AsGormConn(r.db)
	res := db.
		Where("activity_id = ? AND collaborator_id = ?", activityId, collaboratorId).
		First(&participation)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return ParticipationInfo{}, false, nil
		}
		return ParticipationInfo{}, false, res.Error
	}
	return participation, true, nil
}

func (r *gormRepository) Create(participation *ParticipationInfo) error {
	db, _ := database.AsGormConn(r.db)
	return db.Create(participation).Error
}

func (r *gormRepository) CreateTx(tx database.Transaction, participation *ParticipationInfo) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Create(participation).Error
}
