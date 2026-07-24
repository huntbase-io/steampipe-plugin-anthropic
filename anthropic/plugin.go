// Package anthropic implements a Steampipe plugin for querying the Anthropic
// APIs: models, message batches and files (standard API key), organization
// administration (Admin API key), Claude Enterprise compliance data
// (Compliance Access Key), and Claude usage analytics (Admin and Analytics
// API keys).
package anthropic

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

// Plugin creates this (anthropic) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-anthropic",
		DefaultTransform: transform.FromCamel(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"anthropic_analytics_cost_report":           tableAnthropicAnalyticsCostReport(ctx),
			"anthropic_analytics_summary":               tableAnthropicAnalyticsSummary(ctx),
			"anthropic_analytics_usage_report":          tableAnthropicAnalyticsUsageReport(ctx),
			"anthropic_analytics_user_activity":         tableAnthropicAnalyticsUserActivity(ctx),
			"anthropic_analytics_user_cost_report":      tableAnthropicAnalyticsUserCostReport(ctx),
			"anthropic_analytics_user_usage_report":     tableAnthropicAnalyticsUserUsageReport(ctx),
			"anthropic_api_key":                         tableAnthropicApiKey(ctx),
			"anthropic_claude_code_analytics":           tableAnthropicClaudeCodeAnalytics(ctx),
			"anthropic_compliance_activity":             tableAnthropicComplianceActivity(ctx),
			"anthropic_compliance_chat":                 tableAnthropicComplianceChat(ctx),
			"anthropic_compliance_chat_message":         tableAnthropicComplianceChatMessage(ctx),
			"anthropic_compliance_group":                tableAnthropicComplianceGroup(ctx),
			"anthropic_compliance_group_member":         tableAnthropicComplianceGroupMember(ctx),
			"anthropic_compliance_organization":         tableAnthropicComplianceOrganization(ctx),
			"anthropic_compliance_organization_setting": tableAnthropicComplianceOrganizationSetting(ctx),
			"anthropic_compliance_project":              tableAnthropicComplianceProject(ctx),
			"anthropic_compliance_project_attachment":   tableAnthropicComplianceProjectAttachment(ctx),
			"anthropic_compliance_role":                 tableAnthropicComplianceRole(ctx),
			"anthropic_compliance_user":                 tableAnthropicComplianceUser(ctx),
			"anthropic_file":                            tableAnthropicFile(ctx),
			"anthropic_invite":                          tableAnthropicInvite(ctx),
			"anthropic_message_batch":                   tableAnthropicMessageBatch(ctx),
			"anthropic_model":                           tableAnthropicModel(ctx),
			"anthropic_organization_member":             tableAnthropicOrganizationMember(ctx),
			"anthropic_workspace":                       tableAnthropicWorkspace(ctx),
			"anthropic_workspace_member":                tableAnthropicWorkspaceMember(ctx),
		},
	}
	return p
}
