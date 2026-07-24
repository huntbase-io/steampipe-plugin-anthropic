package anthropic

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

// Project fields are passed through as raw JSON plus common extracted fields,
// since the full project schema varies.
type anthropicComplianceProject struct {
	Id               string          `json:"id"`
	Name             string          `json:"name"`
	CreatedAt        *string         `json:"created_at"`
	UpdatedAt        *string         `json:"updated_at"`
	OrganizationUuid string          `json:"organization_uuid"`
	Raw              json.RawMessage `json:"-"`
}

//// TABLE DEFINITION

func tableAnthropicComplianceProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_project",
		Description: "claude.ai projects across the Claude Enterprise parent organization. Requires a Compliance Access Key with the read:compliance_user_data scope.",
		List: &plugin.ListConfig{
			Hydrate: listComplianceProjects,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique project identifier (claude_proj_...)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the project."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the project was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the project was last updated."},
			{Name: "organization_uuid", Type: proto.ColumnType_STRING, Description: "Uuid of the organization the project belongs to."},
			{Name: "project", Type: proto.ColumnType_JSON, Description: "Full raw project record.", Transform: transform.FromField("Raw")},
		},
	}
}

//// LIST FUNCTION

func listComplianceProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getComplianceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_project.listComplianceProjects", "connection_error", err)
		return nil, err
	}

	err = client.listAllPaged(ctx, "/v1/compliance/apps/projects", nil, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var project anthropicComplianceProject
		if err := json.Unmarshal(item, &project); err != nil {
			return false, err
		}
		project.Raw = item
		d.StreamListItem(ctx, project)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_project.listComplianceProjects", "api_error", err)
		return nil, err
	}
	return nil, nil
}
