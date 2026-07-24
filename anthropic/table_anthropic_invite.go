package anthropic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicInvite struct {
	Id        string     `json:"id"`
	Type      string     `json:"type"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	InvitedAt *time.Time `json:"invited_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

//// TABLE DEFINITION

func tableAnthropicInvite(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_invite",
		Description: "Pending and historical invitations to the Anthropic organization. Requires an Admin API key.",
		List: &plugin.ListConfig{
			Hydrate: listInvites,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique invite identifier."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address of the invited user."},
			{Name: "role", Type: proto.ColumnType_STRING, Description: "Organization role granted on acceptance: user, developer, billing, or admin."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the invite: pending, accepted, expired, or deleted."},
			{Name: "invited_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the invite was sent."},
			{Name: "expires_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the invite expires."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Object type, always invite."},
		},
	}
}

//// LIST FUNCTION

func listInvites(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_invite.listInvites", "connection_error", err)
		return nil, err
	}

	err = client.listAll(ctx, "/v1/organizations/invites", d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var invite anthropicInvite
		if err := json.Unmarshal(item, &invite); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, invite)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_invite.listInvites", "api_error", err)
		return nil, err
	}
	return nil, nil
}
