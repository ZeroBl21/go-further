package data

import (
	"strings"

	"github.com/ZeroBl21/go-further/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
	Sort     string

	SortSafeList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	// Page
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be maximum of 10 million")

	// Page Size
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page", "must be maximumo of 100")

	// Sort
	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

// Checks if the sort field matches one of the safelist. If does extract the
// column name fromo the sort field
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

// Return the sort direction depending on the prefix of the sort field
func (f Filters) sortDirection() string {
	if strings.Contains(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
