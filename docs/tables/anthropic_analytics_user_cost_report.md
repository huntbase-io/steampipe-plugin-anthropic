---
title: "Steampipe Table: anthropic_analytics_user_cost_report - Query Claude Enterprise Per-User Costs using SQL"
description: "Allows users to query per-user spend from the Claude Enterprise Analytics API."
---

# Table: anthropic_analytics_user_cost_report - Query Claude Enterprise Per-User Costs using SQL

Ranks users by spend over a reporting window. Amounts are decimal strings in cents USD.

## Table Usage Guide

The `anthropic_analytics_user_cost_report` table returns one row per user for the window. Requires an Analytics API key. Defaults to the last 7 days; filter with the `date` column.

## Examples

### Top spenders in USD
Rank users by cost.

```sql+postgres
select
  user_email,
  sum(amount::numeric) / 100 as usd,
  sum(requests) as requests
from
  anthropic_analytics_user_cost_report
group by
  user_email
order by
  usd desc
limit 10;
```

```sql+sqlite
select
  user_email,
  sum(cast(amount as real)) / 100 as usd,
  sum(requests) as requests
from
  anthropic_analytics_user_cost_report
group by
  user_email
order by
  usd desc
limit 10;
```

### Discounted vs list price
Compare actual spend against list price.

```sql+postgres
select
  user_email,
  sum(amount::numeric) / 100 as actual_usd,
  sum(list_amount::numeric) / 100 as list_usd
from
  anthropic_analytics_user_cost_report
group by
  user_email;
```

```sql+sqlite
select
  user_email,
  sum(cast(amount as real)) / 100 as actual_usd,
  sum(cast(list_amount as real)) / 100 as list_usd
from
  anthropic_analytics_user_cost_report
group by
  user_email;
```

### Find spend outliers
Surface users spending more than three times the average.

```sql+postgres
with spend as (
  select
    user_email,
    sum(amount::numeric) / 100 as usd
  from
    anthropic_analytics_user_cost_report
  group by
    user_email
)
select
  user_email,
  usd
from
  spend
where
  usd > 3 * (select avg(usd) from spend)
order by
  usd desc;
```

```sql+sqlite
with spend as (
  select
    user_email,
    sum(cast(amount as real)) / 100 as usd
  from
    anthropic_analytics_user_cost_report
  group by
    user_email
)
select
  user_email,
  usd
from
  spend
where
  usd > 3 * (select avg(usd) from spend)
order by
  usd desc;
```
