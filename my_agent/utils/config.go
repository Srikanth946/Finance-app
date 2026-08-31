package utils

// ToolConfig holds configuration for tool execution behaviors.
type ToolConfig struct {
	RequireConfirmation bool
}

// DefaultToolConfig provides a standard configuration.
func DefaultToolConfig() *ToolConfig {
	return &ToolConfig{
		RequireConfirmation: false,
	}
}
