package anthropic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicWorkspace struct {
	Id           string     `json:"id"`
	Type         string     `json:"type"`
	Name         string     `json:"name"`
	CreatedAt    *time.Time `json:"created_at"`
	ArchivedAt   *time.Time `json:"archived_at"`
	DisplayColor string     `json:"display_color"`
}

//// TABLE DEFINITION

func tableAnthropicWorkspace(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_workspace",
		Description: "Workspaces in the Anthropic organization. Requires an Admin API key.",
		List: &plugin.ListConfig{
			Hydrate: listWorkspaces,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique workspace identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the workspace."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the workspace was created."},
			{Name: "archived_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the workspace was archived, if archived."},
			{Name: "display_color", Type: proto.ColumnType_STRING, Description: "Hex color code shown in the Anthropic Console."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Object type, always workspace."},
		},
	}
}

//// LIST FUNCTION

func listWorkspaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_workspace.listWorkspaces", "connection_error", err)
		return nil, err
	}

	err = client.listAll(ctx, "/v1/organizations/workspaces", d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var ws anthropicWorkspace
		if err := json.Unmarshal(item, &ws); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, ws)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_workspace.listWorkspaces", "api_error", err)
		return nil, err
	}
	return nil, nil
}
