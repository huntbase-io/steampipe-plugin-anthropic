package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicComplianceGroupMember struct {
	UserId    string     `json:"user_id"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	GroupId   string     `json:"-"`
}

//// TABLE DEFINITION

func tableAnthropicComplianceGroupMember(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_group_member",
		Description: "Members of each group in the Claude Enterprise parent organization. Requires a Compliance Access Key with the read:compliance_org_data and read:compliance_user_data scopes.",
		List: &plugin.ListConfig{
			ParentHydrate: listComplianceGroups,
			Hydrate:       listComplianceGroupMembers,
		},
		Columns: []*plugin.Column{
			{Name: "group_id", Type: proto.ColumnType_STRING, Description: "Id of the group."},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "Id of the member user (user_...)."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email address of the member."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the membership was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the membership was last updated."},
		},
	}
}

//// LIST FUNCTION

func listComplianceGroupMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	group := h.Item.(anthropicComplianceGroup)

	client, err := getComplianceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_group_member.listComplianceGroupMembers", "connection_error", err)
		return nil, err
	}

	path := fmt.Sprintf("/v1/compliance/groups/%s/members", group.Id)
	err = client.listAllPaged(ctx, path, nil, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var member anthropicComplianceGroupMember
		if err := json.Unmarshal(item, &member); err != nil {
			return false, err
		}
		member.GroupId = group.Id
		d.StreamListItem(ctx, member)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_group_member.listComplianceGroupMembers", "api_error", err)
		return nil, err
	}
	return nil, nil
}
