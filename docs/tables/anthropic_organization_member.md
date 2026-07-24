---
title: "Steampipe Table: anthropic_organization_member - Query Anthropic Organization Members using SQL"
description: "Allows users to query users in the Anthropic organization, including emails, roles, and join dates."
---

# Table: anthropic_organization_member - Query Anthropic Organization Members using SQL

Organization members are the users of the Anthropic Console organization, each with an organization-level role.

## Table Usage Guide

The `anthropic_organization_member` table provides insights into organization users and their roles. Requires an Admin API key.

## Examples

### List all members with roles
Review who belongs to the organization.

```sql+postgres
select
  id,
  email,
  name,
  role,
  added_at
from
  anthropic_organization_member;
```

```sql+sqlite
select
  id,
  email,
  name,
  role,
  added_at
from
  anthropic_organization_member;
```

### Count members by role
Audit the distribution of privileges.

```sql+postgres
select
  role,
  count(*)
from
  anthropic_organization_member
group by
  role
order by
  count desc;
```

```sql+sqlite
select
  role,
  count(*)
from
  anthropic_organization_member
group by
  role;
```
