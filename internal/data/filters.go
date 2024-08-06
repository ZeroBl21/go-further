package data

import "github.com/ZeroBl21/go-further/internal/validator"

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
