package rules

import (
	"context"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/uuid"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Rule defines the structure for a rule with a condition, action, and an optional next rule
type Rule struct {
	ID        uuid.UUID
	Name      string
	Condition string
	Action    ActionFunc
	Next      *Rule
}

// ActionFunc defines a function type for actions with support for dependencies and context
type ActionFunc func(ctx context.Context) (context.Context, error)

// RuleEngine manages the execution of rules
type RuleEngine struct {
	storage Storage
	env     *cel.Env
}

// NewRuleEngine creates a new instance of RuleEngine
func NewRuleEngine(storage Storage, decl ...*exprpb.Decl) (*RuleEngine, error) {
	env, err := cel.NewEnv(cel.Declarations(decl...))
	if err != nil {
		return nil, fmt.Errorf("failed to create CEL environment: %w", err)
	}

	return &RuleEngine{
		storage: storage,
		env:     env,
	}, nil
}

// Option defines a functional option for configuring rule execution
type Option func(*execOptions)

// execOptions holds execution settings
type execOptions struct {
	atomic       bool
	returnErrors bool
}

// WithAtomic sets atomic execution: the rule execution stops on the first error
func WithAtomic(atomic bool) Option {
	return func(opts *execOptions) {
		opts.atomic = atomic
	}
}

// WithReturnErrors enables error return via channel
func WithReturnErrors(returnErrors bool) Option {
	return func(opts *execOptions) {
		opts.returnErrors = returnErrors
	}
}

// ExecuteRules executes the rules with given input data, context, and options
func (re *RuleEngine) ExecuteRules(ctx context.Context, inputData map[string]interface{}, opts ...Option) <-chan error {
	errCh := make(chan error)
	options := &execOptions{}

	// Apply functional options
	for _, opt := range opts {
		opt(options)
	}

	go func() {
		defer close(errCh)

		rules, err := re.storage.LoadRules()
		if err != nil {
			errCh <- fmt.Errorf("failed to load rules: %w", err)
			return
		}

		for _, rule := range rules {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				if err = re.executeRule(ctx, &rule, inputData, errCh, options); err != nil && options.atomic {
					return
				}
			}
		}
	}()

	return errCh
}

// executeRule recursively evaluates and executes the current rule and its next rules
func (re *RuleEngine) executeRule(ctx context.Context, rule *Rule, arg map[string]interface{}, errCh chan<- error, options *execOptions) error {
	// Compile the condition
	ast, issues := re.env.Compile(rule.Condition)
	if issues != nil && issues.Err() != nil {
		err := fmt.Errorf("failed to compile for rule '%s' condition: %v", rule.Name, issues.Err())
		errCh <- err
		return err
	}

	// Create the program for execution
	prg, err := re.env.Program(ast)
	if err != nil {
		err = fmt.Errorf("failed to create program: %v", err)
		errCh <- err
		return err
	}

	// Evaluate the condition with the input data
	out, _, err := prg.Eval(arg)
	if err != nil {
		err = fmt.Errorf("failed to evaluate condition: %v", err)
		errCh <- err
		return err
	}

	// If the condition is true, execute the action
	if out == types.True {
		ctx, err = rule.Action(ctx)
		if err != nil {
			errCh <- fmt.Errorf("action failed: %v", err)
			return err
		}

		// Recursively execute the next rule if available
		if rule.Next != nil {
			if err = re.executeRule(ctx, rule.Next, arg, errCh, options); err != nil {
				return err
			}
		}
	}

	return nil
}
