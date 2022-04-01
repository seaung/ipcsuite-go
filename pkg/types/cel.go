package types

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
)

type CelExpression struct {
	celEnvOptions     []cel.EvalOption
	celProgramOptions []cel.ProgramOption
}

func NewEnv() *CelExpression {
	return &CelExpression{}
}

func (c *CelExpression) CompileOptions() []cel.EvalOption {
	return c.celEnvOptions
}

func (c *CelExpression) ProgramOptions() []cel.ProgramOption {
	return c.celProgramOptions
}

func Evaluate(env *cel.Env, expression string, params map[string]interface{}) (ref.Val, error) {
	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return nil, issues.Err()
	}

	program, err := env.Program(ast)
	if err != nil {
		return nil, err
	}

	out, _, err := program.Eval(params)
	if err != nil {
		return nil, err
	}
	return out, nil
}
