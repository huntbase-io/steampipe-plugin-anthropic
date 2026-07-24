---
title: "Steampipe Table: anthropic_analytics_user_usage_report - Query Claude Enterprise Per-User Token Usage using SQL"
description: "Allows users to query per-user token consumption from the Claude Enterprise Analytics API."
---

# Table: anthropic_analytics_user_usage_report - Query Claude Enterprise Per-User Token Usage using SQL

Ranks users by token consumption over a reporting window.

## Table Usage Guide

The `anthropic_analytics_user_usage_report` table returns one row per user for the window. Requires an Analytics API key. Defaults to the last 7 days; filter with the `date` column.

## Examples

### Top token consumers
Identify the heaviest users.

```sql+postgres
select
  user_email,
  total_tokens,
  requests,
  output_tokens
from
  anthropic_analytics_user_usage_report
order by
  total_tokens desc
limit 20;
```

```sql+sqlite
select
  user_email,
  total_tokens,
  requests,
  output_tokens
from
  anthropic_analytics_user_usage_report
order by
  total_tokens desc
limit 20;
```

### Exclude deleted users from a report
Filter out departed accounts.

```sql+postgres
select
  user_email,
  total_tokens
from
  anthropic_analytics_user_usage_report
where
  not user_deleted;
```

```sql+sqlite
select
  user_email,
  total_tokens
from
  anthropic_analytics_user_usage_report
where
  user_deleted = 0;
```
