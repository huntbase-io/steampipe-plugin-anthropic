---
title: "Steampipe Table: anthropic_invite - Query Anthropic Organization Invites using SQL"
description: "Allows users to query pending and historical invitations to the Anthropic organization."
---

# Table: anthropic_invite - Query Anthropic Organization Invites using SQL

Invites grant new users membership of the Anthropic organization with a specified role; they expire after 21 days.

## Table Usage Guide

The `anthropic_invite` table provides insights into invitations and their status. Requires an Admin API key.

## Examples

### List pending invites
Track invitations that have not yet been accepted.

```sql+postgres
select
  email,
  role,
  invited_at,
  expires_at
from
  anthropic_invite
where
  status = 'pending';
```

```sql+sqlite
select
  email,
  role,
  invited_at,
  expires_at
from
  anthropic_invite
where
  status = 'pending';
```

### Find invites granting elevated roles
Review invitations that would grant admin or developer access.

```sql+postgres
select
  email,
  role,
  status,
  invited_at
from
  anthropic_invite
where
  role in ('admin', 'developer');
```

```sql+sqlite
select
  email,
  role,
  status,
  invited_at
from
  anthropic_invite
where
  role in ('admin', 'developer');
```
