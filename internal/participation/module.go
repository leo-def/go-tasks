package participation

import (
	"go-tasks/internal/pkg/database"
)

type Module struct {
	Controller         *Controller
	ControllerActivity *ControllerActivity
	ServiceActivity    *ServiceActivity
	Service            *Service
}

func Initialize(db database.Connection) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	serviceActivity := NewServiceActivity(repo)
	controllerActivity := NewControllerActivity(serviceActivity)
	controller := NewController(service)
	return &Module{
		Service:            &service,
		ServiceActivity:    &serviceActivity,
		Controller:         controller,
		ControllerActivity: controllerActivity}
}
