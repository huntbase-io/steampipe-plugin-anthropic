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

type anthropicAnalyticsSummary struct {
	StartingAt              *time.Time `json:"starting_at"`
	EndingAt                *time.Time `json:"ending_at"`
	AssignedSeatCount       int64      `json:"assigned_seat_count"`
	PendingInviteCount      int64      `json:"pending_invite_count"`
	DailyActiveUserCount    int64      `json:"daily_active_user_count"`
	WeeklyActiveUserCount   int64      `json:"weekly_active_user_count"`
	MonthlyActiveUserCount  int64      `json:"monthly_active_user_count"`
	DailyAdoptionRate       float64    `json:"daily_adoption_rate"`
	WeeklyAdoptionRate      float64    `json:"weekly_adoption_rate"`
	MonthlyAdoptionRate     float64    `json:"monthly_adoption_rate"`
	ChatDailyActiveUsers    int64      `json:"chat_daily_active_user_count"`
	ClaudeCodeDailyActive   int64      `json:"claude_code_daily_active_user_count"`
	ClaudeCodeWeeklyActive  int64      `json:"claude_code_weekly_active_user_count"`
	ClaudeCodeMonthlyActive int64      `json:"claude_code_monthly_active_user_count"`
	Raw                     json.RawMessage
}

type anthropicAnalyticsSummariesResponse struct {
	Summaries []json.RawMessage `json:"summaries"`
}

//// TABLE DEFINITION

func tableAnthropicAnalyticsSummary(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_analytics_summary",
		Description: "Organization-wide daily active users, seats, and adoption rates from the Claude Enterprise Analytics API. Requires an Analytics API key. Defaults to the last 30 days; filter with the date column.",
		List: &plugin.ListConfig{
			Hydrate: listAnalyticsSummaries,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "date", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "Start of the UTC day the summary covers.", Transform: transform.FromField("StartingAt")},
			{Name: "assigned_seat_count", Type: proto.ColumnType_INT, Description: "Number of assigned seats."},
			{Name: "pending_invite_count", Type: proto.ColumnType_INT, Description: "Number of pending invites."},
			{Name: "daily_active_user_count", Type: proto.ColumnType_INT, Description: "Daily active users across all Claude products."},
			{Name: "weekly_active_user_count", Type: proto.ColumnType_INT, Description: "Weekly active users."},
			{Name: "monthly_active_user_count", Type: proto.ColumnType_INT, Description: "Monthly active users."},
			{Name: "daily_adoption_rate", Type: proto.ColumnType_DOUBLE, Description: "Daily active users as a fraction of assigned seats."},
			{Name: "weekly_adoption_rate", Type: proto.ColumnType_DOUBLE, Description: "Weekly adoption rate."},
			{Name: "monthly_adoption_rate", Type: proto.ColumnType_DOUBLE, Description: "Monthly adoption rate."},
			{Name: "chat_daily_active_user_count", Type: proto.ColumnType_INT, Description: "Daily active users in chat.", Transform: transform.FromField("ChatDailyActiveUsers")},
			{Name: "claude_code_daily_active_user_count", Type: proto.ColumnType_INT, Description: "Daily active users in Claude Code.", Transform: transform.FromField("ClaudeCodeDailyActive")},
			{Name: "claude_code_weekly_active_user_count", Type: proto.ColumnType_INT, Description: "Weekly active users in Claude Code.", Transform: transform.FromField("ClaudeCodeWeeklyActive")},
			{Name: "claude_code_monthly_active_user_count", Type: proto.ColumnType_INT, Description: "Monthly active users in Claude Code.", Transform: transform.FromField("ClaudeCodeMonthlyActive")},
			{Name: "summary", Type: proto.ColumnType_JSON, Description: "Full raw summary record including per-product active-user counts (cowork, design, office agent, science).", Transform: transform.FromField("Raw")},
		},
	}
}

//// LIST FUNCTION

func listAnalyticsSummaries(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAnalyticsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_summary.list", "connection_error", err)
		return nil, err
	}

	start, end := getDateRangeFromQuals(d, "date", 30)

	params := url.Values{}
	params.Set("starting_date", start.Format("2006-01-02"))
	// ending_date is exclusive.
	params.Set("ending_date", end.AddDate(0, 0, 1).Format("2006-01-02"))

	var resp anthropicAnalyticsSummariesResponse
	if err := client.getJSON(ctx, "/v1/organizations/analytics/summaries", params, &resp); err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_summary.list", "api_error", err)
		return nil, err
	}

	for _, item := range resp.Summaries {
		var summary anthropicAnalyticsSummary
		if err := json.Unmarshal(item, &summary); err != nil {
			return nil, err
		}
		summary.Raw = item
		d.StreamListItem(ctx, summary)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
