package main

import (
	"context"
	"log"

	"github.com/frops/condactor/actions"
	"github.com/frops/condactor/rules"
	"github.com/google/uuid"
)

type se struct {
}

func (se *se) LoadRules() ([]rules.Rule, error) {
	idempotency := rules.Rule{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Name:      "Check dpd",
		Condition: `true`,
		Action:    actions.IdempotencyAction,
	}
	dpdRule := rules.Rule{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		Name:      "Check dpd",
		Condition: `dpd <= 0`,
		Action: func(ctx context.Context) (context.Context, error) {
			log.Println("dpd is less than or equal to 0")
			return ctx, nil
		},
		Next: &idempotency,
	}
	delinquencyRule := rules.Rule{
		ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Name:      "Check delinquency",
		Condition: `delinquency > 0`,
		Action: func(ctx context.Context) (context.Context, error) {
			log.Println("delinquency is greater than 0")
			return ctx, nil
		},
		Next: &dpdRule,
	}

	return []rules.Rule{
		delinquencyRule,
	}, nil
}
