package rating

import (
	"go-tasks/internal/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	database.Repository
	Create(rating *Rating) error
	CreateTx(tx database.Transaction, rating *Rating) error
	FindByParticipationAndCollaborator(participationID, collaboratorID uuid.UUID) (*Rating, bool, error)
	FindByParticipationAndCollaboratorTx(tx database.Transaction, participationID, collaboratorID uuid.UUID) (*Rating, bool, error)
	Update(rating *Rating) error
	UpdateTx(tx database.Transaction, rating *Rating) error
	Upsert(rating *Rating) error
	UpsertTx(tx database.Transaction, rating *Rating) error
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

func (r *gormRepository) Create(rating *Rating) error {
	db, _ := database.AsGormConn(r.db)
	return db.Create(rating).Error
}

func (r *gormRepository) CreateTx(tx database.Transaction, rating *Rating) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Create(rating).Error
}

func (r *gormRepository) FindByParticipationAndCollaborator(participationID, collaboratorID uuid.UUID) (*Rating, bool, error) {
	db, _ := database.AsGormConn(r.db)
	var rating Rating
	err := db.Preload("Participation").Preload("Collaborator").
		Where("participation_id = ? AND collaborator_id = ?", participationID, collaboratorID).
		First(&rating).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &rating, true, nil
}

func (r *gormRepository) FindByParticipationAndCollaboratorTx(tx database.Transaction, participationID, collaboratorID uuid.UUID) (*Rating, bool, error) {
	txdb, _ := database.AsGormTx(tx)
	var rating Rating
	err := txdb.Preload("Participation").Preload("Collaborator").
		Where("participation_id = ? AND collaborator_id = ?", participationID, collaboratorID).
		First(&rating).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &rating, true, nil
}

func (r *gormRepository) Update(rating *Rating) error {
	db, _ := database.AsGormConn(r.db)
	return db.Save(rating).Error
}

func (r *gormRepository) UpdateTx(tx database.Transaction, rating *Rating) error {
	txdb, _ := database.AsGormTx(tx)
	return txdb.Save(rating).Error
}

func (r *gormRepository) Upsert(rating *Rating) error {
	tx := r.db.BeginTransaction()
	if err := r.UpsertTx(tx, rating); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *gormRepository) UpsertTx(tx database.Transaction, rating *Rating) error {
	existing, found, err := r.FindByParticipationAndCollaboratorTx(tx, rating.ParticipationID, rating.CollaboratorID)
	if err != nil {
		return err
	}

	if found {
		// Update existing rating
		existing.Rate = rating.Rate
		existing.Comment = rating.Comment
		existing.RateType = rating.RateType
		return r.UpdateTx(tx, existing)
	} else {
		// Create new rating
		return r.CreateTx(tx, rating)
	}
}
