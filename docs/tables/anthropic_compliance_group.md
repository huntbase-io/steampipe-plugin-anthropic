---
title: "Steampipe Table: anthropic_compliance_group - Query Anthropic Groups using SQL"
description: "Allows users to query RBAC and SCIM-provisioned groups in the Claude Enterprise parent organization."
---

# Table: anthropic_compliance_group - Query Anthropic Groups using SQL

Groups are collections of users, created manually in claude.ai (`direct`) or synced from an identity provider (`scim`), optionally carrying role assignments.

## Table Usage Guide

The `anthropic_compliance_group` table lists groups with their source and assigned role IDs. Requires a Compliance Access Key with the `read:compliance_org_data` scope.

## Examples

### List groups with their provenance
Reconcile groups against your identity provider.

```sql+postgres
select
  id,
  name,
  source_type,
  roles,
  updated_at
from
  anthropic_compliance_group;
```

```sql+sqlite
select
  id,
  name,
  source_type,
  roles,
  updated_at
from
  anthropic_compliance_group;
```

### Find manually created groups
Identify groups not managed by SCIM.

```sql+postgres
select
  id,
  name,
  created_at
from
  anthropic_compliance_group
where
  source_type = 'direct';
```

```sql+sqlite
select
  id,
  name,
  created_at
from
  anthropic_compliance_group
where
  source_type = 'direct';
```
