package lifecycle

import "errors"

var ErrLifecycleNotFound = errors.New("lifecycle not found")

// Generic status transition error (kept for backward compatibility where needed)
var ErrInvalidLifecycleStatusTransition = errors.New("invalid lifecycle status transition")

// Creation validations
var ErrInvalidCreationStatus = errors.New("invalid creation status: status must be CREATED")
var ErrInvalidParentStatusForCreation = errors.New("invalid parent status for creation: parent must be active")

// Update validations
var ErrInvalidStatusChangeOnUpdate = errors.New("invalid update: status cannot change during update")
var ErrInvalidParentStatusForUpdate = errors.New("invalid parent status for update: parent must be active")

// Status update validations
var ErrInvalidNextStatusTransition = errors.New("invalid status transition: next status not allowed from current")
var ErrInvalidParentStatusForStatusUpdate = errors.New("invalid parent status for status update: parent must be active to set active status")
var ErrInvalidChildStatusForStatusUpdate = errors.New("invalid child status for status update: children must be inactive to set inactive status")

var ErrInvalidLifecycleDate = errors.New("invalid lifecycle date")

// Deletion validations
var ErrInvalidDeletionStatus = errors.New("invalid deletion status")
var ErrActiveLifecycleDeletionNotAllowed = errors.New("invalid deletion: lifecycle is active")
var ErrActiveChildLifecycleDeletionNotAllowed = errors.New("invalid deletion: child lifecycle is active")
