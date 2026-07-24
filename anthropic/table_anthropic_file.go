package anthropic

import (
	"context"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

//// TABLE DEFINITION

func tableAnthropicFile(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_file",
		Description: "Files uploaded to the Anthropic Files API (beta).",
		List: &plugin.ListConfig{
			Hydrate: listFiles,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique file identifier.", Transform: transform.FromField("ID")},
			{Name: "filename", Type: proto.ColumnType_STRING, Description: "Original filename of the uploaded file."},
			{Name: "mime_type", Type: proto.ColumnType_STRING, Description: "MIME type of the file."},
			{Name: "size_bytes", Type: proto.ColumnType_INT, Description: "Size of the file in bytes."},
			{Name: "downloadable", Type: proto.ColumnType_BOOL, Description: "Whether the file can be downloaded."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the file was uploaded."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Object type, always file."},
		},
	}
}

//// LIST FUNCTION

func listFiles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_file.listFiles", "connection_error", err)
		return nil, err
	}

	iter := client.Beta.Files.ListAutoPaging(ctx, sdk.BetaFileListParams{
		Betas: []sdk.AnthropicBeta{sdk.AnthropicBetaFilesAPI2025_04_14},
	})
	for iter.Next() {
		d.StreamListItem(ctx, iter.Current())
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		plugin.Logger(ctx).Error("anthropic_file.listFiles", "api_error", err)
		return nil, err
	}
	return nil, nil
}
