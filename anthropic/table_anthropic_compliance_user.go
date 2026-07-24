package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicComplianceUser struct {
	Id               string     `json:"id"`
	FullName         string     `json:"full_name"`
	Email            string     `json:"email"`
	OrganizationRole string     `json:"organization_role"`
	CreatedAt        *time.Time `json:"created_at"`
	OrganizationUuid string     `json:"-"`
}

//// TABLE DEFINITION

func tableAnthropicComplianceUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_user",
		Description: "Users in each organization linked to the Claude Enterprise parent. Requires a Compliance Access Key with the read:compliance_org_data and read:compliance_user_data scopes.",
		List: &plugin.ListConfig{
			ParentHydrate: listComplianceOrganizations,
			Hydrate:       listComplianceUsers,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique user identifier (user_...)."},
			{Name: "full_name", Type: proto.ColumnType_STRING, Description: "Full name of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address of the user."},
			{Name: "organization_role", Type: proto.ColumnType_STRING, Description: "Built-in membership level: admin, billing, claude_code_user, developer, managed, membership_admin, owner, primary_owner, or user."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the user joined the organization."},
			{Name: "organization_uuid", Type: proto.ColumnType_STRING, Description: "Uuid of the organization the user belongs to."},
		},
	}
}

//// LIST FUNCTION

func listComplianceUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := h.Item.(anthropicComplianceOrganization)

	client, err := getComplianceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_user.listComplianceUsers", "connection_error", err)
		return nil, err
	}

	path := fmt.Sprintf("/v1/compliance/organizations/%s/users", org.Uuid)
	err = client.listAllPaged(ctx, path, nil, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var user anthropicComplianceUser
		if err := json.Unmarshal(item, &user); err != nil {
			return false, err
		}
		user.OrganizationUuid = org.Uuid
		d.StreamListItem(ctx, user)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_user.listComplianceUsers", "api_error", err)
		return nil, err
	}
	return nil, nil
}
