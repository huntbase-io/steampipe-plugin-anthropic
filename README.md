# Anthropic Plugin for Steampipe

Use SQL to query models, message batches, files, organization administration, and Claude Enterprise compliance data from the [Anthropic API](https://platform.claude.com/docs/en/api/).

## Quick start

```sh
make install                       # builds to ~/.steampipe/plugins/local/anthropic/anthropic.plugin
cp config/anthropic.spc ~/.steampipe/config/anthropic.spc
steampipe query "select id, display_name from anthropic_model"
```

## Credentials

The plugin uses three independent credentials, each unlocking a different set of tables. Configure any subset in `~/.steampipe/config/anthropic.spc` or via environment variables.

| Config argument | Env var | Key type | Tables |
|---|---|---|---|
| `api_key` | `ANTHROPIC_API_KEY` | Standard API key (`sk-ant-api...`) | `anthropic_model`, `anthropic_message_batch`, `anthropic_file` |
| `admin_api_key` | `ANTHROPIC_ADMIN_KEY` | Admin API key (`sk-ant-admin...`) | `anthropic_workspace`, `anthropic_workspace_member`, `anthropic_organization_member`, `anthropic_invite`, `anthropic_api_key`, and `anthropic_compliance_activity` (activity feed only) |
| `compliance_api_key` | `ANTHROPIC_COMPLIANCE_ACCESS_KEY` | Compliance Access Key created in claude.ai (Claude Enterprise only) | All `anthropic_compliance_*` tables |
| `analytics_api_key` | `ANTHROPIC_ANALYTICS_API_KEY` | Analytics API key created in claude.ai by the primary owner (Claude Enterprise only) | All `anthropic_analytics_*` tables |

`anthropic_claude_code_analytics` (Claude Code Analytics API) uses the **Admin API key**, not the Analytics key.

Compliance Access Key scopes: `read:compliance_activities` (activity feed), `read:compliance_org_data` (organizations, roles, groups, settings), `read:compliance_user_data` (users, group members, chats, chat messages, projects, attachments).

## Example queries

```sql
-- Which models can I use?
select id, display_name, created_at from anthropic_model;

-- Audit API keys across the org
select name, status, workspace_id, partial_key_hint, created_at
from anthropic_api_key
where status = 'active';

-- Recent compliance activity by user
select created_at, type, actor_email, organization_uuid
from anthropic_compliance_activity
limit 100;

-- Chats updated recently, with owner
select id, name, user_email, model, updated_at
from anthropic_compliance_chat
where updated_at > now() - interval '7 days';

-- Full message content of one chat
select role, created_at, content
from anthropic_compliance_chat_message
where chat_id = 'claude_chat_01H5CWunD7RpVJ5bHa8RCkja';

-- Claude Code productivity per developer, last 30 days
select date, actor_email, num_sessions, lines_added, lines_removed, commits, pull_requests
from anthropic_claude_code_analytics
where date >= now() - interval '30 days'
order by date desc;

-- Enterprise adoption trend
select date, daily_active_user_count, weekly_active_user_count, assigned_seat_count
from anthropic_analytics_summary
where date >= now() - interval '90 days'
order by date;

-- Top spenders this week (amount is cents USD as a decimal string)
select user_email, sum(amount::numeric) / 100 as usd
from anthropic_analytics_user_cost_report
group by user_email
order by usd desc
limit 10;

-- Attest effective retention/redaction settings per org
select o.name, s.name as setting, s.value
from anthropic_compliance_organization o
join anthropic_compliance_organization_setting s on s.organization_uuid = o.uuid
where s.name in ('data_retention_periods', 'content_redaction_enabled');
```

## SQLite extension

The plugin can also be built as a SQLite loadable extension:

```sh
make sqlite    # produces ./steampipe_sqlite_anthropic.so

# macOS system sqlite3 blocks .load — use Homebrew's:
/opt/homebrew/opt/sqlite/bin/sqlite3
sqlite> .load ./steampipe_sqlite_anthropic.so
sqlite> select steampipe_configure_anthropic('
   ...>   admin_api_key = "sk-ant-admin01-..."
   ...> ');
sqlite> select id, name from anthropic_workspace;
```

Configuration uses the same HCL arguments as the `.spc` file; environment variables (`ANTHROPIC_API_KEY`, etc.) also work without calling the configure function.

## Development

```sh
go build ./...      # compile
make install        # build + install into ~/.steampipe/plugins/local/anthropic/
steampipe service restart   # pick up a new binary
```
