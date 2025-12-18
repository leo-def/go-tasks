package participation

import (
	"go-tasks/internal/pkg/database"

	"github.com/google/uuid"
)

type Service interface {
    Create(participation *ParticipationInfo) error
    GetByCollaboratorId(collaboratorId uuid.UUID) ([]CollaboratorParticipation, int64, error)
    GetById(id uuid.UUID) (*Participation, bool, error)
    GetOrCreateByActivityAndCollaboratorTx(tx database.Transaction, activityId, collaboratorId uuid.UUID) (ParticipationInfo, bool, error)
}
type service struct {
	repository Repository
}

func NewService(Repository Repository) Service {
	return &service{Repository}
}

func (s *service) GetByCollaboratorId(collaboratorId uuid.UUID) ([]CollaboratorParticipation, int64, error) {
	return s.repository.GetByCollaboratorId(collaboratorId)
}

func (s *service) GetOrCreateByActivityAndCollaboratorTx(tx database.Transaction, activityId, collaboratorId uuid.UUID) (ParticipationInfo, bool, error) {
	participation, found, err := s.repository.GetInfoByActivityIdAndCollaboratorId(activityId, collaboratorId)
	if err != nil {
		return ParticipationInfo{}, false, err
	}
	if !found {
		participation = ParticipationInfo{
			ActivityID:     activityId,
			CollaboratorID: collaboratorId,
		}
		err = s.repository.CreateTx(tx, &participation)
		if err != nil {
			return ParticipationInfo{}, false, err
		}
	}
	return participation, true, nil
}

func (s *service) GetById(id uuid.UUID) (*Participation, bool, error) {
	return s.repository.GetById(id)
}

func (s *service) Create(participation *ParticipationInfo) error {
	return s.repository.Create(participation)
}
