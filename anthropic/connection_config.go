package anthropic

import (
	"context"
	"errors"
	"os"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicConfig struct {
	// Standard API key (sk-ant-api...) used for data-plane endpoints:
	// models, message batches, files.
	APIKey *string `hcl:"api_key"`

	// Admin API key (sk-ant-admin...) used for organization endpoints:
	// workspaces, members, invites, API keys.
	AdminAPIKey *string `hcl:"admin_api_key"`

	// Compliance Access Key (sk-ant-api01..., created in claude.ai) used for
	// the Compliance API endpoints: activities, organizations, users, roles,
	// groups, chats, projects. Claude Enterprise only.
	ComplianceAPIKey *string `hcl:"compliance_api_key"`

	// Analytics API key (created in claude.ai > Organization settings > API by
	// the primary owner, read:analytics scope) used for the Claude Enterprise
	// Analytics API endpoints under /v1/organizations/analytics/.
	AnalyticsAPIKey *string `hcl:"analytics_api_key"`
}

// ConfigInstance returns a new, empty connection configuration struct.
func ConfigInstance() interface{} {
	return &anthropicConfig{}
}

// GetConfig retrieves and casts the connection configuration.
func GetConfig(connection *plugin.Connection) anthropicConfig {
	if connection == nil {
		return anthropicConfig{}
	}
	config, _ := connection.GetConfig().(anthropicConfig)
	return config
}

// getApiKey resolves the standard API key from config or environment.
func getApiKey(_ context.Context, d *plugin.QueryData) (string, error) {
	config := GetConfig(d.Connection)
	if config.APIKey != nil && *config.APIKey != "" {
		return *config.APIKey, nil
	}
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return key, nil
	}
	return "", errors.New("api_key must be set in the connection config or the ANTHROPIC_API_KEY environment variable must be set")
}

// getAdminApiKey resolves the Admin API key from config or environment.
func getAdminApiKey(_ context.Context, d *plugin.QueryData) (string, error) {
	config := GetConfig(d.Connection)
	if config.AdminAPIKey != nil && *config.AdminAPIKey != "" {
		return *config.AdminAPIKey, nil
	}
	if key := os.Getenv("ANTHROPIC_ADMIN_KEY"); key != "" {
		return key, nil
	}
	return "", errors.New("admin_api_key must be set in the connection config or the ANTHROPIC_ADMIN_KEY environment variable must be set (requires an sk-ant-admin... key)")
}

// getComplianceApiKey resolves the Compliance Access Key from config or environment.
func getComplianceApiKey(_ context.Context, d *plugin.QueryData) (string, error) {
	config := GetConfig(d.Connection)
	if config.ComplianceAPIKey != nil && *config.ComplianceAPIKey != "" {
		return *config.ComplianceAPIKey, nil
	}
	if key := os.Getenv("ANTHROPIC_COMPLIANCE_ACCESS_KEY"); key != "" {
		return key, nil
	}
	return "", errors.New("compliance_api_key must be set in the connection config or the ANTHROPIC_COMPLIANCE_ACCESS_KEY environment variable must be set (requires a Compliance Access Key created in claude.ai; Claude Enterprise only)")
}

// getAnalyticsClient returns a client for the Claude Enterprise Analytics API.
func getAnalyticsClient(_ context.Context, d *plugin.QueryData) (*adminClient, error) {
	config := GetConfig(d.Connection)
	if config.AnalyticsAPIKey != nil && *config.AnalyticsAPIKey != "" {
		return newAdminClient(*config.AnalyticsAPIKey), nil
	}
	if key := os.Getenv("ANTHROPIC_ANALYTICS_API_KEY"); key != "" {
		return newAdminClient(key), nil
	}
	return nil, errors.New("analytics_api_key must be set in the connection config or the ANTHROPIC_ANALYTICS_API_KEY environment variable must be set (requires an Analytics API key created in claude.ai; Claude Enterprise only)")
}

// getComplianceClient returns a client for the Anthropic Compliance API.
// The activity feed also accepts an Admin API key, so fall back to that when
// no Compliance Access Key is configured.
func getComplianceClient(ctx context.Context, d *plugin.QueryData) (*adminClient, error) {
	key, err := getComplianceApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	return newAdminClient(key), nil
}

// getActivityFeedClient prefers the Compliance Access Key but falls back to
// the Admin API key, which is also accepted by GET /v1/compliance/activities.
func getActivityFeedClient(ctx context.Context, d *plugin.QueryData) (*adminClient, error) {
	if key, err := getComplianceApiKey(ctx, d); err == nil {
		return newAdminClient(key), nil
	}
	if key, err := getAdminApiKey(ctx, d); err == nil {
		return newAdminClient(key), nil
	}
	return nil, errors.New("the compliance activity feed requires compliance_api_key (Compliance Access Key) or admin_api_key (sk-ant-admin...) to be set in the connection config, or ANTHROPIC_COMPLIANCE_ACCESS_KEY / ANTHROPIC_ADMIN_KEY in the environment")
}

// getAdminClient returns a client for the Anthropic Admin API.
func getAdminClient(ctx context.Context, d *plugin.QueryData) (*adminClient, error) {
	key, err := getAdminApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	return newAdminClient(key), nil
}

// getClient returns an Anthropic SDK client for data-plane endpoints.
func getClient(ctx context.Context, d *plugin.QueryData) (*sdk.Client, error) {
	key, err := getApiKey(ctx, d)
	if err != nil {
		return nil, err
	}
	client := sdk.NewClient(option.WithAPIKey(key))
	return &client, nil
}
