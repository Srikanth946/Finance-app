package tools

import (
	"fmt"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/toolutils"
	"google.golang.org/genai"
)

type financeTool struct {
	name        string
	description string
	declaration *genai.FunctionDeclaration
	call        func(ctx agent.Context, args map[string]any) (map[string]any, error)
}

func (t *financeTool) Name() string {
	return t.name
}

func (t *financeTool) Description() string {
	return t.description
}

func (t *financeTool) Declaration() *genai.FunctionDeclaration {
	return t.declaration
}

func (t *financeTool) IsLongRunning() bool {
	return false
}

func (t *financeTool) Run(ctx agent.Context, args any) (map[string]any, error) {
	fmt.Printf("[DEBUG] Tool %s called with args: %v\n", t.name, args)

	// Cast args to map[string]any to match the call signature
	toolArgs, ok := args.(map[string]any)
	if !ok {
		if args == nil {
			toolArgs = make(map[string]any)
		} else {
			fmt.Printf("[ERROR] Tool %s received invalid args type: %T\n", t.name, args)
			return nil, fmt.Errorf("invalid arguments type: %T", args)
		}
	}

	result, err := t.call(ctx, toolArgs)
	if err != nil {
		fmt.Printf("[ERROR] Tool %s failed: %v\n", t.name, err)
		return nil, err
	}
	return result, nil
}

func (t *financeTool) ProcessRequest(ctx agent.Context, req *model.LLMRequest) error {
	fmt.Printf("[DEBUG] Packing tool %s into request\n", t.name)
	return toolutils.PackTool(req, t)
}

func NewFinanceTool(name, description string, call func(agent.Context, map[string]any) (map[string]any, error)) tool.Tool {
	return &financeTool{
		name:        name,
		description: description,
		declaration: &genai.FunctionDeclaration{
			Name:        name,
			Description: description,
			Parameters:  &genai.Schema{Type: genai.TypeObject, Properties: make(map[string]*genai.Schema)},
		},
		call: call,
	}
}

type FinanceToolset struct {
	tools []tool.Tool
}

func (s *FinanceToolset) Name() string {
	return "FinanceToolset"
}

func (s *FinanceToolset) Tools(ctx agent.ReadonlyContext) ([]tool.Tool, error) {
	return s.tools, nil
}

func NewFinanceToolset(tools ...tool.Tool) tool.Toolset {
	return &FinanceToolset{tools: tools}
}
