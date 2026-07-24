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

// fromBucketStart maps the "date" column to the bucket's starting_at value.
func fromBucketStart() *transform.ColumnTransforms {
	return transform.FromField("StartingAt")
}

type anthropicAnalyticsCostRow struct {
	StartingAt    string `json:"-"`
	EndingAt      string `json:"-"`
	Amount        string `json:"amount"`
	ListAmount    string `json:"list_amount"`
	Currency      string `json:"currency"`
	CostType      string `json:"cost_type"`
	TokenType     string `json:"token_type"`
	Model         string `json:"model"`
	Product       string `json:"product"`
	ContextWindow string `json:"context_window"`
	InferenceGeo  string `json:"inference_geo"`
	Speed         string `json:"speed"`
	RbacGroupId   string `json:"rbac_group_id"`
	Requests      int64  `json:"requests"`
}

//// TABLE DEFINITION

func tableAnthropicAnalyticsCostReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_analytics_cost_report",
		Description: "Daily cost buckets by model, product, and cost type from the Claude Enterprise Analytics API. Requires an Analytics API key. Defaults to the last 7 days; filter with the date column (max 31 days per query).",
		List: &plugin.ListConfig{
			Hydrate: listAnalyticsCostReport,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "date", Require: plugin.Optional, Operators: []string{"=", ">", ">=", "<", "<="}},
			},
		},
		Columns: []*plugin.Column{
			{Name: "date", Type: proto.ColumnType_TIMESTAMP, Description: "Start of the cost bucket (UTC day).", Transform: fromBucketStart()},
			{Name: "ending_at", Type: proto.ColumnType_TIMESTAMP, Description: "End of the cost bucket."},
			{Name: "amount", Type: proto.ColumnType_STRING, Description: "Cost as a decimal string in cents USD, e.g. 41280.000000 = $412.80."},
			{Name: "list_amount", Type: proto.ColumnType_STRING, Description: "List-price cost as a decimal string in cents USD."},
			{Name: "currency", Type: proto.ColumnType_STRING, Description: "Currency code, currently USD."},
			{Name: "cost_type", Type: proto.ColumnType_STRING, Description: "Cost category: tokens, web_search, or code_execution."},
			{Name: "token_type", Type: proto.ColumnType_STRING, Description: "Token category, when grouped by token_type."},
			{Name: "model", Type: proto.ColumnType_STRING, Description: "Claude model the cost is attributed to."},
			{Name: "product", Type: proto.ColumnType_STRING, Description: "Claude product the cost is attributed to."},
			{Name: "requests", Type: proto.ColumnType_INT, Description: "Number of requests."},
			{Name: "context_window", Type: proto.ColumnType_STRING, Description: "Context window tier."},
			{Name: "inference_geo", Type: proto.ColumnType_STRING, Description: "Inference region."},
			{Name: "speed", Type: proto.ColumnType_STRING, Description: "Inference speed: standard or fast."},
			{Name: "rbac_group_id", Type: proto.ColumnType_STRING, Description: "RBAC group ID, when grouped."},
		},
	}
}

//// LIST FUNCTION

func listAnalyticsCostReport(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAnalyticsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_cost_report.list", "connection_error", err)
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
	params.Add("group_by[]", "cost_type")

	err = client.listAnalyticsBuckets(ctx, "/v1/organizations/analytics/cost_report", params, func(bucket analyticsBucket, item json.RawMessage) (bool, error) {
		var row anthropicAnalyticsCostRow
		if err := json.Unmarshal(item, &row); err != nil {
			return false, err
		}
		row.StartingAt = bucket.StartingAt
		row.EndingAt = bucket.EndingAt
		d.StreamListItem(ctx, row)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_analytics_cost_report.list", "api_error", err)
		return nil, err
	}
	return nil, nil
}
