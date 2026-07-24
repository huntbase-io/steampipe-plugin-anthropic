package anthropic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicOrganizationMember struct {
	Id      string     `json:"id"`
	Type    string     `json:"type"`
	Email   string     `json:"email"`
	Name    string     `json:"name"`
	Role    string     `json:"role"`
	AddedAt *time.Time `json:"added_at"`
}

//// TABLE DEFINITION

func tableAnthropicOrganizationMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_organization_member",
		Description: "Users in the Anthropic organization. Requires an Admin API key.",
		List: &plugin.ListConfig{
			Hydrate: listOrganizationMembers,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique user identifier."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address of the user."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the user."},
			{Name: "role", Type: proto.ColumnType_STRING, Description: "Organization role: user, developer, billing, or admin."},
			{Name: "added_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the user joined the organization."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Object type, always user."},
		},
	}
}

//// LIST FUNCTION

func listOrganizationMembers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_organization_member.listOrganizationMembers", "connection_error", err)
		return nil, err
	}

	err = client.listAll(ctx, "/v1/organizations/users", d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var user anthropicOrganizationMember
		if err := json.Unmarshal(item, &user); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, user)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_organization_member.listOrganizationMembers", "api_error", err)
		return nil, err
	}
	return nil, nil
}
