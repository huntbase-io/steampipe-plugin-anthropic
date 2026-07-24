package anthropic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

type anthropicComplianceChatUser struct {
	Id           string `json:"id"`
	EmailAddress string `json:"email_address"`
}

type anthropicComplianceChat struct {
	Id               string                       `json:"id"`
	Name             string                       `json:"name"`
	CreatedAt        *time.Time                   `json:"created_at"`
	UpdatedAt        *time.Time                   `json:"updated_at"`
	DeletedAt        *time.Time                   `json:"deleted_at"`
	Href             string                       `json:"href"`
	Model            string                       `json:"model"`
	OrganizationUuid string                       `json:"organization_uuid"`
	ProjectId        *string                      `json:"project_id"`
	User             *anthropicComplianceChatUser `json:"user"`
}

//// TABLE DEFINITION

func tableAnthropicComplianceChat(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_chat",
		Description: "claude.ai chats across the Claude Enterprise parent organization. Requires a Compliance Access Key with the read:compliance_user_data scope.",
		List: &plugin.ListConfig{
			Hydrate: listComplianceChats,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique chat identifier (claude_chat_...)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the chat."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the chat was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the chat was last updated."},
			{Name: "deleted_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the chat was soft-deleted in claude.ai, if deleted."},
			{Name: "href", Type: proto.ColumnType_STRING, Description: "URL of the chat in claude.ai."},
			{Name: "model", Type: proto.ColumnType_STRING, Description: "Model used in the chat."},
			{Name: "organization_uuid", Type: proto.ColumnType_STRING, Description: "Uuid of the organization the chat belongs to."},
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "Id of the project the chat belongs to, if any."},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "Id of the user who owns the chat.", Transform: transform.FromField("User.Id")},
			{Name: "user_email", Type: proto.ColumnType_STRING, Description: "Email address of the user who owns the chat.", Transform: transform.FromField("User.EmailAddress")},
		},
	}
}

//// LIST FUNCTION

func listComplianceChats(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getComplianceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_chat.listComplianceChats", "connection_error", err)
		return nil, err
	}

	err = client.listAll(ctx, "/v1/compliance/apps/chats", d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var chat anthropicComplianceChat
		if err := json.Unmarshal(item, &chat); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, chat)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_chat.listComplianceChats", "api_error", err)
		return nil, err
	}
	return nil, nil
}
