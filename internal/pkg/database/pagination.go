package database

import (
    "fmt"

    "go-tasks/internal/pkg/httpx"
    "gorm.io/gorm"
)

// Paginate applies count, sort, and page to the given query.
// - T is the destination type (model struct)
// - model is the table model reference (e.g., &Company{})
// - apply is an optional filter hook to extend the base query
func Paginate[T any](db *gorm.DB, model any, p httpx.PaginationParams, apply func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
    var items []T
    var count int64

    q := db.Model(model)
    if apply != nil {
        q = apply(q)
    }

    if err := q.Count(&count).Error; err != nil {
        return nil, 0, err
    }

    if p.SortBy != "" {
        order := fmt.Sprintf("%s %s", p.SortBy, p.SortOrder)
        q = q.Order(order)
    }

    if err := q.Offset(p.Skip).Limit(p.Limit).Find(&items).Error; err != nil {
        return nil, 0, err
    }
    return items, count, nil
}