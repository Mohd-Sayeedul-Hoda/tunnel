package request

import (
	"context"
)

type Validator interface {
	Valid(ctx context.Context) map[string]string
}
