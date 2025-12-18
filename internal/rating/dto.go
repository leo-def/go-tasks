package rating

import (
	"go-tasks/internal/collaborator"

	"github.com/google/uuid"
)

type RatingInfoDTO struct {
	ID              uuid.UUID `json:"id" example:"a13e4567-e89b-12d3-a456-426614174000"`
	Rate            float64   `json:"rate" example:"4.5"`
	ParticipationID uuid.UUID `json:"participation_id" example:"b13e4567-e89b-12d3-a456-426614174000"`
	CollaboratorID  uuid.UUID `json:"collaborator_id" example:"c13e4567-e89b-12d3-a456-426614174000"`
}

type RatingDTO struct {
	RatingInfoDTO
	Participation RatingParticipationDTO `json:"participation"`
	Collaborator  RatingCollaboratorDTO  `json:"collaborator"`
}

type RatingParticipationDTO struct {
	ID             uuid.UUID `json:"id" example:"d13e4567-e89b-12d3-a456-426614174000"`
	CollaboratorID uuid.UUID `json:"collaborator_id" example:"e13e4567-e89b-12d3-a456-426614174000"`
	ActivityID     uuid.UUID `json:"activity_id" example:"f13e4567-e89b-12d3-a456-426614174000"`
}

type RatingCollaboratorDTO struct {
	ID        uuid.UUID                     `json:"id" example:"013e4567-e89b-12d3-a456-426614174000"`
	CompanyID uuid.UUID                     `json:"company_id" example:"113e4567-e89b-12d3-a456-426614174000"`
	AccountID uuid.UUID                     `json:"account_id" example:"213e4567-e89b-12d3-a456-426614174000"`
	Role      collaborator.CollaboratorRole `json:"role" example:"OWNER"`
	Active    bool                          `json:"active" example:"true"`
}

// Create
type RatingCreateDTO struct {
	Rate     float64 `json:"rate" example:"5"`
	Comment  string  `json:"comment"`
	RateType string  `json:"rate_type"`
}

type ActivityRatingCreateDTO struct {
	RatingCreateDTO
	ActivityID uuid.UUID `json:"activity_id" example:"f13e4567-e89b-12d3-a456-426614174000"`
}

func FromCreateDTO(dto RatingCreateDTO) Rating {
	return Rating{
		RatingInfo: RatingInfo{
			Rate:     dto.Rate,
			Comment:  dto.Comment,
			RateType: dto.RateType,
		},
	}
}

func FromActivityRatingCreateDTO(dto ActivityRatingCreateDTO) Rating {
	return FromCreateDTO(dto.RatingCreateDTO)
}

func ToRatingDTO(rating *Rating) RatingDTO {
	return RatingDTO{
		RatingInfoDTO: RatingInfoDTO{
			ID:              rating.ID,
			Rate:            rating.Rate,
			ParticipationID: rating.ParticipationID,
			CollaboratorID:  rating.CollaboratorID,
		},
		Participation: ToRatingParticipationDTO(rating.Participation),
		Collaborator:  ToRatingCollaboratorDTO(rating.Collaborator),
	}
}

func ToRatingParticipationDTO(p RatingParticipation) RatingParticipationDTO {
	return RatingParticipationDTO{
		ID:             p.ID,
		CollaboratorID: p.CollaboratorID,
		ActivityID:     p.ActivityID,
	}
}

func ToRatingCollaboratorDTO(c RatingCollaborator) RatingCollaboratorDTO {
	return RatingCollaboratorDTO{
		ID:        c.ID,
		CompanyID: c.CompanyID,
		AccountID: c.AccountID,
		Role:      c.Role,
		Active:    c.Active,
	}
}
