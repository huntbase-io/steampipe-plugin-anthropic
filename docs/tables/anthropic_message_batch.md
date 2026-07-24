---
title: "Steampipe Table: anthropic_message_batch - Query Anthropic Message Batches using SQL"
description: "Allows users to query Message Batches API jobs, including processing status, request counts, and expiry."
---

# Table: anthropic_message_batch - Query Anthropic Message Batches using SQL

The Message Batches API processes large volumes of Messages requests asynchronously. This table lists batches in the workspace of the API key.

## Table Usage Guide

The `anthropic_message_batch` table provides insights into batch jobs: their processing status, per-result counts, and lifecycle timestamps. Requires a standard API key.

## Examples

### List batches that are still processing
Track in-flight batch jobs.

```sql+postgres
select
  id,
  processing_status,
  created_at,
  expires_at,
  request_counts
from
  anthropic_message_batch
where
  processing_status = 'in_progress';
```

```sql+sqlite
select
  id,
  processing_status,
  created_at,
  expires_at,
  request_counts
from
  anthropic_message_batch
where
  processing_status = 'in_progress';
```

### Summarize request outcomes for ended batches
Check error rates across completed batches.

```sql+postgres
select
  id,
  ended_at,
  request_counts ->> 'succeeded' as succeeded,
  request_counts ->> 'errored' as errored,
  request_counts ->> 'expired' as expired
from
  anthropic_message_batch
where
  processing_status = 'ended';
```

```sql+sqlite
select
  id,
  ended_at,
  request_counts ->> 'succeeded' as succeeded,
  request_counts ->> 'errored' as errored,
  request_counts ->> 'expired' as expired
from
  anthropic_message_batch
where
  processing_status = 'ended';
```
