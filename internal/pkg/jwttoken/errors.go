package jwttoken

import (
	"errors"
)

var ErrJWTSecretNotSet = errors.New("JWT_SECRET is not set")
var ErrInvalidTokenClaims = errors.New("invalid token claims")
