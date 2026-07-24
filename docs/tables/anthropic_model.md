---
title: "Steampipe Table: anthropic_model - Query Anthropic Models using SQL"
description: "Allows users to query Claude models available to the API key, including model IDs, display names, and release dates."
---

# Table: anthropic_model - Query Anthropic Models using SQL

The Models API lists the Claude models available to your API key, with their identifiers and release dates.

## Table Usage Guide

The `anthropic_model` table provides insights into the Claude models your organization can use. Requires a standard API key (`api_key` or `ANTHROPIC_API_KEY`).

## Examples

### List all available models
Discover the models available to your API key.

```sql+postgres
select
  id,
  display_name,
  created_at
from
  anthropic_model;
```

```sql+sqlite
select
  id,
  display_name,
  created_at
from
  anthropic_model;
```

### Find the newest models
Identify recently released models to plan upgrades.

```sql+postgres
select
  id,
  display_name,
  created_at
from
  anthropic_model
order by
  created_at desc
limit 5;
```

```sql+sqlite
select
  id,
  display_name,
  created_at
from
  anthropic_model
order by
  created_at desc
limit 5;
```
