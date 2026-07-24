package anthropic

import (
	"context"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAnthropicMessageBatch(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_message_batch",
		Description: "Message batches created in the workspace of the API key.",
		List: &plugin.ListConfig{
			Hydrate: listMessageBatches,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique batch identifier.", Transform: transform.FromField("ID")},
			{Name: "processing_status", Type: proto.ColumnType_STRING, Description: "Processing status: in_progress, canceling, or ended."},
			{Name: "request_counts", Type: proto.ColumnType_JSON, Description: "Counts of requests by result: processing, succeeded, errored, canceled, expired."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the batch was created."},
			{Name: "ended_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which processing ended."},
			{Name: "expires_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the batch will expire and end processing."},
			{Name: "archived_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the batch was archived."},
			{Name: "cancel_initiated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which cancellation was initiated."},
			{Name: "results_url", Type: proto.ColumnType_STRING, Description: "URL to download the results of the batch.", Transform: transform.FromField("ResultsURL")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Object type, always message_batch."},
		},
	}
}

//// LIST FUNCTION

func listMessageBatches(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_message_batch.listMessageBatches", "connection_error", err)
		return nil, err
	}

	iter := client.Messages.Batches.ListAutoPaging(ctx, sdk.MessageBatchListParams{})
	for iter.Next() {
		d.StreamListItem(ctx, iter.Current())
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		plugin.Logger(ctx).Error("anthropic_message_batch.listMessageBatches", "api_error", err)
		return nil, err
	}
	return nil, nil
}
