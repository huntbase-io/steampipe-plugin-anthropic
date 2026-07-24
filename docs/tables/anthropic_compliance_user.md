---
title: "Steampipe Table: anthropic_compliance_user - Query Anthropic Compliance Users using SQL"
description: "Allows users to query users of each organization linked to the Claude Enterprise parent."
---

# Table: anthropic_compliance_user - Query Anthropic Compliance Users using SQL

Lists active users of every linked organization, with organization-level roles and join dates.

## Table Usage Guide

The `anthropic_compliance_user` table iterates all linked organizations automatically. Requires a Compliance Access Key with the `read:compliance_org_data` and `read:compliance_user_data` scopes.

## Examples

### List users across all linked organizations
Seed an eDiscovery custodian list.

```sql+postgres
select
  id,
  full_name,
  email,
  organization_role,
  organization_uuid
from
  anthropic_compliance_user;
```

```sql+sqlite
select
  id,
  full_name,
  email,
  organization_role,
  organization_uuid
from
  anthropic_compliance_user;
```

### Find owners and admins
Audit elevated roles per linked organization.

```sql+postgres
select
  o.name as organization,
  u.email,
  u.organization_role
from
  anthropic_compliance_user u
  join anthropic_compliance_organization o on o.uuid = u.organization_uuid
where
  u.organization_role in ('owner', 'primary_owner', 'admin');
```

```sql+sqlite
select
  o.name as organization,
  u.email,
  u.organization_role
from
  anthropic_compliance_user u
  join anthropic_compliance_organization o on o.uuid = u.organization_uuid
where
  u.organization_role in ('owner', 'primary_owner', 'admin');
```
