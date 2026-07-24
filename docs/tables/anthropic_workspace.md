---
title: "Steampipe Table: anthropic_workspace - Query Anthropic Workspaces using SQL"
description: "Allows users to query workspaces in the Anthropic organization, including names, colors, and archival state."
---

# Table: anthropic_workspace - Query Anthropic Workspaces using SQL

Workspaces segment API keys and spend within an Anthropic organization. This table lists all workspaces.

## Table Usage Guide

The `anthropic_workspace` table provides insights into organization workspaces. Requires an Admin API key (`admin_api_key` or `ANTHROPIC_ADMIN_KEY`).

## Examples

### List all workspaces
Get an overview of workspaces in the organization.

```sql+postgres
select
  id,
  name,
  display_color,
  created_at,
  archived_at
from
  anthropic_workspace;
```

```sql+sqlite
select
  id,
  name,
  display_color,
  created_at,
  archived_at
from
  anthropic_workspace;
```

### Find archived workspaces
Identify workspaces that are no longer active.

```sql+postgres
select
  id,
  name,
  archived_at
from
  anthropic_workspace
where
  archived_at is not null;
```

```sql+sqlite
select
  id,
  name,
  archived_at
from
  anthropic_workspace
where
  archived_at is not null;
```
