package participation

import (
	"github.com/google/uuid"
)

// Response

type ParticipationInfoDTO struct {
	ID             uuid.UUID `json:"id" example:"313e4567-e89b-12d3-a456-426614174000"`
	CollaboratorID uuid.UUID `json:"collaborator_id" example:"413e4567-e89b-12d3-a456-426614174000"`
	ActivityID     uuid.UUID `json:"activity_id" example:"513e4567-e89b-12d3-a456-426614174000"`
}

type ParticipationDTO struct {
	ParticipationInfoDTO
	Collaborator ParticipationCollaboratorDTO `json:"collaborator"`
	Activity     ParticipationActivityDTO     `json:"activity"`
}

type ParticipationCollaboratorDTO struct {
	ID        uuid.UUID               `json:"id" example:"613e4567-e89b-12d3-a456-426614174000"`
	Title     string                  `json:"title" example:"Team Lead"`
	AccountId uuid.UUID               `json:"account_id" example:"713e4567-e89b-12d3-a456-426614174000"`
	Account   ParticipationAccountDTO `json:"account"`
}

type ParticipationAccountDTO struct {
	ID       uuid.UUID `json:"id" example:"813e4567-e89b-12d3-a456-426614174000"`
	Username string    `json:"username" example:"collab_bob"`
	Name     string    `json:"name" example:"Bob Collaborator"`
	Phone    string    `json:"phone" example:"+1-202-555-0166"`
	Email    string    `json:"email" example:"bob@example.com"`
}

type ParticipationActivityDTO struct {
	ID    uuid.UUID `json:"id" example:"913e4567-e89b-12d3-a456-426614174000"`
	Title string    `json:"title" example:"Activity Alpha"`
}

type CollaboratorParticipationDTO struct {
	ParticipationInfo
	Activity ParticipationActivity `json:"activity"`
}

type ActivityParticipationDTO struct {
	ParticipationInfo
	Collaborator ParticipationCollaborator `json:"collaborator"`
}

type PaginatedCollaboratorParticipationDTO struct {
	Items []CollaboratorParticipationDTO `json:"items"`
	Count int64                          `json:"count"`
}

type PaginatedActivityParticipationDTO struct {
	Items []ActivityParticipationDTO `json:"items"`
	Count int64                      `json:"count"`
}

// Create Param
type ParticipationCreateDTO struct {
	CollaboratorID uuid.UUID `json:"collaborator_id"`
}

func FromCreateDTO(dto ParticipationCreateDTO) ParticipationInfo {
	return ParticipationInfo{
		CollaboratorID: dto.CollaboratorID,
	}
}

func ToParticipationInfoDTO(participation ParticipationInfo) ParticipationInfoDTO {
	return ParticipationInfoDTO{
		ID:             participation.ID,
		CollaboratorID: participation.CollaboratorID,
		ActivityID:     participation.ActivityID,
	}
}

func ToParticipationDTO(participation Participation) ParticipationDTO {
	return ParticipationDTO{
		ParticipationInfoDTO: ToParticipationInfoDTO(participation.ParticipationInfo),
		Collaborator:         ToParticipationCollaboratorDTO(participation.Collaborator),
		Activity:             ToParticipationActivityDTO(participation.Activity),
	}
}

func ToParticipationDTOs(list []Participation) []ParticipationDTO {
	dtos := make([]ParticipationDTO, 0, len(list))
	for _, item := range list {
		dtos = append(dtos, ToParticipationDTO(item))
	}
	return dtos
}

func ToParticipationCollaboratorDTO(collaborator ParticipationCollaborator) ParticipationCollaboratorDTO {
	return ParticipationCollaboratorDTO{
		ID:        collaborator.ID,
		Title:     collaborator.Title,
		AccountId: collaborator.AccountID,
		Account:   ToParticipationAccountDTO(collaborator.Account),
	}
}

func ToParticipationAccountDTO(account ParticipationAccount) ParticipationAccountDTO {
	return ParticipationAccountDTO(account)
}

func ToParticipationActivityDTO(activity ParticipationActivity) ParticipationActivityDTO {
	return ParticipationActivityDTO(activity)
}

func ToActivityParticipationDTO(activityParticipation ActivityParticipation) ActivityParticipationDTO {
	return ActivityParticipationDTO(activityParticipation)
}

func ToActivityParticipationDTOs(list []ActivityParticipation) []ActivityParticipationDTO {
	dtos := make([]ActivityParticipationDTO, 0, len(list))
	for _, item := range list {
		dtos = append(dtos, ToActivityParticipationDTO(item))
	}
	return dtos
}

func ToCollaboratorParticipationDTO(collaboratorParticipation CollaboratorParticipation) CollaboratorParticipationDTO {
	return CollaboratorParticipationDTO(collaboratorParticipation)
}

func ToCollaboratorParticipationDTOs(list []CollaboratorParticipation) []CollaboratorParticipationDTO {
	dtos := make([]CollaboratorParticipationDTO, 0, len(list))
	for _, item := range list {
		dtos = append(dtos, ToCollaboratorParticipationDTO(item))
	}
	return dtos
}
