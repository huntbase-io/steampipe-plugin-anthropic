---
title: "Steampipe Table: anthropic_analytics_summary - Query Claude Enterprise Activity Summaries using SQL"
description: "Allows users to query organization-wide daily active users, seats, and adoption rates from the Claude Enterprise Analytics API."
---

# Table: anthropic_analytics_summary - Query Claude Enterprise Activity Summaries using SQL

Daily organization-level snapshots of active users, seat counts, pending invites, and adoption rates across Claude products.

## Table Usage Guide

The `anthropic_analytics_summary` table returns one row per day. Requires an Analytics API key (Claude Enterprise). Defaults to the last 30 days; filter with the `date` column. Per-product breakdowns beyond chat and Claude Code are in the `summary` JSON column.

## Examples

### Adoption trend over the last quarter
Track active users against assigned seats.

```sql+postgres
select
  date,
  daily_active_user_count,
  weekly_active_user_count,
  assigned_seat_count,
  weekly_adoption_rate
from
  anthropic_analytics_summary
where
  date >= now() - interval '90 days'
order by
  date;
```

```sql+sqlite
select
  date,
  daily_active_user_count,
  weekly_active_user_count,
  assigned_seat_count,
  weekly_adoption_rate
from
  anthropic_analytics_summary
where
  date >= datetime('now', '-90 days')
order by
  date;
```

### Claude Code adoption
Compare Claude Code active users to overall activity.

```sql+postgres
select
  date,
  claude_code_daily_active_user_count,
  daily_active_user_count
from
  anthropic_analytics_summary
order by
  date desc;
```

```sql+sqlite
select
  date,
  claude_code_daily_active_user_count,
  daily_active_user_count
from
  anthropic_analytics_summary
order by
  date desc;
```
