package rating

import (
	"go-tasks/internal/participation"

	"github.com/google/uuid"
)

type Service interface {
	CreateForCollaborator(rating *Rating, assigneeID, activityID, collaboratorID uuid.UUID) error
	CreateForParticipation(rating *Rating, assigneeID, participationID uuid.UUID) error
}

type service struct {
	repository           Repository
	participationService participation.Service
}

func NewService(repository Repository, participationService participation.Service) Service {
	return &service{repository: repository, participationService: participationService}
}

func (s *service) CreateForParticipation(rating *Rating, assigneeID, participationID uuid.UUID) error {
	rating.ParticipationID = participationID
	rating.CollaboratorID = assigneeID
	return s.repository.UpsertTx(s.repository.GetConnection().BeginTransaction(), rating)
}

func (s *service) CreateForCollaborator(rating *Rating, assigneeID, activityID, collaboratorID uuid.UUID) error {
	tx := s.repository.GetConnection().BeginTransaction()
	participation, found, err := s.participationService.GetOrCreateByActivityAndCollaboratorTx(tx, activityID, collaboratorID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !found {
		println("participation not found, created new activityID: ", activityID.String(), " collaboratorID: ", collaboratorID.String())
	}
	rating.ParticipationID = participation.ID
	rating.CollaboratorID = assigneeID
	if err := s.repository.UpsertTx(tx, rating); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
