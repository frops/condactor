package actions

import (
	"context"
	"fmt"
	"log"

	"github.com/frops/condactor/rules"
)

var IdempotencyAction rules.ActionFunc = func(ctx context.Context) (context.Context, error) {
	customerID := ctx.Value("customer_id")
	if customerID == nil {
		return ctx, fmt.Errorf("customer_id not found in context")
	}

	ik := fmt.Sprintf("customer_id:%d", ctx.Value("customer_id").(int))

	log.Println("identifying customer with key:", ik)
	return ctx, nil
}
