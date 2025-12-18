package app

import (
	"go-tasks/internal/account"
	"go-tasks/internal/activity"
	"go-tasks/internal/assignment"
	"go-tasks/internal/auth"
	"go-tasks/internal/collaborator"
	"go-tasks/internal/company"
	"go-tasks/internal/health"
	"go-tasks/internal/lifecycle"
	"go-tasks/internal/participation"
	"go-tasks/internal/pkg/database"
	"go-tasks/internal/rating"
	"go-tasks/internal/task"
)

type Module struct {
	HealthModule        *health.Module
	AccountModule       *account.Module
	CollaboratorModule  *collaborator.Module
	AuthModule          *auth.Module
	CompanyModule       *company.Module
	LifecycleModule     *lifecycle.Module
	ParticipationModule *participation.Module
	RatingModule        *rating.Module
	AssignmentModule    *assignment.Module
	ActivityModule      *activity.Module
	TaskModule          *task.Module
}

func Initialize(db database.Connection) *Module {
	healthModule := health.Initialize()
	accountModule := account.Initialize(db)
	collaboratorModule := collaborator.Initialize(db, *accountModule.Service)
	authModule := auth.Initialize(*accountModule.Service, *collaboratorModule.Service)
	companyModule := company.Initialize(db, *accountModule.Service)
	lifecycleModule := lifecycle.Initialize(db)
	participationModule := participation.Initialize(db)
	ratingModule := rating.Initialize(db, *participationModule.Service)
	assignmentModule := assignment.Initialize(db, *participationModule.Service)
	activityModule := activity.Initialize(db, lifecycleModule.Service)
	taskModule := task.Initialize(db, lifecycleModule.Service, *assignmentModule.Service, *activityModule.Service)
	return &Module{
		HealthModule:        healthModule,
		AccountModule:       accountModule,
		CollaboratorModule:  collaboratorModule,
		AuthModule:          authModule,
		CompanyModule:       companyModule,
		LifecycleModule:     lifecycleModule,
		ParticipationModule: participationModule,
		RatingModule:        ratingModule,
		AssignmentModule:    assignmentModule,
		ActivityModule:      activityModule,
		TaskModule:          taskModule,
	}
}
