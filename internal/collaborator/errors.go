package collaborator

import "errors"

var ErrForbiddenNotSubordinate = errors.New("access forbidden: the targeted collaborator is not a subordinate")

var ErrForbiddenNotInCompany = errors.New("access forbidden: the targeted collaborator is not in the company")

var ErrCollaboratorNotFound = errors.New("collaborator not found")
