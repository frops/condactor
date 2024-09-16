package actions

import (
	"context"
	"fmt"

	"github.com/frops/condactor/rules"
)

// ActionConfig defines the structure for action configurations
type ActionConfig struct {
	URL   string `json:"url"`   // Example config for fetching data
	Value string `json:"value"` // Example config for logging
}

// CreateAction creates an ActionFunc based on the action type
func CreateAction(actionType string, config ActionConfig) rules.ActionFunc {
	switch actionType {
	case "fetch_data":
		return func(ctx context.Context) (context.Context, error) {
			// Example fetch_data action
			fmt.Printf("Fetching data from %s\n", config.URL)
			return ctx, nil
		}
	case "log_message":
		return func(ctx context.Context) (context.Context, error) {
			// Example log_message action
			fmt.Printf("Logging value: %s\n", config.Value)
			return ctx, nil
		}
	default:
		return func(ctx context.Context) (context.Context, error) {
			return ctx, nil
		}
	}
}
