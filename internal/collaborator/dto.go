package collaborator

import (
	"github.com/google/uuid"
)

type ManagedCollaboratorRoleDTO string

const (
	ManagedCollaboratorRoleDTOManager ManagedCollaboratorRoleDTO = "MANAGER"
	ManagedCollaboratorRoleDTOOps     ManagedCollaboratorRoleDTO = "OPS"
)

type CollaboratorRoleDTO string

const (
	CollaboratorRoleDTOOwner   CollaboratorRoleDTO = "OWNER"
	CollaboratorRoleDTOManager CollaboratorRoleDTO = "MANAGER"
	CollaboratorRoleDTOOps     CollaboratorRoleDTO = "OPS"
)

type CreateCollaboratorRole string

// Response

type CollaboratorInfoDTO struct {
	ID        uuid.UUID           `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CompanyID uuid.UUID           `json:"company_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	AccountID uuid.UUID           `json:"account_id" example:"123e4567-e89b-12d3-a456-426614174002"`
	Role      CollaboratorRoleDTO `json:"role" example:"OWNER"`
	Active    bool                `json:"active" example:"true"`
}

type CollaboratorDTO struct {
	CollaboratorInfoDTO
	Company CollaboratorCompanyDTO `json:"company"`
	Account CollaboratorAccountDTO `json:"account"`
}

type PaginatedAccountCollaborator struct {
	Items []AccountCollaboratorDTO `json:"items"`
	Count int64                    `json:"count"`
}

type PaginatedCollaboratorDTO struct {
	Items []CollaboratorDTO `json:"items"`
	Count int64             `json:"count"`
}

type CollaboratorAccountDTO struct {
	ID       uuid.UUID `json:"id" example:"223e4567-e89b-12d3-a456-426614174000"`
	Username string    `json:"username" example:"jdoe"`
	Name     string    `json:"name" example:"John Doe"`
	Phone    string    `json:"phone" example:"+1-202-555-0123"`
	Email    string    `json:"email" example:"john.doe@example.com"`
}

type CollaboratorCompanyDTO struct {
	ID    uuid.UUID `json:"id" example:"323e4567-e89b-12d3-a456-426614174000"`
	Title string    `json:"title" example:"Acme Corp"`
}

type CompanyCollaboratorDTO struct {
	CollaboratorInfoDTO
	Account CollaboratorAccountDTO `json:"account"`
}

type AccountCollaboratorDTO struct {
	CollaboratorInfoDTO
	Company CollaboratorCompanyDTO `json:"company"`
}

type PaginatedCompanyCollaboratorDTO struct {
	Items []CompanyCollaboratorDTO `json:"items"`
	Count int64                    `json:"count"`
}

// Create Param
type CollaboratorCreateDTO struct {
	AccountID uuid.UUID                  `json:"account_id" example:"223e4567-e89b-12d3-a456-426614174000"`
	Role      ManagedCollaboratorRoleDTO `json:"role" binding:"required" enums:"MANAGER,OPS" example:"MANAGER"`
}

type CollaboratorWithAccountCreateDTO struct {
	Role    ManagedCollaboratorRoleDTO   `json:"role" binding:"required" enums:"MANAGER,OPS" example:"MANAGER"`
	Account CollaboratorAccountCreateDTO `json:"account" binding:"required"`
}

type CollaboratorAccountCreateDTO struct {
	Username string `json:"username" example:"owner_joe"`
	Password string `json:"password" example:"S3cureP@ss!"`
	Name     string `json:"name" example:"Joe Owner"`
	Phone    string `json:"phone" example:"+1-202-555-0199"`
	Email    string `json:"email" example:"owner@example.com"`
}

// Update Param
type CollaboratorUpdateDTO struct {
	Role    ManagedCollaboratorRoleDTO   `json:"role" binding:"required" enums:"MANAGER,OPS" example:"MANAGER"`
	Account CollaboratorAccountUpdateDTO `json:"account" binding:"required"`
}

type CollaboratorAccountUpdateDTO struct {
	Name string `json:"name" example:"Jane Manager"`
}

// When create a new one with a already existing account
func FromCreateDTOWithCompanyID(dto CollaboratorCreateDTO, companyID uuid.UUID) CollaboratorInfo {
	return CollaboratorInfo{
		Role:      CollaboratorRole(dto.Role),
		AccountID: dto.AccountID,
		CompanyID: companyID,
	}
}

