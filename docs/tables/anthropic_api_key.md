---
title: "Steampipe Table: anthropic_api_key - Query Anthropic API Keys using SQL"
description: "Allows users to query API keys across the Anthropic organization, including status, workspace, and creator."
---

# Table: anthropic_api_key - Query Anthropic API Keys using SQL

API keys authenticate requests to the Anthropic API. Each key belongs to a workspace (or the default workspace) and has a lifecycle status.

## Table Usage Guide

The `anthropic_api_key` table provides insights into API keys for security audits and hygiene. Requires an Admin API key.

## Examples

### List active API keys
Inventory keys currently in use.

```sql+postgres
select
  id,
  name,
  workspace_id,
  partial_key_hint,
  created_at
from
  anthropic_api_key
where
  status = 'active';
```

```sql+sqlite
select
  id,
  name,
  workspace_id,
  partial_key_hint,
  created_at
from
  anthropic_api_key
where
  status = 'active';
```

### Find old active keys
Identify keys that may be due for rotation.

```sql+postgres
select
  name,
  created_at,
  partial_key_hint
from
  anthropic_api_key
where
  status = 'active'
  and created_at < now() - interval '180 days';
```

```sql+sqlite
select
  name,
  created_at,
  partial_key_hint
from
  anthropic_api_key
where
  status = 'active'
  and created_at < datetime('now', '-180 days');
```

### Show who created each key
Attribute keys to their creators.

```sql+postgres
select
  name,
  status,
  created_by ->> 'id' as created_by_id
from
  anthropic_api_key;
```

```sql+sqlite
select
  name,
  status,
  created_by ->> 'id' as created_by_id
from
  anthropic_api_key;
```

### Find active keys in the default workspace
Keys not scoped to a workspace bypass workspace-level spend controls.

```sql+postgres
select
  name,
  status,
  partial_key_hint,
  created_at
from
  anthropic_api_key
where
  workspace_id is null
  and status = 'active';
```

```sql+sqlite
select
  name,
  status,
  partial_key_hint,
  created_at
from
  anthropic_api_key
where
  workspace_id is null
  and status = 'active';
```
