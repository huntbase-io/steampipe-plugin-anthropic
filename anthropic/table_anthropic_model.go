package anthropic

import (
	"context"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAnthropicModel(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_model",
		Description: "Claude models available to the API key.",
		List: &plugin.ListConfig{
			Hydrate: listModels,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique model identifier, e.g. claude-opus-4-8.", Transform: transform.FromField("ID")},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "Human-readable name of the model."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the model was released."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Object type, always model."},
		},
	}
}

//// LIST FUNCTION

func listModels(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_model.listModels", "connection_error", err)
		return nil, err
	}

	iter := client.Models.ListAutoPaging(ctx, sdk.ModelListParams{})
	for iter.Next() {
		d.StreamListItem(ctx, iter.Current())
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		plugin.Logger(ctx).Error("anthropic_model.listModels", "api_error", err)
		return nil, err
	}
	return nil, nil
}