// When create a new one with a new account
func CompanyCollaboratorFromCreateNewAccountDTO(dto CollaboratorWithAccountCreateDTO, companyID uuid.UUID) CompanyCollaborator {
	return CompanyCollaborator{
		CollaboratorInfo: CollaboratorInfo{
			Role:      CollaboratorRole(dto.Role),
			CompanyID: companyID,
		},
		Account: CollaboratorAccount{
			Username: dto.Account.Username,
			Password: dto.Account.Password,
			Name:     dto.Account.Name,
			Phone:    dto.Account.Phone,
			Email:    dto.Account.Email,
		},
	}
}

// When update a existing one
func CompanyCollaboratorFromUpdateDTO(dto CollaboratorUpdateDTO, companyID uuid.UUID) CompanyCollaborator {
	return CompanyCollaborator{
		CollaboratorInfo: CollaboratorInfo{
			Role: CollaboratorRole(dto.Role),
		},
		Account: CollaboratorAccount{
			Name: dto.Account.Name,
		},
	}
}

// FromUpdateDTOWithID converts update DTO plus an ID into a domain Collaborator.
// Standardized interface: id is uuid.UUID.
func CompanyCollaboratorFromUpdateDTOWithID(dto CollaboratorUpdateDTO, id, companyID uuid.UUID) CompanyCollaborator {
	c := CompanyCollaboratorFromUpdateDTO(dto, companyID)
	c.ID = id
	return c
}

// To response DTO
func ToCollaboratorInfoDTO(collaboratorInfo CollaboratorInfo) CollaboratorInfoDTO {
	return CollaboratorInfoDTO{
		ID:        collaboratorInfo.ID,
		CompanyID: collaboratorInfo.CompanyID,
		AccountID: collaboratorInfo.AccountID,
		Role:      CollaboratorRoleDTO(collaboratorInfo.Role),
		Active:    collaboratorInfo.Active,
	}
}

func ToCollaboratorDTO(collaborator Collaborator) CollaboratorDTO {
	return CollaboratorDTO{
		CollaboratorInfoDTO: ToCollaboratorInfoDTO(collaborator.CollaboratorInfo),
		Account:             ToCollaboratorAccountDTO(collaborator.Account),
		Company:             ToCollaboratorCompanyDTO(collaborator.Company),
	}
}

// ToCollaboratorDTOs converts a slice of domain Collaborator into DTOs.
func ToCollaboratorDTOs(ms []Collaborator) []CollaboratorDTO {
	out := make([]CollaboratorDTO, 0, len(ms))
	for _, m := range ms {
		out = append(out, ToCollaboratorDTO(m))
	}
	return out
}

func ToCompanyCollaboratorDTO(collaborator CompanyCollaborator) CompanyCollaboratorDTO {
	return CompanyCollaboratorDTO{
		CollaboratorInfoDTO: ToCollaboratorInfoDTO(collaborator.CollaboratorInfo),
		Account:             ToCollaboratorAccountDTO(collaborator.Account),
	}
}

// ToCompanyCollaboratorDTOs converts a slice of domain CompanyCollaborator into DTOs.
func ToCompanyCollaboratorDTOs(ms []CompanyCollaborator) []CompanyCollaboratorDTO {
	out := make([]CompanyCollaboratorDTO, 0, len(ms))
	for _, m := range ms {
		out = append(out, ToCompanyCollaboratorDTO(m))
	}
	return out
}

func ToAccountCollaboratorDTO(collaborator AccountCollaborator) AccountCollaboratorDTO {
	return AccountCollaboratorDTO{
		CollaboratorInfoDTO: ToCollaboratorInfoDTO(collaborator.CollaboratorInfo),
		Company:             ToCollaboratorCompanyDTO(collaborator.Company),
	}
}

// ToAccountCollaboratorDTOs converts a slice of domain AccountCollaborator into DTOs.
func ToAccountCollaboratorDTOs(ms []AccountCollaborator) []AccountCollaboratorDTO {
	out := make([]AccountCollaboratorDTO, 0, len(ms))
	for _, m := range ms {
		out = append(out, ToAccountCollaboratorDTO(m))
	}
	return out
}

func ToCollaboratorAccountDTO(account CollaboratorAccount) CollaboratorAccountDTO {
	return CollaboratorAccountDTO{
		ID:       account.ID,
		Username: account.Username,
		Name:     account.Name,
		Phone:    account.Phone,
		Email:    account.Email,
	}
}

func ToCollaboratorCompanyDTO(company CollaboratorCompany) CollaboratorCompanyDTO {
	return CollaboratorCompanyDTO(company)
}
