package anthropic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

type anthropicWorkspaceMember struct {
	Type          string `json:"type"`
	UserId        string `json:"user_id"`
	WorkspaceId   string `json:"workspace_id"`
	WorkspaceRole string `json:"workspace_role"`
}

//// TABLE DEFINITION

func tableAnthropicWorkspaceMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_workspace_member",
		Description: "Members of each workspace in the Anthropic organization. Requires an Admin API key.",
		List: &plugin.ListConfig{
			ParentHydrate: listWorkspaces,
			Hydrate:       listWorkspaceMembers,
		},
		Columns: []*plugin.Column{
			{Name: "workspace_id", Type: proto.ColumnType_STRING, Description: "Id of the workspace."},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "Id of the organization member."},
			{Name: "workspace_role", Type: proto.ColumnType_STRING, Description: "Role of the member in the workspace: workspace_user, workspace_developer, workspace_admin, or workspace_billing.", Transform: transform.FromField("WorkspaceRole")},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Object type, always workspace_member."},
		},
	}
}

//// LIST FUNCTION

func listWorkspaceMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workspace := h.Item.(anthropicWorkspace)

	client, err := getAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_workspace_member.listWorkspaceMembers", "connection_error", err)
		return nil, err
	}

	path := fmt.Sprintf("/v1/organizations/workspaces/%s/members", workspace.Id)
	err = client.listAll(ctx, path, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var member anthropicWorkspaceMember
		if err := json.Unmarshal(item, &member); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, member)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_workspace_member.listWorkspaceMembers", "api_error", err)
		return nil, err
	}
	return nil, nil
}
