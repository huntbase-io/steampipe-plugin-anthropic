package anthropic

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

type anthropicAnalyticsUserActivity struct {
	User *struct {
		Id           string `json:"id"`
		EmailAddress string `json:"email_address"`
	} `json:"user"`
	ChatMetrics       json.RawMessage `json:"chat_metrics"`
	ClaudeCodeMetrics json.RawMessage `json:"claude_code_metrics"`
	CoworkMetrics     json.RawMessage `json:"cowork_metrics"`
	DesignMetrics     json.RawMessage `json:"design_metrics"`
	OfficeMetrics     json.RawMessage `json:"office_metrics"`
	ScienceMetrics    json.RawMessage `json:"science_metrics"`
	WebSearchCount    int64           `json:"web_search_count"`
	LastActivityDate  string          `json:"last_activity_date"`
	RbacGroupId       string          `json:"rbac_group_id"`
	RbacGroupName     string          `json:"rbac_group_name"`
}

//// TABLE DEFINITION

func tableAnthropicAnalyticsUserActivity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_analytics_user_activity",
		Description: "Per-user activity metrics across Claude products (chat, Claude Code, Cowork, design, office, science) from the Claude Enterprise Analytics API. Requires an Analytics API key. Defaults to the last 7 days; filter with the date column.",
		List: &plugin.ListConfig{
			Hydrate: listAnalyticsUserActivity,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "date", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "Date filter for the activity window (query parameter, not returned per row)."},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "Id of the user (user_...).", Transform: transform.FromField("User.Id")},
			{Name: "user_email", Type: proto.ColumnType_STRING, Description: "Email address of the user.", Transform: transform.FromField("User.EmailAddress")},
			{Name: "last_activity_date", Type: proto.ColumnType_STRING, Description: "Most recent date the user was active."},
			{Name: "web_search_count", Type: proto.ColumnType_INT, Description: "Web searches performed."},
			{Name: "rbac_group_id", Type: proto.ColumnType_STRING, Description: "RBAC group Id, when grouped."},
			{Name: "rbac_group_name", Type: proto.ColumnType_STRING, Description: "RBAC group name, when grouped."},
			{Name: "chat_metrics", Type: proto.ColumnType_JSON, Description: "Chat metrics: conversations, messages, projects, files, artifacts, skills, connectors."},
			{Name: "claude_code_metrics", Type: proto.ColumnType_JSON, Description: "Claude Code metrics: sessions, commits, pull requests, lines of code, tool actions."},
			{Name: "cowork_metrics", Type: proto.ColumnType_JSON, Description: "Cowork metrics: sessions, messages, actions, file edits, skills, plugins."},
			{Name: "design_metrics", Type: proto.ColumnType_JSON, Description: "Claude design metrics."},
			{Name: "office_metrics", Type: proto.ColumnType_JSON, Description: "Office agent metrics (excel, outlook, powerpoint, word)."},
			{Name: "science_metrics", Type: proto.ColumnType_JSON, Description: "Claude for science metrics."},
		},
	}
}

//// LIST FUNCTION

func listAnalyticsUserActivity(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAnalyticsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_user_activity.list", "connection_error", err)
		return nil, err
	}

	start, end := getDateRangeFromQuals(d, "date", 7)

	params := url.Values{}
	if start.Equal(end) {
		params.Set("date", start.Format("2006-01-02"))
	} else {
		params.Set("starting_date", start.Format("2006-01-02"))
		params.Set("ending_date", end.AddDate(0, 0, 1).Format("2006-01-02"))
	}

	err = client.listAllPaged(ctx, "/v1/organizations/analytics/users", params, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var activity anthropicAnalyticsUserActivity
		if err := json.Unmarshal(item, &activity); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, activity)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_user_activity.list", "api_error", err)
		return nil, err
	}
	return nil, nil
}
