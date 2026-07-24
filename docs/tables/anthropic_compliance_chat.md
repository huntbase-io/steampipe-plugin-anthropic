---
title: "Steampipe Table: anthropic_compliance_chat - Query claude.ai Chats using SQL"
description: "Allows users to query chat metadata across the Claude Enterprise organization for eDiscovery and DLP."
---

# Table: anthropic_compliance_chat - Query claude.ai Chats using SQL

Lists claude.ai chat metadata (name, owner, model, timestamps) organization-wide. Soft-deleted chats appear with `deleted_at` set; hard-deleted chats are not retrievable.

## Table Usage Guide

The `anthropic_compliance_chat` table provides chat metadata only — use `anthropic_compliance_chat_message` for content. Requires a Compliance Access Key with the `read:compliance_user_data` scope.

## Examples

### List recently updated chats
Keep an export current.

```sql+postgres
select
  id,
  name,
  user_email,
  model,
  updated_at
from
  anthropic_compliance_chat
where
  updated_at > now() - interval '7 days';
```

```sql+sqlite
select
  id,
  name,
  user_email,
  model,
  updated_at
from
  anthropic_compliance_chat
where
  updated_at > datetime('now', '-7 days');
```

### Find chats owned by a specific user
Scope a legal hold to a custodian.

```sql+postgres
select
  id,
  name,
  created_at,
  deleted_at
from
  anthropic_compliance_chat
where
  user_email = 'user@example.com';
```

```sql+sqlite
select
  id,
  name,
  created_at,
  deleted_at
from
  anthropic_compliance_chat
where
  user_email = 'user@example.com';
```

### Count chats per user
Understand usage distribution.

```sql+postgres
select
  user_email,
  count(*)
from
  anthropic_compliance_chat
group by
  user_email
order by
  count desc;
```

```sql+sqlite
select
  user_email,
  count(*)
from
  anthropic_compliance_chat
group by
  user_email;
```

### List soft-deleted chats
Chats users deleted in claude.ai remain visible to compliance until the retention window expires — useful in exfiltration investigations.

```sql+postgres
select
  id,
  name,
  user_email,
  created_at,
  deleted_at
from
  anthropic_compliance_chat
where
  deleted_at is not null
order by
  deleted_at desc;
```

```sql+sqlite
select
  id,
  name,
  user_email,
  created_at,
  deleted_at
from
  anthropic_compliance_chat
where
  deleted_at is not null
order by
  deleted_at desc;
```
