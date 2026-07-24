package anthropic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type anthropicApiKeyCreatedBy struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type anthropicApiKey struct {
	Id             string                   `json:"id"`
	Type           string                   `json:"type"`
	Name           string                   `json:"name"`
	WorkspaceId    *string                  `json:"workspace_id"`
	CreatedAt      *time.Time               `json:"created_at"`
	CreatedBy      anthropicApiKeyCreatedBy `json:"created_by"`
	PartialKeyHint string                   `json:"partial_key_hint"`
	Status         string                   `json:"status"`
}

//// TABLE DEFINITION

func tableAnthropicApiKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "anthropic_api_key",
		Description: "API keys in the Anthropic organization. Requires an Admin API key.",
		List: &plugin.ListConfig{
			Hydrate: listApiKeys,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique API key identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the API key."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the API key: active, inactive, or archived."},
			{Name: "workspace_id", Type: proto.ColumnType_STRING, Description: "Id of the workspace the key belongs to, null for the default workspace."},
			{Name: "partial_key_hint", Type: proto.ColumnType_STRING, Description: "Partially redacted hint of the API key value."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time at which the API key was created."},
			{Name: "created_by", Type: proto.ColumnType_JSON, Description: "The organization member who created the API key."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Object type, always api_key."},
		},
	}
}

//// LIST FUNCTION

func listApiKeys(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getAdminClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_api_key.listApiKeys", "connection_error", err)
		return nil, err
	}

	err = client.listAll(ctx, "/v1/organizations/api_keys", d.QueryContext.GetLimit(), func(item json.RawMessage) (bool, error) {
		var key anthropicApiKey
		if err := json.Unmarshal(item, &key); err != nil {
			return false, err
		}
		d.StreamListItem(ctx, key)
		return d.RowsRemaining(ctx) != 0, nil
	})
	if err != nil {
		plugin.Logger(ctx).Error("anthropic_api_key.listApiKeys", "api_error", err)
		return nil, err
	}
	return nil, nil
}
