# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Steampipe plugin (Go) that exposes the Anthropic API as Postgres tables. It covers three API surfaces, each with its own credential:

1. **Data plane** (standard API key, `sk-ant-api...`): models, message batches, files — implemented via the official `anthropic-sdk-go`.
2. **Admin API** (`sk-ant-admin...`): workspaces, workspace/org members, invites, API keys — the Go SDK does not cover `/v1/organizations/*`, so these use the hand-rolled REST client in `anthropic/admin_client.go`.
3. **Compliance API** (Compliance Access Key created in claude.ai, Claude Enterprise only): activity feed, linked organizations, users, roles, groups, settings, chats, chat messages, projects, attachments — all under `/v1/compliance/*`, also via the REST client. This is the primary area of interest for this project.
4. **Analytics** — two separate APIs with different keys: the **Claude Code Analytics API** (`/v1/organizations/usage_report/claude_code`, Admin key, one UTC day per request — the table iterates the date range) and the **Claude Enterprise Analytics API** (`/v1/organizations/analytics/*`, Analytics API key created in claude.ai). The `usage_report`/`cost_report` endpoints nest rows under `data[].results[]` time buckets — handled by `listAnalyticsBuckets`.

## Commands

- `go build ./...` — compile.
- `make sqlite` — build the SQLite loadable extension (`steampipe_sqlite_anthropic.so`). Clones `turbot/steampipe-sqlite` into `build/`, renders its templates for this plugin, and adds a `replace` directive to this repo (the module isn't published). **The plugin must stay on `steampipe-plugin-sdk/v6`** — steampipe-sqlite compiles the plugin into its own binary against v6; a v5 plugin won't type-check. Test with Homebrew sqlite3 (`/opt/homebrew/opt/sqlite/bin/sqlite3`, system sqlite3 blocks `.load`): `.load ./steampipe_sqlite_anthropic.so`, then `select steampipe_configure_anthropic('admin_api_key = "..."')`, then query tables by their normal names.
- `make install` — build and install to `~/.steampipe/plugins/local/anthropic/anthropic.plugin`. The local-plugin path convention matters: the `.spc` uses `plugin = "local/anthropic"`, which maps to that exact directory. Installing under `hub.steampipe.io/...` results in a "plugin not installed" connection error.
- `steampipe service restart` — required after reinstalling for the running service to pick up the new binary; then `steampipe query "..."` to test. Schema load is async — if `information_schema.tables` shows no `anthropic` tables right after restart, wait a few seconds. Connection errors surface in `steampipe_internal.steampipe_connection_state`.
- Connection config lives at `~/.steampipe/config/anthropic.spc` (template in `config/anthropic.spc`).
- There are no tests; verification is compiling plus live queries with real keys (env vars: `ANTHROPIC_API_KEY`, `ANTHROPIC_ADMIN_KEY`, `ANTHROPIC_COMPLIANCE_ACCESS_KEY`). Note the Steampipe service does not inherit the interactive shell's env — put keys in the `.spc` or the service's environment.

## Architecture

- `anthropic/plugin.go` registers every table; `anthropic/connection_config.go` resolves the three credentials (config value falls back to env var) and builds clients. `getActivityFeedClient` is special: the activity feed accepts either a Compliance Access Key or an Admin key, so it tries both.
- `anthropic/admin_client.go` is the shared REST client for Admin + Compliance endpoints. It implements **two pagination schemes**, and picking the right one per endpoint is the main correctness trap:
  - `listAll` — cursor scheme (`has_more`/`first_id`/`last_id`, pass `last_id` back as `after_id`): Admin endpoints, activity feed, compliance chats.
  - `listAllPaged` — token scheme (`has_more`/`next_page`, pass token back as `page`): compliance organizations, users, roles, groups, group members, projects, attachments.
  - Chat messages (`/v1/compliance/apps/chats/{id}/messages`) use a third shape: the envelope is the chat object with a `chat_messages` array plus cursor fields — handled inline in its table file.
- Tables are one file each (`table_anthropic_*.go`). Admin/Compliance rows are plain structs with `json` tags decoded from raw JSON; the plugin's `DefaultTransform: transform.FromCamel()` maps column `foo_bar` → struct field `FooBar`, so extra `Transform:` entries are only needed for nested fields (e.g. `Actor.EmailAddress`) or fields whose JSON name doesn't match (e.g. `Raw`).
- Parent/child tables use `ParentHydrate` (e.g. `anthropic_workspace_member` lists workspaces first; `anthropic_compliance_user` iterates organizations). `anthropic_compliance_chat_message` instead requires a `chat_id` qual (`KeyColumns: plugin.SingleColumn("chat_id")`) because listing messages for every chat would be prohibitively expensive.
- Every list hydrate honors `d.QueryContext.GetLimit()` / `d.RowsRemaining(ctx)` so `LIMIT` queries stop paginating early.
- Analytics tables take optional `date` quals (`=`, `>`, `>=`, `<`, `<=`) resolved by `getDateRangeFromQuals` in `anthropic/utils.go`; with no qual they default to a trailing 7- or 30-day window. `ending_date`/`ending_at` request params are exclusive, so the resolved inclusive end date is sent as `end + 1 day`.
- `listAllPaged`'s `has_more` is a `*bool` on purpose: some analytics endpoints omit it and signal the last page via a null/empty `next_page` alone. Don't "simplify" it back to `bool`.

## API references

- Compliance API docs: https://platform.claude.com/docs/en/manage-claude/compliance-api (sub-pages: compliance-activity-feed, compliance-org-data, compliance-content-data, compliance-errors). All compliance endpoints share one 600 req/min rate limit per parent org.
- Admin API docs: https://platform.claude.com/docs/en/manage-claude/admin-api
- Compliance endpoint shapes were transcribed from docs, not an SDK — when adding columns or new compliance tables, verify field names against the live docs rather than guessing.
