package company

import "github.com/google/uuid"

// CompanyDTO represents the outward-facing API representation of a company
// Keep DTOs free of ORM-specific tags and use json tags only
// Add/rename fields here without impacting persistence layer
// Response
type CompanyDTO struct {
	ID    uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title string    `json:"title" example:"Acme Corp"`
}

type PaginatedCompanyDTO struct {
	Items []CompanyDTO `json:"items"`
	Count int64        `json:"count"`
}

// CompanyCreateDTO is the payload for creating a company.
// It includes an owner with account fields. On success, only the Company is returned.
// Create Param
type CompanyCreateDTO struct {
	Title string                `json:"title" binding:"required" example:"Acme Corp"`
	Owner CompanyOwnerCreateDTO `json:"owner" binding:"required"`
}

type CompanyOwnerCreateDTO struct {
	// Account contains credentials and profile details for the owner.
	// The password is hashed by the account service inside a DB transaction.
	Account CompanyOwnerAccountCreateDTO `json:"account" binding:"required"`
}

type CompanyOwnerAccountCreateDTO struct {
	Username string `json:"username" binding:"required" example:"owner_joe"`
	Password string `json:"password" binding:"required" example:"S3cureP@ss!"`
	Name     string `json:"name" binding:"required" example:"Joe Owner"`
	Phone    string `json:"phone" binding:"required" example:"+1-202-555-0123"`
	Email    string `json:"email" binding:"required" example:"owner@example.com"`
}

// CompanyUpdateDTO represents the payload for updating a company
// No ID here; ID is taken from path param at the controller level
// Update Param
type CompanyUpdateDTO struct {
	Title string `json:"title" binding:"required" example:"Acme Corp v2"`
}

// Mapping helpers
// Response
func ToCompanyDTO(m Company) CompanyDTO {
	return CompanyDTO{ID: m.ID, Title: m.Title}
}

func ToCompanyDTOs(ms []Company) []CompanyDTO {
	out := make([]CompanyDTO, 0, len(ms))
	for _, m := range ms {
		out = append(out, ToCompanyDTO(m))
	}
	return out
}

// FromCreateDTO converts a CompanyCreateDTO into a CompanyWithOwner struct.
// This is used when creating a new company with an owner.
// Create Param
func FromCreateDTO(d CompanyCreateDTO) CompanyWithOwner {
	return CompanyWithOwner{
		CompanyInfo: CompanyInfo{Title: d.Title},
		Owner: CompanyOwner{
			Account: CompanyOwnerAccount{
				Username: d.Owner.Account.Username,
				Password: d.Owner.Account.Password,
				Name:     d.Owner.Account.Name,
				Phone:    d.Owner.Account.Phone,
				Email:    d.Owner.Account.Email,
			},
		},
	}
}

// FromUpdateDTOWithID converts a CompanyUpdateDTO into a Company struct.
// This is used when updating an existing company.
// Update Param
func FromUpdateDTOWithID(id uuid.UUID, d CompanyUpdateDTO) Company {
	return Company{CompanyInfo: CompanyInfo{ID: id, Title: d.Title}}
}
