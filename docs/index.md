---
organization: Huntbase
category: ["ai"]
brand_color: "#D97757"
display_name: "Anthropic"
short_name: "anthropic"
description: "Steampipe plugin to query models, organization administration, compliance data, and usage analytics from the Anthropic API."
og_description: "Query the Anthropic API with SQL! Open source CLI. No DB required."
---

# Anthropic + Steampipe

[Anthropic](https://www.anthropic.com) is an AI safety company and the maker of Claude. The [Anthropic API](https://platform.claude.com/docs/en/api/) exposes models, organization administration, Claude Enterprise compliance data, and usage analytics.

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

For example:

```sql
select
  id,
  display_name,
  created_at
from
  anthropic_model;
```

```
+------------------+-----------------+----------------------+
| id               | display_name    | created_at           |
+------------------+-----------------+----------------------+
| claude-opus-4-8  | Claude Opus 4.8 | 2026-05-01T00:00:00Z |
| claude-sonnet-5  | Claude Sonnet 5 | 2026-04-01T00:00:00Z |
+------------------+-----------------+----------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/huntbase/anthropic/tables)**

## Quick start

### Install

Build from source and install as a local plugin:

```sh
make install
cp config/anthropic.spc ~/.steampipe/config/anthropic.spc
```

### Credentials

The plugin uses up to four independent credentials; configure any subset — each unlocks its own group of tables.

| Item | Description |
|---|---|
| Standard API key | `sk-ant-api...` from the [Claude Console](https://platform.claude.com/settings/keys). Tables: `anthropic_model`, `anthropic_message_batch`, `anthropic_file`. |
| Admin API key | `sk-ant-admin...` from [Console Admin keys](https://platform.claude.com/settings/admin-keys). Tables: `anthropic_workspace*`, `anthropic_organization_member`, `anthropic_invite`, `anthropic_api_key`, `anthropic_claude_code_analytics`, and `anthropic_compliance_activity`. |
| Compliance Access Key | Created in claude.ai (Claude Enterprise only). Tables: `anthropic_compliance_*`. Scopes: `read:compliance_activities`, `read:compliance_org_data`, `read:compliance_user_data`. |
| Analytics API key | Created in claude.ai by the primary owner (Claude Enterprise only, `read:analytics` scope). Tables: `anthropic_analytics_*`. |

### Configuration

Installing the latest anthropic plugin will create a config file (`~/.steampipe/config/anthropic.spc`) with a single connection named `anthropic`:

```hcl
connection "anthropic" {
  plugin = "local/anthropic"

  # Standard API key. Env: ANTHROPIC_API_KEY
  # api_key = "sk-ant-api03-..."

  # Admin API key. Env: ANTHROPIC_ADMIN_KEY
  # admin_api_key = "sk-ant-admin01-..."

  # Compliance Access Key (Claude Enterprise). Env: ANTHROPIC_COMPLIANCE_ACCESS_KEY
  # compliance_api_key = "sk-ant-api01-..."

  # Analytics API key (Claude Enterprise). Env: ANTHROPIC_ANALYTICS_API_KEY
  # analytics_api_key = "..."
}
```

Environment variables are used as a fallback when a config argument is not set.
