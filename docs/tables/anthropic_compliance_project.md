---
title: "Steampipe Table: anthropic_compliance_project - Query claude.ai Projects using SQL"
description: "Allows users to query claude.ai projects across the Claude Enterprise organization."
---

# Table: anthropic_compliance_project - Query claude.ai Projects using SQL

Projects bundle related chats with custom instructions and knowledge-base content. This table lists project metadata organization-wide.

## Table Usage Guide

The `anthropic_compliance_project` table provides project metadata; the full raw record is in the `project` JSON column. Requires a Compliance Access Key with the `read:compliance_user_data` scope.

## Examples

### List projects
Enumerate projects across the organization.

```sql+postgres
select
  id,
  name,
  organization_uuid,
  created_at
from
  anthropic_compliance_project;
```

```sql+sqlite
select
  id,
  name,
  organization_uuid,
  created_at
from
  anthropic_compliance_project;
```

### Find chats belonging to a project
Join projects with chats.

```sql+postgres
select
  p.name as project,
  c.name as chat,
  c.user_email
from
  anthropic_compliance_project p
  join anthropic_compliance_chat c on c.project_id = p.id;
```

```sql+sqlite
select
  p.name as project,
  c.name as chat,
  c.user_email
from
  anthropic_compliance_project p
  join anthropic_compliance_chat c on c.project_id = p.id;
```
