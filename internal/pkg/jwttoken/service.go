package jwttoken

import (
	"go-tasks/internal/pkg/env"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(authData *AuthData) (string, error)
	ParseToken(token string) (*AuthData, error)
}

type tokenService struct {
}

func NewService() TokenService {
	return &tokenService{}
}

func (s *tokenService) GenerateToken(authData *AuthData) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtHourExpiration := env.GetIntEnv("JWT_HOUR_EXPIRATION", 1)
	if jwtSecret == "" {
		return "", ErrJWTSecretNotSet
	}

	claims := &TokenData{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtHourExpiration) * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    authData.Id.String(),
			Subject:   authData.Username,
			ID:        uuid.NewString(),
			Audience:  []string{authData.Collaborator.CompanyID.String()},
		},
		Role:             authData.Role,
		CollaboratorID:   authData.Collaborator.ID.String(),
		CompanyID:        authData.Collaborator.CompanyID.String(),
		CollaboratorRole: authData.Collaborator.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	return signedToken, err
}

func (s *tokenService) ParseToken(tokenStr string) (*AuthData, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, ErrJWTSecretNotSet
	}
	token, err := jwt.ParseWithClaims(tokenStr, &TokenData{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*TokenData); ok {
		accountID, err := uuid.Parse(claims.Issuer)
		if err != nil {
			return nil, err
		}
		collabID, err := uuid.Parse(claims.CollaboratorID)
		if err != nil {
			return nil, err
		}
		companyID, err := uuid.Parse(claims.CompanyID)
		if err != nil {
			return nil, err
		}
		return &AuthData{
			Id:        accountID,
			Username:  claims.Subject,
			Role:      claims.Role,
			SessionID: claims.ID,
			Collaborator: CollaboratorData{
				ID:        collabID,
				CompanyID: companyID,
				Role:      claims.CollaboratorRole,
			},
		}, nil
	} else {
		return nil, ErrInvalidTokenClaims
	}
}
