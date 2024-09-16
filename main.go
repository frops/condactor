package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/frops/condactor/rules"
	"github.com/google/cel-go/checker/decls"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func main() {
	vars := []*exprpb.Decl{
		decls.NewVar("delinquency", decls.Int),
		decls.NewVar("dpd", decls.Int),
	}

	engine, err := rules.NewRuleEngine(&se{}, vars...)
	if err != nil {
		log.Fatalf("Failed to create rule engine: %v", err)
	}

	// Example input data
	inputData := map[string]interface{}{
		"delinquency": 17,
		"dpd":         -1,
	}

	// Create a cancellable context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute rules with options: atomic execution and error channel enabled
	errCh := engine.ExecuteRules(
		ctx,
		inputData,
		rules.WithAtomic(true),
		rules.WithReturnErrors(true),
	)

	// Handle errors from channel
	for err = range errCh {
		if err != nil {
			fmt.Printf("Error occurred: %v\n", err)
		}
	}
}
