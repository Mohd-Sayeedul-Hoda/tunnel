package request

import (
	"context"
)

type Filters struct {
	Page  int
	Limit int
}

func (f *Filters) Valid(ctx context.Context, v *Valid) *Valid {

	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.Limit > 0, "limit", "must be greater than zero")
	v.Check(f.Limit <= 100, "limit", "must be a maximum of 100")

	return v
}
