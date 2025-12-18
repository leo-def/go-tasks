package task

import (
    "go-tasks/internal/activity"
    "go-tasks/internal/assignment"
    "go-tasks/internal/lifecycle"
    "go-tasks/internal/pkg/database"
)

type Module struct {
    ControllerActivity *ControllerActivity
    ServiceActivity    *ServiceActivity
    ControllerActivityOwner *ControllerActivityOwner
    ServiceActivityOwner    *ServiceActivityOwner
}

func Initialize(db database.Connection, lifecycleService lifecycle.Service, assignmentService assignment.ServiceTask, activityService activity.Service) *Module {
    repository := NewRepository(db)
    service := NewService(repository, lifecycleService, assignmentService, activityService)
    controller := NewControllerActivity(service)
    ownerService := NewServiceActivityOwner(service, activityService)
    ownerController := NewControllerActivityOwner(ownerService)
    return &Module{ControllerActivity: controller, ServiceActivity: &service, ControllerActivityOwner: ownerController, ServiceActivityOwner: &ownerService}
}
