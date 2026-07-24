package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicComplianceRole struct {
	Id               string     `json:"id"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	OrganizationUuid string     `json:"-"`
}

//// TABLE DEFINITION

func tableAnthropicComplianceRole(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_role",
		Description: "RBAC roles defined on each organization linked to the Claude Enterprise parent. Requires a Compliance Access Key with the read:compliance_org_data scope.",
		List: &plugin.ListConfig{
			ParentHydrate: listComplianceOrganizations,
			Hydrate:       listComplianceRoles,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique role identifier (rbac_role_...)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the role."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the role."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the role was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the role was last updated."},
			{Name: "organization_uuid", Type: proto.ColumnType_STRING, Description: "Uuid of the organization the role is defined on."},
		},
	}
}

//// LIST FUNCTION

func listComplianceRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := h.Item.(anthropicComplianceOrganization)

	client, err := getComplianceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_role.listComplianceRoles", "connection_error", err)
		return nil, err
	}

	path := fmt.Sprintf("/v1/compliance/organizations/%s/roles", org.Uuid)
	err = client.listAllPaged(ctx, path, nil, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var role anthropicComplianceRole
		if err := json.Unmarshal(item, &role); err != nil {
			return false, err
		}
		role.OrganizationUuid = org.Uuid
		d.StreamListItem(ctx, role)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_role.listComplianceRoles", "api_error", err)
		return nil, err
	}
	return nil, nil
}
