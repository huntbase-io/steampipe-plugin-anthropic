---
title: "Steampipe Table: anthropic_claude_code_analytics - Query Claude Code Analytics using SQL"
description: "Allows users to query daily per-user Claude Code productivity metrics: sessions, lines of code, commits, PRs, tool acceptance, and cost."
---

# Table: anthropic_claude_code_analytics - Query Claude Code Analytics using SQL

The Claude Code Analytics API reports daily aggregated productivity metrics per user for Claude Code usage on the Claude API.

## Table Usage Guide

The `anthropic_claude_code_analytics` table returns one row per user per day. Requires an Admin API key. Defaults to the last 7 days; filter with the `date` column (`=`, `>=`, `<=`). Claude Code usage on claude.ai Enterprise seats is reported by the `anthropic_analytics_*` tables instead.

## Examples

### Developer productivity for the last 30 days
Track sessions, code changes, and commits per developer.

```sql+postgres
select
  date,
  actor_email,
  num_sessions,
  lines_added,
  lines_removed,
  commits,
  pull_requests
from
  anthropic_claude_code_analytics
where
  date >= now() - interval '30 days'
order by
  date desc;
```

```sql+sqlite
select
  date,
  actor_email,
  num_sessions,
  lines_added,
  lines_removed,
  commits,
  pull_requests
from
  anthropic_claude_code_analytics
where
  date >= datetime('now', '-30 days')
order by
  date desc;
```

### Edit tool acceptance rate per user
Measure how often proposals are accepted.

```sql+postgres
select
  actor_email,
  sum((tool_actions -> 'edit_tool' ->> 'accepted')::int) as accepted,
  sum((tool_actions -> 'edit_tool' ->> 'rejected')::int) as rejected
from
  anthropic_claude_code_analytics
where
  date >= now() - interval '30 days'
group by
  actor_email;
```

```sql+sqlite
select
  actor_email,
  sum(cast(tool_actions -> 'edit_tool' ->> 'accepted' as integer)) as accepted,
  sum(cast(tool_actions -> 'edit_tool' ->> 'rejected' as integer)) as rejected
from
  anthropic_claude_code_analytics
where
  date >= datetime('now', '-30 days')
group by
  actor_email;
```

### Estimated cost by model
Break down spend per model (amount is cents USD).

```sql+postgres
select
  m ->> 'model' as model,
  sum((m -> 'estimated_cost' ->> 'amount')::numeric) / 100 as usd
from
  anthropic_claude_code_analytics,
  jsonb_array_elements(model_breakdown) as m
where
  date >= now() - interval '30 days'
group by
  model;
```

```sql+sqlite
select
  json_extract(m.value, '$.model') as model,
  sum(json_extract(m.value, '$.estimated_cost.amount')) / 100 as usd
from
  anthropic_claude_code_analytics,
  json_each(model_breakdown) as m
where
  date >= datetime('now', '-30 days')
group by
  model;
```
