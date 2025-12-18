package main

import (
	_ "go-tasks/docs"
	"go-tasks/internal/pkg/logger"
)

// @title GoTasks API
// @version 1.0
// @description GoTasks is a robust API for managing company tasks and collaborations.
// @description It provides comprehensive features for:
// @description - **Company and Collaborator Management**: Onboard companies and manage staff with hierarchical roles.
// @description - **Activity and Task Lifecycle Tracking**: Create and monitor activities and tasks through various stages.
// @description - **Role-Based Access Control (RBAC)**: Secure endpoints with granular permissions for Admins, Ops, Owners, and Managers.
// @description - **Performance Ratings and Assignments**: Assign tasks to collaborators and rate their performance.
// @description - **Secure Authentication**: JWT-based authentication for secure access.
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api/v1

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes http https
// @tag.name Auth | Public
// @tag.description Authentication endpoints accessible without a token
// @tag.name Auth | Protected
// @tag.description Authentication endpoints requiring a valid token
// @tag.name Company | Admin or Ops
// @tag.description Company management endpoints requiring Admin or Ops role
// @tag.name Collaborator | Company Owner
// @tag.description Collaborator endpoints requiring Company Owner role
// @tag.name Collaborator | Company
// @tag.description Collaborator endpoints for company members
// @tag.name Activity | Company Manager or Owner
// @tag.description Activity endpoints requiring Company Manager or Owner role
// @tag.name Participation | Company Manager or Owner
// @tag.description Participation endpoints requiring Company Manager or Owner role
// @tag.name Participation | Protected
// @tag.description Participation endpoints for authenticated users
// @tag.name Rating | Company Manager or Owner
// @tag.description Rating endpoints requiring Company Manager or Owner role
// @tag.name Task | Company Manager or Owner
// @tag.description Task endpoints requiring Company Manager or Owner role
// @tag.name Task | Company
// @tag.description Task endpoints for company members
// @tag.name Assignment | Company Manager or Owner
// @tag.description Assignment endpoints requiring Company Manager or Owner role
func main() {
	logger.Info("Initializing GoTasks API...\n", nil)
	application := InitializeApp()
	logger.Info("Running GoTasks API...\n", nil)
	application.Run()
}
