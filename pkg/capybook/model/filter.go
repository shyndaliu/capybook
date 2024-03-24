package model

import (
	"strings"

	"github.com/shyndaliu/capybook/pkg/capybook/validator"
)

type Filters struct {
	Page         int
	Limit        int
	Sort         string
	SortSafelist []string
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.Limit
}
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

// Return the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.Limit > 0, "limit", "must be greater than zero")
	v.Check(f.Limit <= 100, "limit", "must be a maximum of 100")
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
