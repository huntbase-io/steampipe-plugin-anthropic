package anthropic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicComplianceGroup struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SourceType  string     `json:"source_type"`
	Roles       []string   `json:"roles"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

//// TABLE DEFINITION

func tableAnthropicComplianceGroup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_group",
		Description: "RBAC and SCIM-provisioned groups in the Claude Enterprise parent organization. Requires a Compliance Access Key with the read:compliance_org_data scope.",
		List: &plugin.ListConfig{
			Hydrate: listComplianceGroups,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique group identifier (rbac_group_...)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the group."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the group."},
			{Name: "source_type", Type: proto.ColumnType_STRING, Description: "How the group is managed: direct (created in claude.ai) or scim (synced from an identity provider)."},
			{Name: "roles", Type: proto.ColumnType_JSON, Description: "Role IDs assigned to the group."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the group was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the group was last updated."},
		},
	}
}

//// LIST FUNCTION

func listComplianceGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getComplianceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_group.listComplianceGroups", "connection_error", err)
		return nil, err
	}

	err = client.listAllPaged(ctx, "/v1/compliance/groups", nil, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var group anthropicComplianceGroup
		if err := json.Unmarshal(item, &group); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, group)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_group.listComplianceGroups", "api_error", err)
		return nil, err
	}
	return nil, nil
}
