---
title: "Steampipe Table: anthropic_compliance_chat_message - Query claude.ai Chat Messages using SQL"
description: "Allows users to query the full message content of a claude.ai chat, including files and artifacts."
---

# Table: anthropic_compliance_chat_message - Query claude.ai Chat Messages using SQL

Returns the messages of one chat with content blocks, uploaded files, generated files, and artifacts.

## Table Usage Guide

The `anthropic_compliance_chat_message` table **requires a `chat_id` qualifier** in the WHERE clause. Requires a Compliance Access Key with the `read:compliance_user_data` scope.

## Examples

### Get all messages of a chat
Export a chat transcript.

```sql+postgres
select
  role,
  created_at,
  content
from
  anthropic_compliance_chat_message
where
  chat_id = 'claude_chat_01H5CWunD7RpVJ5bHa8RCkja'
order by
  created_at;
```

```sql+sqlite
select
  role,
  created_at,
  content
from
  anthropic_compliance_chat_message
where
  chat_id = 'claude_chat_01H5CWunD7RpVJ5bHa8RCkja'
order by
  created_at;
```

### List file attachments in a chat
Find uploads to feed DLP review.

```sql+postgres
select
  id,
  role,
  files
from
  anthropic_compliance_chat_message
where
  chat_id = 'claude_chat_01H5CWunD7RpVJ5bHa8RCkja'
  and files is not null;
```

```sql+sqlite
select
  id,
  role,
  files
from
  anthropic_compliance_chat_message
where
  chat_id = 'claude_chat_01H5CWunD7RpVJ5bHa8RCkja'
  and files is not null;
```
