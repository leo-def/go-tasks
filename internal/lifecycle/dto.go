package lifecycle

import (
	"time"

	"github.com/google/uuid"
)

type LifecycleStatusDTO string

const (
	LifecycleStatusCreatedDTO    LifecycleStatusDTO = "CREATED"
	LifecycleStatusInProgressDTO LifecycleStatusDTO = "IN_PROGRESS"
	LifecycleStatusBlockedDTO    LifecycleStatusDTO = "BLOCKED"
	LifecycleStatusCompletedDTO  LifecycleStatusDTO = "COMPLETED"
	LifecycleStatusCancelledDTO  LifecycleStatusDTO = "CANCELLED"
)

type LifecycleInfoDTO struct {
    ID       uuid.UUID          `json:"id"`
    InitDate time.Time          `json:"init_date"`
    DueDate  time.Time          `json:"due_date"`
    ParentID uuid.UUID          `json:"parent_id"`
    Status   LifecycleStatusDTO `json:"status"`
}

type LifecycleDTO struct {
    LifecycleInfoDTO
    Updates []LifecycleStatusUpdateDTO `json:"updates"`
    Parent  LifecycleInfoDTO           `json:"parent"`
    Children []LifecycleInfoDTO        `json:"children"`
}

type LifecycleStatusUpdateDTO struct {
    ID           uuid.UUID          `json:"id"`
    LifecycleID  uuid.UUID          `json:"lifecycle_id"`
    StatusBefore LifecycleStatusDTO `json:"status_before"`
    StatusAfter  LifecycleStatusDTO `json:"status_after"`
    UpdateDate   time.Time          `json:"update_date"`
}

type UpdateLifecycleStatusDTO struct {
    ID     uuid.UUID          `json:"id"`
    Status LifecycleStatusDTO `json:"status"`
}

// Mapping helpers
func ToLifecycleInfoDTO(m LifecycleInfo) LifecycleInfoDTO {
    return LifecycleInfoDTO{
        ID:       m.ID,
        InitDate: m.InitDate,
        DueDate:  m.DueDate,
        ParentID: m.ParentID,
        Status:   LifecycleStatusDTO(m.Status),
    }
}

func ToLifecycleInfoDTOs(ms []LifecycleInfo) []LifecycleInfoDTO {
    out := make([]LifecycleInfoDTO, 0, len(ms))
    for _, m := range ms {
        out = append(out, ToLifecycleInfoDTO(m))
    }
    return out
}

func ToLifecycleStatusUpdateDTO(u LifecycleStatusUpdate) LifecycleStatusUpdateDTO {
    return LifecycleStatusUpdateDTO{
        ID:           u.ID,
        LifecycleID:  u.LifecycleID,
        StatusBefore: LifecycleStatusDTO(u.StatusBefore),
        StatusAfter:  LifecycleStatusDTO(u.StatusAfter),
        UpdateDate:   u.UpdateDate,
    }
}

func ToLifecycleStatusUpdateDTOs(us []LifecycleStatusUpdate) []LifecycleStatusUpdateDTO {
    out := make([]LifecycleStatusUpdateDTO, 0, len(us))
    for _, u := range us {
        out = append(out, ToLifecycleStatusUpdateDTO(u))
    }
    return out
}

func ToLifecycleDTO(l Lifecycle) LifecycleDTO {
    return LifecycleDTO{
        LifecycleInfoDTO: ToLifecycleInfoDTO(l.LifecycleInfo),
        Updates:           ToLifecycleStatusUpdateDTOs(l.Updates),
        Parent:            ToLifecycleInfoDTO(l.Parent),
        Children:          ToLifecycleInfoDTOs(l.Children),
    }
}
