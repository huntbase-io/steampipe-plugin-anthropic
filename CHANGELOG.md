# Changelog

## v0.1.0 [unreleased]

_What's new?_

- New tables:
  - Data plane (standard API key): [anthropic_model](docs/tables/anthropic_model.md), [anthropic_message_batch](docs/tables/anthropic_message_batch.md), [anthropic_file](docs/tables/anthropic_file.md)
  - Admin API: [anthropic_workspace](docs/tables/anthropic_workspace.md), [anthropic_workspace_member](docs/tables/anthropic_workspace_member.md), [anthropic_organization_member](docs/tables/anthropic_organization_member.md), [anthropic_invite](docs/tables/anthropic_invite.md), [anthropic_api_key](docs/tables/anthropic_api_key.md), [anthropic_claude_code_analytics](docs/tables/anthropic_claude_code_analytics.md)
  - Compliance API (Claude Enterprise): [anthropic_compliance_activity](docs/tables/anthropic_compliance_activity.md), [anthropic_compliance_organization](docs/tables/anthropic_compliance_organization.md), [anthropic_compliance_organization_setting](docs/tables/anthropic_compliance_organization_setting.md), [anthropic_compliance_user](docs/tables/anthropic_compliance_user.md), [anthropic_compliance_role](docs/tables/anthropic_compliance_role.md), [anthropic_compliance_group](docs/tables/anthropic_compliance_group.md), [anthropic_compliance_group_member](docs/tables/anthropic_compliance_group_member.md), [anthropic_compliance_chat](docs/tables/anthropic_compliance_chat.md), [anthropic_compliance_chat_message](docs/tables/anthropic_compliance_chat_message.md), [anthropic_compliance_project](docs/tables/anthropic_compliance_project.md), [anthropic_compliance_project_attachment](docs/tables/anthropic_compliance_project_attachment.md)
  - Enterprise Analytics API (Claude Enterprise): [anthropic_analytics_summary](docs/tables/anthropic_analytics_summary.md), [anthropic_analytics_user_activity](docs/tables/anthropic_analytics_user_activity.md), [anthropic_analytics_usage_report](docs/tables/anthropic_analytics_usage_report.md), [anthropic_analytics_cost_report](docs/tables/anthropic_analytics_cost_report.md), [anthropic_analytics_user_usage_report](docs/tables/anthropic_analytics_user_usage_report.md), [anthropic_analytics_user_cost_report](docs/tables/anthropic_analytics_user_cost_report.md)
- Four independent credentials (`api_key`, `admin_api_key`, `compliance_api_key`, `analytics_api_key`), each with environment-variable fallback.
- Automatic backoff and retry on rate-limit (429) and server (5xx/529) errors, honoring `retry-after`.
- Builds as a Postgres plugin (`make install`) and as a SQLite loadable extension (`make sqlite`).
