package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicComplianceProjectAttachment struct {
	Id        string     `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	Filename  string     `json:"filename"`
	MimeType  string     `json:"mime_type"`
	Type      string     `json:"type"`
	ProjectId string     `json:"-"`
}

//// TABLE DEFINITION

func tableAnthropicComplianceProjectAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_compliance_project_attachment",
		Description: "Files and documents attached to claude.ai projects. Requires a Compliance Access Key with the read:compliance_user_data scope.",
		List: &plugin.ListConfig{
			ParentHydrate: listComplianceProjects,
			Hydrate:       listComplianceProjectAttachments,
		},
		Columns: []*plugin.Column{
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "Id of the project the attachment belongs to."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Attachment identifier: claude_file_... for binary uploads, claude_proj_doc_... for text documents."},
			{Name: "filename", Type: proto.ColumnType_STRING, Description: "Filename of the attachment."},
			{Name: "mime_type", Type: proto.ColumnType_STRING, Description: "MIME type of the attachment."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Attachment kind: project_file (binary upload) or project_doc (plain-text document)."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the attachment was added."},
		},
	}
}

//// LIST FUNCTION

func listComplianceProjectAttachments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project := h.Item.(anthropicComplianceProject)

	client, err := getComplianceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_project_attachment.list", "connection_error", err)
		return nil, err
	}

	path := fmt.Sprintf("/v1/compliance/apps/projects/%s/attachments", project.Id)
	err = client.listAllPaged(ctx, path, nil, d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var attachment anthropicComplianceProjectAttachment
		if err := json.Unmarshal(item, &attachment); err != nil {
			return false, err
		}
		attachment.ProjectId = project.Id
		d.StreamListItem(ctx, attachment)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_compliance_project_attachment.list", "api_error", err)
		return nil, err
	}
	return nil, nil
}
