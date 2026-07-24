---
title: "Steampipe Table: anthropic_analytics_user_activity - Query Claude Enterprise User Activity using SQL"
description: "Allows users to query per-user activity metrics across Claude products from the Claude Enterprise Analytics API."
---

# Table: anthropic_analytics_user_activity - Query Claude Enterprise User Activity using SQL

Per-user metrics across chat, Claude Code, Cowork, design, office agents, and science, with per-product metric blocks.

## Table Usage Guide

The `anthropic_analytics_user_activity` table returns one row per user for the queried window. Requires an Analytics API key (Claude Enterprise). Defaults to the last 7 days; filter with the `date` column. Per-product metric objects are always present — organizations without usage of a product see zeros, not nulls.

## Examples

### Most active chat users
Rank users by message volume.

```sql+postgres
select
  user_email,
  chat_metrics ->> 'message_count' as messages,
  chat_metrics ->> 'distinct_conversation_count' as conversations,
  last_activity_date
from
  anthropic_analytics_user_activity
order by
  (chat_metrics ->> 'message_count')::int desc
limit 20;
```

```sql+sqlite
select
  user_email,
  chat_metrics ->> 'message_count' as messages,
  chat_metrics ->> 'distinct_conversation_count' as conversations,
  last_activity_date
from
  anthropic_analytics_user_activity
order by
  cast(chat_metrics ->> 'message_count' as integer) desc
limit 20;
```

### Claude Code activity per user
Surface sessions and commits from the nested metrics block.

```sql+postgres
select
  user_email,
  claude_code_metrics -> 'core_metrics' ->> 'distinct_session_count' as sessions,
  claude_code_metrics -> 'core_metrics' ->> 'commit_count' as commits,
  claude_code_metrics -> 'core_metrics' ->> 'pull_request_count' as pull_requests
from
  anthropic_analytics_user_activity;
```

```sql+sqlite
select
  user_email,
  claude_code_metrics -> 'core_metrics' ->> 'distinct_session_count' as sessions,
  claude_code_metrics -> 'core_metrics' ->> 'commit_count' as commits,
  claude_code_metrics -> 'core_metrics' ->> 'pull_request_count' as pull_requests
from
  anthropic_analytics_user_activity;
```
