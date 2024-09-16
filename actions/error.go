package actions

import (
	"context"
	"fmt"

	"github.com/frops/condactor/rules"
)

var ErrorAction rules.ActionFunc = func(ctx context.Context) (context.Context, error) {
	return ctx, fmt.Errorf("error occurred")
}
