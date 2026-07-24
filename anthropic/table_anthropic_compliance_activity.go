package anthropic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

type anthropicActivityActor struct {
	Type         string `json:"type"`
	EmailAddress string `json:"email_address"`
	UserId       string `json:"user_id"`
	IpAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
}

type anthropicActivity struct {
	Id               string                  `json:"id"`
	CreatedAt        *time.Time              `json:"created_at"`
	OrganizationId   string                  `json:"organization_id"`
	OrganizationUuid string                  `json:"organization_uuid"`
	Actor            *anthropicActivityActor `json:"actor"`
	Type             string                  `json:"type"`
	// Raw holds the full event payload including event-type-specific fields
	// (claude_chat_id, claude_project_id, etc.).
	Raw json.RawMessage `json:"-"`
}

//// TABLE DEFINITION

func tableAnthropicComplianceActivity(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_activity",
		Description: "Activity Feed events for the organization from the Compliance API. Accepts a Compliance Access Key or an Admin API key.",
		List: &plugin.ListConfig{
			Hydrate: listComplianceActivities,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique activity event identifier."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Activity event type, e.g. claude_chat_created."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the event occurred."},
			{Name: "organization_id", Type: proto.ColumnType_STRING, Description: "org_-prefixed organization identifier."},
			{Name: "organization_uuid", Type: proto.ColumnType_STRING, Description: "Uuid of the organization, joinable with anthropic_compliance_organization."},
			{Name: "actor_type", Type: proto.ColumnType_STRING, Description: "Type of the actor that generated the event.", Transform: transform.FromField("Actor.Type")},
			{Name: "actor_email", Type: proto.ColumnType_STRING, Description: "Email address of the acting user.", Transform: transform.FromField("Actor.EmailAddress")},
			{Name: "actor_user_id", Type: proto.ColumnType_STRING, Description: "User Id of the acting user.", Transform: transform.FromField("Actor.UserId")},
			{Name: "actor_ip_address", Type: proto.ColumnType_STRING, Description: "IP address of the actor.", Transform: transform.FromField("Actor.IpAddress")},
			{Name: "actor_user_agent", Type: proto.ColumnType_STRING, Description: "User agent of the actor.", Transform: transform.FromField("Actor.UserAgent")},
			{Name: "actor", Type: proto.ColumnType_JSON, Description: "Full actor object."},
			{Name: "event", Type: proto.ColumnType_JSON, Description: "Full raw event payload including event-type-specific fields such as claude_chat_id and claude_project_id.", Transform: transform.FromField("Raw")},
		},
	}
}

//// LIST FUNCTION

func listComplianceActivities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getActivityFeedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_activity.listComplianceActivities", "connection_error", err)
		return nil, err
	}

	err = client.listAll(ctx, "/v1/compliance/activities", d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var activity anthropicActivity
		if err := json.Unmarshal(item, &activity); err != nil {
			return false, err
		}
		activity.Raw = item
		d.StreamListItem(ctx, activity)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_activity.listComplianceActivities", "api_error", err)
		return nil, err
	}
	return nil, nil
}
