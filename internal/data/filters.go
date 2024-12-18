package data

import (
	"strings"

	"github.com/kharljhon14/greenlight/internal/validator"
)

type Filter struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func ValidateFilters(v *validator.Validator, f Filter) {
	v.Check(f.Page > 0, "page", "must be grater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Check that the sort parameter matches a value in the safelist
	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

func (f Filter) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filter) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filter) limit() int {
	return f.PageSize
}

func (f Filter) offset() int {
	return (f.Page - 1) * f.PageSize
}
