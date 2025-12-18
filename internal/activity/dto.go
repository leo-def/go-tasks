package activity

import (
	"go-tasks/internal/lifecycle"
	"time"

	"github.com/google/uuid"
)

// Response
type ActivityInfoDTO struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	LifecycleID uuid.UUID `json:"lifecycle_id"`
	CompanyID   uuid.UUID `json:"company_id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedByID uuid.UUID `json:"created_by_id"`
	Active      bool      `json:"active"`
}

type ActivityDTO struct {
	ActivityInfoDTO
	Company   ActivityCompanyDTO      `json:"company"`
	Lifecycle ActivityLifecycleDTO    `json:"lifecycle"`
	Owner     ActivityCollaboratorDTO `json:"owner"`
	CreatedBy ActivityCollaboratorDTO `json:"created_by"`
}

type CompanyActivityDTO struct {
	ActivityInfoDTO
	Lifecycle ActivityLifecycleDTO    `json:"lifecycle"`
	Owner     ActivityCollaboratorDTO `json:"owner"`
	CreatedBy ActivityCollaboratorDTO `json:"created_by"`
}

type OwnActivityDTO struct {
	ActivityInfoDTO
	Lifecycle ActivityLifecycleDTO    `json:"lifecycle"`
	CreatedBy ActivityCollaboratorDTO `json:"created_by"`
}

type PaginatedActivityDTO struct {
	Items []ActivityDTO `json:"items"`
	Count int64         `json:"count"`
}

type ActivityCompanyDTO struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type ActivityLifecycleDTO struct {
	ID       uuid.UUID                 `json:"id"`
	InitDate string                    `json:"init_date"`
	DueDate  string                    `json:"due_date"`
	Status   string                    `json:"status"`
	Updates  []ActivityStatusUpdateDTO `json:"updates"`
}

type ActivityStatusUpdateDTO struct {
	ID           uuid.UUID `json:"id"`
	LifecycleID  uuid.UUID `json:"lifecycle_id"`
	StatusBefore string    `json:"status_before"`
	StatusAfter  string    `json:"status_after"`
	UpdateDate   string    `json:"update_date"`
}

type ActivityCollaboratorDTO struct {
	ID        uuid.UUID          `json:"id"`
	Role      string             `json:"role"`
	AccountID uuid.UUID          `json:"account_id"`
	CompanyID uuid.UUID          `json:"company_id"`
	Account   ActivityAccountDTO `json:"account"`
}

type ActivityAccountDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
}

// Create Param
type ActivityCreateDTO struct {
	Title     string                     `json:"title"`
	OwnerID   uuid.UUID                  `json:"owner_id"`
	Lifecycle ActivityLifecycleCreateDTO `json:"lifecycle"`
}
type OwnActivityCreateDTO struct {
	Title     string                     `json:"title"`
	Lifecycle ActivityLifecycleCreateDTO `json:"lifecycle"`
}

type ActivityLifecycleCreateDTO struct {
	InitDate time.Time `json:"init_date"`
	DueDate  time.Time `json:"due_date"`
}

// Update Param
type ActivityUpdateDTO struct {
	Title     string                     `json:"title"`
	OwnerID   uuid.UUID                  `json:"owner_id"`
	Lifecycle ActivityLifecycleUpdateDTO `json:"lifecycle"`
}

type OwnActivityUpdateDTO struct {
	Title     string                     `json:"title"`
	Lifecycle ActivityLifecycleUpdateDTO `json:"lifecycle"`
}

type ActivityLifecycleUpdateDTO struct {
	InitDate time.Time `json:"init_date"`
	DueDate  time.Time `json:"due_date"`
}

// Update Status
type ActivityUpdateStatusDTO struct {
	Status lifecycle.LifecycleStatusDTO `json:"status"`
}

func CompanyActivityFromCreateDTO(dto ActivityCreateDTO) CompanyActivity {
	return CompanyActivity{
		ActivityInfo: ActivityInfo{
			Title:   dto.Title,
			OwnerID: dto.OwnerID,
		},
		Lifecycle: ActivityLifecycle{
			InitDate: dto.Lifecycle.InitDate,
			DueDate:  dto.Lifecycle.DueDate,
		},
	}
}

func FromCreateDTO(dto ActivityCreateDTO) Activity {
	return Activity{
		ActivityInfo: ActivityInfo{
			Title:   dto.Title,
			OwnerID: dto.OwnerID,
		},
		Lifecycle: ActivityLifecycle{
			InitDate: dto.Lifecycle.InitDate,
			DueDate:  dto.Lifecycle.DueDate,
		},
	}
}

func OwnActivityFromCreateDTO(dto OwnActivityCreateDTO) OwnActivity {
	return OwnActivity{
		ActivityInfo: ActivityInfo{
			Title: dto.Title,
		},
		Lifecycle: ActivityLifecycle{
			InitDate: dto.Lifecycle.InitDate,
			DueDate:  dto.Lifecycle.DueDate,
		},
	}
}

func CompanyActivityFromUpdateDTOWithID(id uuid.UUID, dto ActivityUpdateDTO) CompanyActivity {
	return CompanyActivity{
		ActivityInfo: ActivityInfo{
			ID:      id,
			Title:   dto.Title,
			OwnerID: dto.OwnerID,
		},
		Lifecycle: ActivityLifecycle{
			InitDate: dto.Lifecycle.InitDate,
			DueDate:  dto.Lifecycle.DueDate,
		},
	}
}

func FromUpdateDTOWithID(id uuid.UUID, dto ActivityUpdateDTO) Activity {
	return Activity{
		ActivityInfo: ActivityInfo{
			ID:      id,
			Title:   dto.Title,
			OwnerID: dto.OwnerID,
		},
		Lifecycle: ActivityLifecycle{
			InitDate: dto.Lifecycle.InitDate,
			DueDate:  dto.Lifecycle.DueDate,
		},
	}
}

func OwnActivityFromUpdateDTOWithID(id uuid.UUID, dto OwnActivityUpdateDTO) OwnActivity {
	return OwnActivity{
		ActivityInfo: ActivityInfo{
			ID:    id,
			Title: dto.Title,
		},
		Lifecycle: ActivityLifecycle{
			InitDate: dto.Lifecycle.InitDate,
			DueDate:  dto.Lifecycle.DueDate,
		},
	}
}

func ToActivityDTO(activity *Activity) ActivityDTO {
	return ActivityDTO{
		ActivityInfoDTO: ToActivityInfoDTO(&activity.ActivityInfo),
		Company:         ToActivityCompanyDTO(&activity.Company),
		Lifecycle:       ToActivityLifecycleDTO(&activity.Lifecycle),
		Owner:           ToActivityCollaboratorDTO(&activity.Owner),
		CreatedBy:       ToActivityCollaboratorDTO(&activity.CreatedBy),
	}
}

func ToCompanyActivityDTO(activity *CompanyActivity) CompanyActivityDTO {
	return CompanyActivityDTO{
		ActivityInfoDTO: ToActivityInfoDTO(&activity.ActivityInfo),
		Lifecycle:       ToActivityLifecycleDTO(&activity.Lifecycle),
		Owner:           ToActivityCollaboratorDTO(&activity.Owner),
		CreatedBy:       ToActivityCollaboratorDTO(&activity.CreatedBy),
	}
}

func ToOwnActivityDTO(activity *OwnActivity) OwnActivityDTO {
	return OwnActivityDTO{
		ActivityInfoDTO: ToActivityInfoDTO(&activity.ActivityInfo),
		Lifecycle:       ToActivityLifecycleDTO(&activity.Lifecycle),
		CreatedBy:       ToActivityCollaboratorDTO(&activity.CreatedBy),
	}
}

func ToActivityInfoDTO(activity *ActivityInfo) ActivityInfoDTO {
	return ActivityInfoDTO{
		ID:          activity.ID,
		Title:       activity.Title,
		LifecycleID: activity.LifecycleID,
		CompanyID:   activity.CompanyID,
		OwnerID:     activity.OwnerID,
		CreatedByID: activity.CreatedByID,
		Active:      activity.Active,
	}
}

func ToActivityCompanyDTO(company *ActivityCompany) ActivityCompanyDTO {
	return ActivityCompanyDTO{
		ID:    company.ID,
		Title: company.Title,
	}
}

func ToActivityLifecycleDTO(lifecycle *ActivityLifecycle) ActivityLifecycleDTO {
	return ActivityLifecycleDTO{
		ID:       lifecycle.ID,
		InitDate: lifecycle.InitDate.Format(time.RFC3339),
		DueDate:  lifecycle.DueDate.Format(time.RFC3339),
		Status:   string(lifecycle.Status),
		Updates:  ToActivityStatusUpdateDTOs(lifecycle.Updates),
	}
}

func ToActivityStatusUpdateDTO(update ActivityStatusUpdate) ActivityStatusUpdateDTO {
	return ActivityStatusUpdateDTO{
		ID:           update.ID,
		LifecycleID:  update.LifecycleID,
		StatusBefore: string(update.StatusBefore),
		StatusAfter:  string(update.StatusAfter),
		UpdateDate:   update.UpdateDate.Format(time.RFC3339),
	}
}

func ToActivityStatusUpdateDTOs(updates []ActivityStatusUpdate) []ActivityStatusUpdateDTO {
	dtos := make([]ActivityStatusUpdateDTO, 0, len(updates))
	for _, u := range updates {
		dtos = append(dtos, ToActivityStatusUpdateDTO(u))
	}
	return dtos
}

func ToActivityAccountDTO(account *ActivityAccount) ActivityAccountDTO {
	return ActivityAccountDTO{
		ID:       account.ID,
		Username: account.Username,
		Name:     account.Name,
		Phone:    account.Phone,
		Email:    account.Email,
	}
}

func ToActivityCollaboratorDTO(collaborator *ActivityCollaborator) ActivityCollaboratorDTO {
	return ActivityCollaboratorDTO{
		ID:        collaborator.ID,
		Role:      collaborator.Role,
		AccountID: collaborator.AccountID,
		CompanyID: collaborator.CompanyID,
		Account:   ToActivityAccountDTO(&collaborator.Account),
	}
}

func ToActivityDTOs(activities []Activity) []ActivityDTO {
	dtos := make([]ActivityDTO, 0, len(activities))
	for _, activity := range activities {
		dtos = append(dtos, ToActivityDTO(&activity))
	}
	return dtos
}

func ToCompanyActivityDTOs(activities []CompanyActivity) []CompanyActivityDTO {
	dtos := make([]CompanyActivityDTO, 0, len(activities))
	for _, activity := range activities {
		dtos = append(dtos, ToCompanyActivityDTO(&activity))
	}
	return dtos
}

func ToOwnActivityDTOs(activities []OwnActivity) []OwnActivityDTO {
	dtos := make([]OwnActivityDTO, 0, len(activities))
	for _, activity := range activities {
		dtos = append(dtos, ToOwnActivityDTO(&activity))
	}
	return dtos
}
