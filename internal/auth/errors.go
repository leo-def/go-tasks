package auth

import (
	"errors"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

var ErrNotCollaboratorInCompany = errors.New("not a collaborator in this company")

var ErrNotCollaborator = errors.New("not a collaborator")
