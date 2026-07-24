---
title: "Steampipe Table: anthropic_file - Query Anthropic Files using SQL"
description: "Allows users to query files uploaded to the Anthropic Files API, including filenames, MIME types, and sizes."
---

# Table: anthropic_file - Query Anthropic Files using SQL

The Files API (beta) stores files for reuse across Messages API requests. This table lists uploaded files.

## Table Usage Guide

The `anthropic_file` table provides insights into uploaded files: names, MIME types, sizes, and whether they are downloadable. Requires a standard API key.

## Examples

### List all files with size
Review stored files and their storage footprint.

```sql+postgres
select
  id,
  filename,
  mime_type,
  size_bytes,
  created_at
from
  anthropic_file
order by
  size_bytes desc;
```

```sql+sqlite
select
  id,
  filename,
  mime_type,
  size_bytes,
  created_at
from
  anthropic_file
order by
  size_bytes desc;
```

### Find large PDF uploads
Identify sizeable documents that may be candidates for cleanup.

```sql+postgres
select
  id,
  filename,
  size_bytes
from
  anthropic_file
where
  mime_type = 'application/pdf'
  and size_bytes > 10 * 1024 * 1024;
```

```sql+sqlite
select
  id,
  filename,
  size_bytes
from
  anthropic_file
where
  mime_type = 'application/pdf'
  and size_bytes > 10 * 1024 * 1024;
```
