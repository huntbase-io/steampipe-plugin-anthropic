package anthropic

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicAnalyticsUsageRow struct {
	StartingAt           string          `json:"-"`
	EndingAt             string          `json:"-"`
	Model                string          `json:"model"`
	Product              string          `json:"product"`
	ContextWindow        string          `json:"context_window"`
	InferenceGeo         string          `json:"inference_geo"`
	Speed                string          `json:"speed"`
	RbacGroupId          string          `json:"rbac_group_id"`
	Requests             int64           `json:"requests"`
	UncachedInputTokens  int64           `json:"uncached_input_tokens"`
	OutputTokens         int64           `json:"output_tokens"`
	CacheReadInputTokens int64           `json:"cache_read_input_tokens"`
	CacheCreation        json.RawMessage `json:"cache_creation"`
	ServerToolUse        json.RawMessage `json:"server_tool_use"`
}

//// TABLE DEFINITION

func tableAnthropicAnalyticsUsageReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_analytics_usage_report",
		Description: "Daily token usage buckets by model and product from the Claude Enterprise Analytics API. Requires an Analytics API key. Defaults to the last 7 days; filter with the date column (max 31 days per query).",
		List: &plugin.ListConfig{
			Hydrate: listAnalyticsUsageReport,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "date", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "Start of the usage bucket (UTC day).", Transform: fromBucketStart()},
			{Name: "ending_at", Type: proto.ColumnType_TIMESTAMP, Description: "End of the usage bucket."},
			{Name: "model", Type: proto.ColumnType_STRING, Description: "Claude model the usage is attributed to."},
			{Name: "product", Type: proto.ColumnType_STRING, Description: "Claude product: chat, claude_code, cowork, office_agent, claude_in_chrome, claude_design, claude-in-slack."},
			{Name: "requests", Type: proto.ColumnType_INT, Description: "Number of requests."},
			{Name: "uncached_input_tokens", Type: proto.ColumnType_INT, Description: "Uncached input tokens."},
			{Name: "output_tokens", Type: proto.ColumnType_INT, Description: "Output tokens."},
			{Name: "cache_read_input_tokens", Type: proto.ColumnType_INT, Description: "Input tokens served from cache."},
			{Name: "cache_creation", Type: proto.ColumnType_JSON, Description: "Cache-write token counts by TTL."},
			{Name: "server_tool_use", Type: proto.ColumnType_JSON, Description: "Server tool usage counts, e.g. web_search_requests."},
			{Name: "context_window", Type: proto.ColumnType_STRING, Description: "Context window tier: 0-200k or 200k-1M."},
			{Name: "inference_geo", Type: proto.ColumnType_STRING, Description: "Inference region: global, us, or not_available."},
			{Name: "speed", Type: proto.ColumnType_STRING, Description: "Inference speed: standard or fast."},
			{Name: "rbac_group_id", Type: proto.ColumnType_STRING, Description: "RBAC group ID, when grouped."},
		},
	}
}

//// LIST FUNCTION

func listAnalyticsUsageReport(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAnalyticsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_usage_report.list", "connection_error", err)
		return nil, err
	}

	start, end := getDateRangeFromQuals(d, "date", 7)

	params := url.Values{}
	params.Set("starting_at", start.Format(time.RFC3339))
	params.Set("ending_at", end.AddDate(0, 0, 1).Format(time.RFC3339))
	params.Set("bucket_width", "1d")
	params.Set("limit", "31")
	params.Add("group_by[]", "model")
	params.Add("group_by[]", "product")

	err = client.listAnalyticsBuckets(ctx, "/v1/organizations/analytics/usage_report", params, func(bucket analyticsBucket, item json.RawMessage) (bool, error) {
		var row anthropicAnalyticsUsageRow
		if err := json.Unmarshal(item, &row); err != nil {
			return false, err
		}
		row.StartingAt = bucket.StartingAt
		row.EndingAt = bucket.EndingAt
		d.StreamListItem(ctx, row)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_usage_report.list", "api_error", err)
		return nil, err
	}
	return nil, nil
}
