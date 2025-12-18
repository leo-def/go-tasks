package httpx

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestResolvePaginationDefaults(t *testing.T) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request = httptest.NewRequest("GET", "/x", nil)
	p, err := ResolvePagination(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Skip != 0 || p.Limit != 20 || string(p.SortOrder) != "ASC" || p.SortBy != "id" {
		t.Fatalf("unexpected defaults: %+v", p)
	}
}

func TestResolvePaginationQuery(t *testing.T) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request = httptest.NewRequest("GET", "/x?skip=5&limit=10&sortBy=id&sortOrder=DESC&filter=name%3Afoo", nil)
	p, err := ResolvePagination(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Skip != 5 || p.Limit != 10 || string(p.SortOrder) != "DESC" || p.SortBy != "id" || p.Filter != "name:foo" {
		t.Fatalf("unexpected parsed values: %+v", p)
	}
}

func TestResolvePaginationErrors(t *testing.T) {
	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Request = httptest.NewRequest("GET", "/x?skip=-1", nil)
	if _, err := ResolvePagination(c); err == nil {
		t.Fatalf("expected error for negative skip")
	}
	c.Request = httptest.NewRequest("GET", "/x?limit=0", nil)
	if _, err := ResolvePagination(c); err == nil {
		t.Fatalf("expected error for non-positive limit")
	}
	c.Request = httptest.NewRequest("GET", "/x?sortOrder=BAD", nil)
	if _, err := ResolvePagination(c); err == nil {
		t.Fatalf("expected error for invalid sortOrder")
	}
}
