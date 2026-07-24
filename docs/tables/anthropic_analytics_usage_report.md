---
title: "Steampipe Table: anthropic_analytics_usage_report - Query Claude Enterprise Token Usage using SQL"
description: "Allows users to query daily token usage buckets by model and product from the Claude Enterprise Analytics API."
---

# Table: anthropic_analytics_usage_report - Query Claude Enterprise Token Usage using SQL

Daily token usage broken down by model and product for a Claude Enterprise organization.

## Table Usage Guide

The `anthropic_analytics_usage_report` table returns one row per day per model/product combination. Requires an Analytics API key. Defaults to the last 7 days; the API caps a single range at 31 days.

## Examples

### Daily token usage by model
Track consumption trends per model.

```sql+postgres
select
  date,
  model,
  sum(uncached_input_tokens) as input_tokens,
  sum(output_tokens) as output_tokens,
  sum(cache_read_input_tokens) as cache_read
from
  anthropic_analytics_usage_report
group by
  date, model
order by
  date;
```

```sql+sqlite
select
  date,
  model,
  sum(uncached_input_tokens) as input_tokens,
  sum(output_tokens) as output_tokens,
  sum(cache_read_input_tokens) as cache_read
from
  anthropic_analytics_usage_report
group by
  date, model
order by
  date;
```

### Usage by product
Compare chat vs Claude Code vs other products.

```sql+postgres
select
  product,
  sum(requests) as requests,
  sum(output_tokens) as output_tokens
from
  anthropic_analytics_usage_report
where
  date >= now() - interval '30 days'
group by
  product;
```

```sql+sqlite
select
  product,
  sum(requests) as requests,
  sum(output_tokens) as output_tokens
from
  anthropic_analytics_usage_report
where
  date >= datetime('now', '-30 days')
group by
  product;
```
