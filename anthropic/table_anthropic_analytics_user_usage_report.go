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

type anthropicAnalyticsActor struct {
	Type    string `json:"type"`
	UserId  string `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
}

type anthropicAnalyticsUserUsageRow struct {
	Actor                *anthropicAnalyticsActor `json:"actor"`
	StartingAt           *time.Time               `json:"starting_at"`
	EndingAt             *time.Time               `json:"ending_at"`
	Model                string                   `json:"model"`
	Product              string                   `json:"product"`
	ContextWindow        string                   `json:"context_window"`
	InferenceGeo         string                   `json:"inference_geo"`
	Speed                string                   `json:"speed"`
	RbacGroupId          string                   `json:"rbac_group_id"`
	Requests             int64                    `json:"requests"`
	TotalTokens          int64                    `json:"total_tokens"`
	UncachedInputTokens  int64                    `json:"uncached_input_tokens"`
	OutputTokens         int64                    `json:"output_tokens"`
	CacheReadInputTokens int64                    `json:"cache_read_input_tokens"`
	CacheCreation        json.RawMessage          `json:"cache_creation"`
	ServerToolUse        json.RawMessage          `json:"server_tool_use"`
}

//// TABLE DEFINITION

func tableAnthropicAnalyticsUserUsageReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_analytics_user_usage_report",
		Description: "Per-user token usage ranking from the Claude Enterprise Analytics API. Requires an Analytics API key. Defaults to the last 7 days; filter with the date column.",
		List: &plugin.ListConfig{
			Hydrate: listAnalyticsUserUsageReport,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "date", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "Start of the reporting window.", Transform: transform.FromField("StartingAt")},
			{Name: "ending_at", Type: proto.ColumnType_TIMESTAMP, Description: "End of the reporting window."},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "ID of the user.", Transform: transform.FromField("Actor.UserId")},
			{Name: "user_email", Type: proto.ColumnType_STRING, Description: "Email address of the user.", Transform: transform.FromField("Actor.Email")},
			{Name: "user_name", Type: proto.ColumnType_STRING, Description: "Name of the user.", Transform: transform.FromField("Actor.Name")},
			{Name: "user_deleted", Type: proto.ColumnType_BOOL, Description: "Whether the user account has been deleted.", Transform: transform.FromField("Actor.Deleted")},
			{Name: "total_tokens", Type: proto.ColumnType_INT, Description: "Total tokens consumed."},
			{Name: "requests", Type: proto.ColumnType_INT, Description: "Number of requests."},
			{Name: "uncached_input_tokens", Type: proto.ColumnType_INT, Description: "Uncached input tokens."},
			{Name: "output_tokens", Type: proto.ColumnType_INT, Description: "Output tokens."},
			{Name: "cache_read_input_tokens", Type: proto.ColumnType_INT, Description: "Input tokens served from cache."},
			{Name: "cache_creation", Type: proto.ColumnType_JSON, Description: "Cache-write token counts by TTL."},
			{Name: "server_tool_use", Type: proto.ColumnType_JSON, Description: "Server tool usage counts."},
			{Name: "model", Type: proto.ColumnType_STRING, Description: "Model, when grouped."},
			{Name: "product", Type: proto.ColumnType_STRING, Description: "Product, when grouped."},
			{Name: "context_window", Type: proto.ColumnType_STRING, Description: "Context window tier, when grouped."},
			{Name: "inference_geo", Type: proto.ColumnType_STRING, Description: "Inference region, when grouped."},
			{Name: "speed", Type: proto.ColumnType_STRING, Description: "Inference speed, when grouped."},
			{Name: "rbac_group_id", Type: proto.ColumnType_STRING, Description: "RBAC group ID, when grouped."},
		},
	}
}

//// LIST FUNCTION

func listAnalyticsUserUsageReport(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAnalyticsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_user_usage_report.list", "connection_error", err)
		return nil, err
	}

	start, end := getDateRangeFromQuals(d, "date", 7)

	params := url.Values{}
	params.Set("starting_at", start.Format(time.RFC3339))
	params.Set("ending_at", end.AddDate(0, 0, 1).Format(time.RFC3339))

	err = client.listAllPaged(ctx, "/v1/organizations/analytics/user_usage_report", params, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var row anthropicAnalyticsUserUsageRow
		if err := json.Unmarshal(item, &row); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, row)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_user_usage_report.list", "api_error", err)
		return nil, err
	}
	return nil, nil
}
