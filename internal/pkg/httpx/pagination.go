package httpx

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Query parameter keys centralized here; change them once to affect all endpoints
const (
	ParamSkip      = "skip"
	ParamLimit     = "limit"
	ParamSortBy    = "sortBy"
	ParamSortOrder = "sortOrder"
	ParamFilter    = "filter"
)

type SortField string

// SortOrder allowed values
const (
	SortASC  SortField = "ASC"
	SortDESC SortField = "DESC"
)

// PaginationParams models common list query parameters
// SortBy can be "field.subField"; SortOrder is ASC or DESC
// Filter meaning is endpoint-specific but passed through consistently
type PaginationParams struct {
	Skip      int
	Limit     int
	SortBy    string
	SortOrder string
	Filter    string
}

// ResolvePagination parses query parameters with defaults and validation
// Defaults: skip=0, limit=20, sortOrder=ASC; sortBy and filter optional
func ResolvePagination(ctx *gin.Context) (PaginationParams, error) {
    sortBy := ctx.Query(ParamSortBy)
    if sortBy == "" {
        sortBy = "id"
    }
    q := PaginationParams{
        Skip:      0,
        Limit:     20,
        SortBy:    sortBy,
        SortOrder: string(SortASC),
        Filter:    ctx.Query(ParamFilter),
    }
	if s := ctx.Query(ParamSkip); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 0 {
			return q, fmt.Errorf("invalid %s: %s", ParamSkip, s)
		}
		q.Skip = v
	}
	if s := ctx.Query(ParamLimit); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v <= 0 {
			return q, fmt.Errorf("invalid %s: %s", ParamLimit, s)
		}
		q.Limit = v
	}
	if order := ctx.Query(ParamSortOrder); order != "" {
		if order != string(SortASC) && order != string(SortDESC) {
			return q, fmt.Errorf("invalid %s: %s", ParamSortOrder, order)
		}
		q.SortOrder = order
	}
	return q, nil
}

func ResolveUintParam(ctx *gin.Context, param string) (uint, error) {
	s := ctx.Param(param)
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %s", param, s)
	}
	return uint(v), nil
}

// New: ResolveUUIDParam parses a path param as UUID with a consistent error message.
func ResolveUUIDParam(ctx *gin.Context, param string) (uuid.UUID, error) {
	s := ctx.Param(param)
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
