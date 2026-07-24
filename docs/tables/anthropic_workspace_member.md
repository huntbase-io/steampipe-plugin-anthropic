---
title: "Steampipe Table: anthropic_workspace_member - Query Anthropic Workspace Members using SQL"
description: "Allows users to query workspace membership and roles across the Anthropic organization."
---

# Table: anthropic_workspace_member - Query Anthropic Workspace Members using SQL

Workspace members map organization users into workspaces with a workspace-level role.

## Table Usage Guide

The `anthropic_workspace_member` table lists members of every workspace (it iterates workspaces automatically). Requires an Admin API key.

## Examples

### List members of every workspace
Audit who has access to which workspace.

```sql+postgres
select
  workspace_id,
  user_id,
  workspace_role
from
  anthropic_workspace_member;
```

```sql+sqlite
select
  workspace_id,
  user_id,
  workspace_role
from
  anthropic_workspace_member;
```

### Join with organization members to resolve emails
Show workspace access by user email.

```sql+postgres
select
  w.name as workspace,
  m.email,
  wm.workspace_role
from
  anthropic_workspace_member wm
  join anthropic_workspace w on w.id = wm.workspace_id
  join anthropic_organization_member m on m.id = wm.user_id;
```

```sql+sqlite
select
  w.name as workspace,
  m.email,
  wm.workspace_role
from
  anthropic_workspace_member wm
  join anthropic_workspace w on w.id = wm.workspace_id
  join anthropic_organization_member m on m.id = wm.user_id;
```

### Find workspace admins
Identify elevated access within workspaces.

```sql+postgres
select
  workspace_id,
  user_id
from
  anthropic_workspace_member
where
  workspace_role = 'workspace_admin';
```

```sql+sqlite
select
  workspace_id,
  user_id
from
  anthropic_workspace_member
where
  workspace_role = 'workspace_admin';
```

### Find orphaned workspace grants
Workspace access held by users who are no longer organization members.

```sql+postgres
select
  wm.workspace_id,
  wm.user_id,
  wm.workspace_role
from
  anthropic_workspace_member wm
  left join anthropic_organization_member m on m.id = wm.user_id
where
  m.id is null;
```

```sql+sqlite
select
  wm.workspace_id,
  wm.user_id,
  wm.workspace_role
from
  anthropic_workspace_member wm
  left join anthropic_organization_member m on m.id = wm.user_id
where
  m.id is null;
```
