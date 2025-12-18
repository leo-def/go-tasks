package activity

import "errors"

var ErrOnwerNotFound = errors.New("owner not found")
var ErrForbiddenNotOwner = errors.New("access forbidden: the activity is not owned by the user")
var ErrForbiddenNotInCompany = errors.New("access forbidden: the activity is not in the company")
var ErrActivityNotFound = errors.New("activity not found")
