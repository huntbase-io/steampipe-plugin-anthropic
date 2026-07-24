package anthropic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicComplianceSettingsResponse struct {
	Type           string `json:"type"`
	OrganizationId string `json:"organization_id"`
	Settings       []struct {
		Name  string          `json:"name"`
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	} `json:"settings"`
}

type anthropicComplianceOrganizationSetting struct {
	OrganizationUuid string
	Name             string
	Type             string
	Value            json.RawMessage
}

//// TABLE DEFINITION

func tableAnthropicComplianceOrganizationSetting(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_organization_setting",
		Description: "Effective settings in force for each organization linked to the Claude Enterprise parent (retention, redaction, SSO, IP allowlist, etc.). Requires a Compliance Access Key with the read:compliance_org_data scope; the settings endpoint is enabled per parent organization.",
		List: &plugin.ListConfig{
			ParentHydrate: listComplianceOrganizations,
			Hydrate:       listComplianceOrganizationSettings,
		},
		Columns: []*plugin.Column{
			{Name: "organization_uuid", Type: proto.ColumnType_STRING, Description: "Uuid of the organization the setting applies to."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the setting, e.g. data_retention_periods, content_redaction_enabled."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the setting value: boolean, integer, string_list, provisioning_mode, or data_retention."},
			{Name: "value", Type: proto.ColumnType_JSON, Description: "Effective value of the setting."},
		},
	}
}

//// LIST FUNCTION

func listComplianceOrganizationSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := h.Item.(anthropicComplianceOrganization)

	client, err := getComplianceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_organization_setting.list", "connection_error", err)
		return nil, err
	}

	var resp anthropicComplianceSettingsResponse
	path := fmt.Sprintf("/v1/compliance/organizations/%s/settings", org.Uuid)
	if err := client.getJSON(ctx, path, nil, &resp); err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_organization_setting.list", "api_error", err)
		return nil, err
	}

	for _, s := range resp.Settings {
		d.StreamListItem(ctx, anthropicComplianceOrganizationSetting{
			OrganizationUuid: org.Uuid,
			Name:             s.Name,
			Type:             s.Type,
			Value:            s.Value,
		})
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
