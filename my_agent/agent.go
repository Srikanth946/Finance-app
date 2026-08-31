package main

import (
	"context"
	"finance_app/my_agent/models"
	"finance_app/my_agent/tools"
	"finance_app/my_agent/utils"

	"log"
	"os"

	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/cmd/launcher"
	"google.golang.org/adk/v2/cmd/launcher/full"
	"google.golang.org/adk/v2/tool"
)

func main() {

	ctx := context.Background()
	model := models.NewGeminiModel(ctx)

	transactionsTool := NewFinanceTool(
		"GetTransactions",
		"Fetches all financial transactions for the current user directly from the system. Requires no data from user while calling.",
		func(ctx agent.Context, args map[string]any) (map[string]any, error) {
			res, err := tools.GetTransactions(ctx)
			if err != nil {
				return nil, err
			}
			return map[string]any{"result": res}, nil
		},
	)

	dashboardTool := NewFinanceTool(
		"GetDashboardSummary",
		"Fetches the current user's financial dashboard summary (previous month interest, total amount) directly from the system. Requires no arguments.",
		func(ctx agent.Context, args map[string]any) (map[string]any, error) {
			res, err := tools.GetDashboardSummary(ctx)
			if err != nil {
				return nil, err
			}
			return map[string]any{"result": res}, nil
		},
	)

	// Tool configuration
	toolCfg := utils.DefaultToolConfig()

	// Wrap tools in a Toolset and apply confirmation logic if required
	financeToolset := NewFinanceToolset(transactionsTool, dashboardTool)
	var finalTools []tool.Tool

	if toolCfg.RequireConfirmation {
		confirmedToolset := tool.WithConfirmation(financeToolset, true, nil)
		var err error
		finalTools, err = confirmedToolset.Tools(nil)
		if err != nil {
			log.Fatalf("error extracting confirmed tools: %v", err)
		}
	} else {
		var err error
		finalTools, err = financeToolset.Tools(nil)
		if err != nil {
			log.Fatalf("error extracting tools: %v", err)
		}
	}

	FinanceAgent, err := llmagent.New(llmagent.Config{
		Name:        "Finance Agent",
		Description: "An agent that can answer questions about finance and investments.",
		Instruction: "You are a helpful assistant that can answer questions about finance and investments. You have direct access to the user's personal financial data via the provided tools. When the user asks about their transactions, balance, or dashboard, you MUST use the provided tools to fetch the data before answering. Do not tell the user you cannot access private accounts, because you have tools specifically for that purpose. Give answers in a concise and clear only by using received data from tools.",
		Model:       model,
		Tools:       finalTools,
	},
	)

	if err != nil {
		log.Fatalf("error in agent creation - %s", err.Error())
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(FinanceAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}

}
