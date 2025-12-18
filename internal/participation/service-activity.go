package participation

import (
    "github.com/google/uuid"
)

type ServiceActivity interface {
    Create(participation *ParticipationInfo, activityId uuid.UUID) error
    Get(activityId uuid.UUID) ([]ActivityParticipation, int64, error)
    GetById(id uuid.UUID, activityId uuid.UUID) (*ActivityParticipation, bool, error)
}
type serviceActivity struct {
    repository Repository
}

func NewServiceActivity(Repository Repository) ServiceActivity {
    return &serviceActivity{Repository}
}

func (s *serviceActivity) Get(activityId uuid.UUID) ([]ActivityParticipation, int64, error) {
	return s.repository.GetByActivityId(activityId)
}

func (s *serviceActivity) GetById(id uuid.UUID, activityId uuid.UUID) (*ActivityParticipation, bool, error) {
	return s.repository.GetActivityParticipationById(id, activityId)
}

func (s *serviceActivity) Create(participation *ParticipationInfo, activityId uuid.UUID) error {
	participation.ActivityID = activityId
	return s.repository.Create(participation)
}
