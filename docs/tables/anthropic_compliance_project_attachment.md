---
title: "Steampipe Table: anthropic_compliance_project_attachment - Query claude.ai Project Attachments using SQL"
description: "Allows users to query files and documents attached to claude.ai projects."
---

# Table: anthropic_compliance_project_attachment - Query claude.ai Project Attachments using SQL

Project attachments are binary uploads (`project_file`, IDs `claude_file_...`) or plain-text documents (`project_doc`, IDs `claude_proj_doc_...`).

## Table Usage Guide

The `anthropic_compliance_project_attachment` table iterates all projects automatically. Requires a Compliance Access Key with the `read:compliance_user_data` scope.

## Examples

### List attachments across all projects
Inventory project knowledge-base content.

```sql+postgres
select
  project_id,
  id,
  filename,
  mime_type,
  type
from
  anthropic_compliance_project_attachment;
```

```sql+sqlite
select
  project_id,
  id,
  filename,
  mime_type,
  type
from
  anthropic_compliance_project_attachment;
```

### Find binary uploads only
Separate binary files from text documents for download tooling.

```sql+postgres
select
  project_id,
  id,
  filename
from
  anthropic_compliance_project_attachment
where
  type = 'project_file';
```

```sql+sqlite
select
  project_id,
  id,
  filename
from
  anthropic_compliance_project_attachment
where
  type = 'project_file';
```
