<p align="center">
  <img src="docs/images/anthropic.svg" width="100" alt="Anthropic logo" />
</p>

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

## Security & IT audit queries

```sql
-- 1. Stale active API keys (rotate anything older than 90 days)
select name, partial_key_hint, created_at,
  extract(day from now() - created_at) as age_days
from anthropic_api_key
where status = 'active'
  and created_at < now() - interval '90 days'
order by created_at;

-- 2. API keys not scoped to a workspace (live in the default workspace)
select name, status, partial_key_hint, created_at
from anthropic_api_key
where workspace_id is null
  and status = 'active';

-- 3. Who holds elevated org roles right now
select email, role, added_at
from anthropic_organization_member
where role in ('admin', 'developer')
order by role, added_at;

-- 4. Pending invites that grant elevated access (approve or revoke before they land)
select email, role, invited_at, expires_at
from anthropic_invite
where status = 'pending'
  and role in ('admin', 'developer');

-- 5. Workspace access held by users no longer in the organization (orphaned grants)
select wm.workspace_id, wm.user_id, wm.workspace_role
from anthropic_workspace_member wm
  left join anthropic_organization_member m on m.id = wm.user_id
where m.id is null;

-- 6. Admin-surface changes in the audit feed (key creation, role changes, SSO edits)
select created_at, type, actor_email, actor_ip_address
from anthropic_compliance_activity
where type like 'admin_api_key%'
   or type like 'rbac_%'
   or type like 'org_sso%'
order by created_at desc;

-- 7. Users signing in from many distinct IPs (session sharing / credential theft signal)
select actor_email, count(distinct actor_ip_address) as distinct_ips,
  array_agg(distinct actor_ip_address) as ips
from anthropic_compliance_activity
where actor_email is not null
group by actor_email
having count(distinct actor_ip_address) > 3
order by distinct_ips desc;

-- 8. Activity outside business hours (UTC) by user
select actor_email, count(*) as events
from anthropic_compliance_activity
where extract(hour from created_at) not between 7 and 19
group by actor_email
order by events desc;

-- 9. Chats users soft-deleted themselves (visible to compliance until retention expiry)
select id, name, user_email, created_at, deleted_at
from anthropic_compliance_chat
where deleted_at is not null
order by deleted_at desc;

-- 10. Attest security controls per linked organization (redaction, IP allowlist, SSO)
select o.name as organization, s.name as setting, s.value
from anthropic_compliance_organization_setting s
  join anthropic_compliance_organization o on o.uuid = s.organization_uuid
where s.name in ('content_redaction_enabled', 'ip_allowlist_enabled',
  'ip_allowlist_ip_ranges', 'sso_provisioning_mode', 'data_retention_periods');

-- 11. Groups managed by hand rather than SCIM (membership drift risk)
select g.name, g.source_type, count(gm.user_id) as members
from anthropic_compliance_group g
  left join anthropic_compliance_group_member gm on gm.group_id = g.id
where g.source_type = 'direct'
group by g.name, g.source_type;

-- 12. Deleted users still showing token usage (offboarding gap)
select user_email, user_name, total_tokens, requests
from anthropic_analytics_user_usage_report
where user_deleted;

-- 13. Spend outliers: users more than 3x the average
with spend as (
  select user_email, sum(amount::numeric) / 100 as usd
  from anthropic_analytics_user_cost_report
  group by user_email
)
select user_email, usd
from spend
where usd > 3 * (select avg(usd) from spend)
order by usd desc;

-- 14. Project knowledge-base uploads (DLP review surface)
select p.name as project, a.filename, a.mime_type, a.created_at
from anthropic_compliance_project_attachment a
  join anthropic_compliance_project p on p.id = a.project_id
order by a.created_at desc;
```

## SQLite extension

The plugin can also be built as a SQLite loadable extension:

```sh
make sqlite      # produces ./steampipe_sqlite_anthropic.so
make postgres    # produces build/postgres/ (steampipe_postgres_anthropic FDW; needs pg_config)

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

---

The Anthropic name and logo are trademarks of Anthropic, PBC, used here for identification only. This is a community plugin and is not affiliated with or endorsed by Anthropic.
