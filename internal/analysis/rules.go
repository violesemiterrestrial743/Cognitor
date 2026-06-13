package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type Rule interface {
	ID() string
	Evaluate(context.Context, model.SemanticChange) []model.Finding
}

type Engine struct {
	rules []Rule
}

func NewEngine(rules ...Rule) Engine {
	return Engine{rules: rules}
}

func DefaultEngine() Engine {
	return NewEngine(
		AccessCheckRule{},
		BoundsCheckRule{},
		NativeAPIRule{},
		HandleValidationRule{},
		ObjectLifetimeRule{},
		TokenFlowRule{},
		IOCTLRule{},
		RPCRule{},
		COMRule{},
		MarshallingRule{},
		ALPCRule{},
		RegistryRule{},
		ServiceRule{},
	)
}

func (e Engine) Evaluate(ctx context.Context, changes []model.SemanticChange) []model.Finding {
	var findings []model.Finding
	for _, change := range changes {
		if ctx.Err() != nil {
			return findings
		}
		for _, rule := range e.rules {
			findings = append(findings, rule.Evaluate(ctx, change)...)
		}
	}
	return findings
}
