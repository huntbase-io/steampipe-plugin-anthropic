package anthropic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

type anthropicComplianceOrganization struct {
	Uuid      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
}

//// TABLE DEFINITION

func tableAnthropicComplianceOrganization(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_organization",
		Description: "Organizations linked to the Claude Enterprise parent organization. Requires a Compliance Access Key with the read:compliance_org_data scope.",
		List: &plugin.ListConfig{
			Hydrate: listComplianceOrganizations,
		},
		Columns: []*plugin.Column{
			{Name: "uuid", Type: proto.ColumnType_STRING, Description: "Canonical organization identifier.", Transform: transform.FromField("Uuid")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the organization."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the organization was created."},
		},
	}
}

//// LIST FUNCTION

func listComplianceOrganizations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getComplianceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_organization.listComplianceOrganizations", "connection_error", err)
		return nil, err
	}

	err = client.listAllPaged(ctx, "/v1/compliance/organizations", nil, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var org anthropicComplianceOrganization
		if err := json.Unmarshal(item, &org); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, org)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_organization.listComplianceOrganizations", "api_error", err)
		return nil, err
	}
	return nil, nil
}
