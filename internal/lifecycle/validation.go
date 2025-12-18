package lifecycle

import (
	"slices"
)

var LifecycleNextStatus = map[LifecycleStatus][]LifecycleStatus{
	LifecycleStatusCreated:    {LifecycleStatusBlocked, LifecycleStatusCancelled, LifecycleStatusInProgress},
	LifecycleStatusInProgress: {LifecycleStatusCompleted, LifecycleStatusBlocked, LifecycleStatusCancelled},
	LifecycleStatusBlocked:    {LifecycleStatusInProgress, LifecycleStatusCancelled},
	LifecycleStatusCompleted:  {},
	LifecycleStatusCancelled:  {},
}

var ActiveLifecycleStatus = []LifecycleStatus{
	LifecycleStatusCreated,
	LifecycleStatusInProgress,
	LifecycleStatusBlocked,
}

var InactiveLifecycleStatus = []LifecycleStatus{
	LifecycleStatusCompleted,
	LifecycleStatusCancelled,
}

type LifecycleValidation interface {
	ValidateCreation(lifecycle *LifecycleInfo, parent *LifecycleInfo) error
	ValidateUpdate(old *LifecycleInfo, lifecycle *Lifecycle) error
	ValidateDeletion(lifecycle *Lifecycle) error
	ValidateStatusUpdate(old *Lifecycle, status LifecycleStatus) error
}

type lifecycleValidationService struct {
}

func NewValidationService() LifecycleValidation {
	return &lifecycleValidationService{}
}

func (v *lifecycleValidationService) ValidateCreation(lifecycle *LifecycleInfo, parent *LifecycleInfo) error {
    err := v.validateLifecycleDate(lifecycle)
    if err != nil {
        return err
    }
    if lifecycle.Status != LifecycleStatusCreated {
        return ErrInvalidCreationStatus
    }
    if parent.HasID() {
        err = v.validateDateBounds(lifecycle, parent)
        if err != nil {
            return err
        }
        if !v.isStatusActive(parent.Status) {
            return ErrInvalidParentStatusForCreation
        }
    }
    return nil
}

func (v *lifecycleValidationService) ValidateUpdate(old *LifecycleInfo, lifecycle *Lifecycle) error {
    if old.Status != lifecycle.Status {
        return ErrInvalidStatusChangeOnUpdate
    }
    err := v.validateLifecycleDate(&lifecycle.LifecycleInfo)
    if err != nil {
        return err
    }
    if lifecycle.HasParentID() {
        err = v.validateDateBounds(&lifecycle.LifecycleInfo, &lifecycle.Parent)
        if err != nil {
            return err
        }
        if !v.isStatusActive(lifecycle.Parent.Status) {
            return ErrInvalidParentStatusForUpdate
        }
    }
    err = v.validateChildren(lifecycle, func(child *LifecycleInfo) error {
        return v.validateDateBounds(child, &lifecycle.LifecycleInfo)
    })
    if err != nil {
        return err
    }
    return nil
}

func (v *lifecycleValidationService) ValidateDeletion(lifecycle *Lifecycle) error {
    err := v.validateChildren(lifecycle, func(child *LifecycleInfo) error {
        if v.isStatusActive(child.Status) {
            return ErrActiveChildLifecycleDeletionNotAllowed
        }
        return nil
    })
    if err != nil {
        return err
    }
    if v.isStatusActive(lifecycle.Status) {
        return ErrActiveLifecycleDeletionNotAllowed
    }
    return nil
}

func (v *lifecycleValidationService) ValidateStatusUpdate(old *Lifecycle, status LifecycleStatus) error {
    err := v.validateNextStatus(&old.LifecycleInfo, status)
    if err != nil {
        return err
    }

    activeStatus := v.isStatusActive(status)
    if activeStatus {
        if old.HasParentID() {
            if !v.isStatusActive(old.Parent.Status) {
                return ErrInvalidParentStatusForStatusUpdate
            }
        }
    } else {
        err = v.validateChildren(old, func(child *LifecycleInfo) error {
            if v.isStatusActive(child.Status) {
                return ErrInvalidChildStatusForStatusUpdate
            }
            return nil
        })
        if err != nil {
            return err
        }
        return nil
    }
    return nil
}

// private

func (*lifecycleValidationService) validateDateBounds(lifecycle *LifecycleInfo, parent *LifecycleInfo) error {
	if !parent.HasID() {
		return nil
	}
	if lifecycle.InitDate.Before(parent.InitDate) || lifecycle.InitDate.After(parent.DueDate) {
		return ErrInvalidLifecycleDate
	}
	if lifecycle.DueDate.Before(parent.InitDate) || lifecycle.DueDate.After(parent.DueDate) {
		return ErrInvalidLifecycleDate
	}
	return nil
}

func (*lifecycleValidationService) validateChildren(lifecycle *Lifecycle, validation func(child *LifecycleInfo) error) error {
	for _, child := range lifecycle.Children {
		err := validation(&child)
		if err != nil {
			return err
		}
	}
	return nil
}

func (*lifecycleValidationService) validateLifecycleDate(lifecycleInfo *LifecycleInfo) error {
	if lifecycleInfo.InitDate.After(lifecycleInfo.DueDate) {
		return ErrInvalidLifecycleDate
	}
	return nil
}

func (*lifecycleValidationService) validateNextStatus(lifecycleInfo *LifecycleInfo, status LifecycleStatus) error {
    nextStatusOptions := LifecycleNextStatus[lifecycleInfo.Status]
    if !slices.Contains(nextStatusOptions, status) {
        return ErrInvalidNextStatusTransition
    }
    return nil
}

func (*lifecycleValidationService) isStatusActive(status LifecycleStatus) bool {
	return slices.Contains(ActiveLifecycleStatus, status)
}
