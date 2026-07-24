---
title: "Steampipe Table: anthropic_analytics_cost_report - Query Claude Enterprise Costs using SQL"
description: "Allows users to query daily cost buckets by model, product, and cost type from the Claude Enterprise Analytics API."
---

# Table: anthropic_analytics_cost_report - Query Claude Enterprise Costs using SQL

Daily cost broken down by model, product, and cost type. Amounts are decimal strings in cents USD.

## Table Usage Guide

The `anthropic_analytics_cost_report` table returns one row per day per model/product/cost-type combination. Requires an Analytics API key. Defaults to the last 7 days (31-day range cap). Cost data can be revised for up to 30 days; for invoicing-grade totals query dates at least 30 days in the past.

## Examples

### Daily spend in USD
Convert cents-denominated amounts to dollars.

```sql+postgres
select
  date,
  sum(amount::numeric) / 100 as usd
from
  anthropic_analytics_cost_report
group by
  date
order by
  date;
```

```sql+sqlite
select
  date,
  sum(cast(amount as real)) / 100 as usd
from
  anthropic_analytics_cost_report
group by
  date
order by
  date;
```

### Spend by product and model
Attribute cost across products.

```sql+postgres
select
  product,
  model,
  sum(amount::numeric) / 100 as usd
from
  anthropic_analytics_cost_report
where
  date >= now() - interval '30 days'
group by
  product, model
order by
  usd desc;
```

```sql+sqlite
select
  product,
  model,
  sum(cast(amount as real)) / 100 as usd
from
  anthropic_analytics_cost_report
where
  date >= datetime('now', '-30 days')
group by
  product, model
order by
  usd desc;
```
