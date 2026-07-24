package anthropic

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

type claudeCodeActor struct {
	Type         string `json:"type"`
	EmailAddress string `json:"email_address"`
	ApiKeyName   string `json:"api_key_name"`
}

type claudeCodeAnalyticsRecord struct {
	Date           *time.Time       `json:"date"`
	Actor          *claudeCodeActor `json:"actor"`
	OrganizationId string           `json:"organization_id"`
	CustomerType   string           `json:"customer_type"`
	TerminalType   string           `json:"terminal_type"`
	CoreMetrics    *struct {
		NumSessions int64 `json:"num_sessions"`
		LinesOfCode struct {
			Added   int64 `json:"added"`
			Removed int64 `json:"removed"`
		} `json:"lines_of_code"`
		CommitsByClaudeCode      int64 `json:"commits_by_claude_code"`
		PullRequestsByClaudeCode int64 `json:"pull_requests_by_claude_code"`
	} `json:"core_metrics"`
	ToolActions    json.RawMessage `json:"tool_actions"`
	ModelBreakdown json.RawMessage `json:"model_breakdown"`
}

//// TABLE DEFINITION

func tableAnthropicClaudeCodeAnalytics(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_claude_code_analytics",
		Description: "Daily aggregated Claude Code productivity metrics per user. Requires an Admin API key. Defaults to the last 7 days; filter with the date column.",
		List: &plugin.ListConfig{
			Hydrate: listClaudeCodeAnalytics,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "date", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "UTC day the metrics cover."},
			{Name: "actor_type", Type: proto.ColumnType_STRING, Description: "How the user authenticated: user_actor (OAuth) or api_actor (API key).", Transform: transform.FromField("Actor.Type")},
			{Name: "actor_email", Type: proto.ColumnType_STRING, Description: "Email address of the user (user_actor).", Transform: transform.FromField("Actor.EmailAddress")},
			{Name: "actor_api_key_name", Type: proto.ColumnType_STRING, Description: "API key name (api_actor).", Transform: transform.FromField("Actor.ApiKeyName")},
			{Name: "organization_id", Type: proto.ColumnType_STRING, Description: "Organization Uuid."},
			{Name: "customer_type", Type: proto.ColumnType_STRING, Description: "api (pay-as-you-go) or subscription (Pro/Team)."},
			{Name: "terminal_type", Type: proto.ColumnType_STRING, Description: "Terminal or environment where Claude Code was used, e.g. vscode, iTerm.app, tmux."},
			{Name: "num_sessions", Type: proto.ColumnType_INT, Description: "Number of distinct Claude Code sessions.", Transform: transform.FromField("CoreMetrics.NumSessions")},
			{Name: "lines_added", Type: proto.ColumnType_INT, Description: "Lines of code added by Claude Code.", Transform: transform.FromField("CoreMetrics.LinesOfCode.Added")},
			{Name: "lines_removed", Type: proto.ColumnType_INT, Description: "Lines of code removed by Claude Code.", Transform: transform.FromField("CoreMetrics.LinesOfCode.Removed")},
			{Name: "commits", Type: proto.ColumnType_INT, Description: "Git commits created through Claude Code.", Transform: transform.FromField("CoreMetrics.CommitsByClaudeCode")},
			{Name: "pull_requests", Type: proto.ColumnType_INT, Description: "Pull requests created through Claude Code.", Transform: transform.FromField("CoreMetrics.PullRequestsByClaudeCode")},
			{Name: "tool_actions", Type: proto.ColumnType_JSON, Description: "Accepted/rejected counts per tool (edit_tool, multi_edit_tool, write_tool, notebook_edit_tool)."},
			{Name: "model_breakdown", Type: proto.ColumnType_JSON, Description: "Per-model token counts and estimated cost (cents USD)."},
		},
	}
}

//// LIST FUNCTION

func listClaudeCodeAnalytics(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_claude_code_analytics.list", "connection_error", err)
		return nil, err
	}

	start, end := getDateRangeFromQuals(d, "date", 7)

	// The endpoint returns one UTC day per request; iterate the range.
	for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {
		params := url.Values{}
		params.Set("starting_at", day.Format("2006-01-02"))
		err = client.listAllPaged(ctx, "/v1/organizations/usage_report/claude_code", params, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
			var rec claudeCodeAnalyticsRecord
			if err := json.Unmarshal(item, &rec); err != nil {
				return false, err
			}
			d.StreamListItem(ctx, rec)
			return d.RowsRemaining(ctx) != 0, nil
		})
		if err != nil {
			plugin.Logger(ctx).Error("anthropic_claude_code_analytics.list", "api_error", err)
			return nil, err
		}
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
