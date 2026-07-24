---
title: "Steampipe Table: anthropic_compliance_role - Query Anthropic RBAC Roles using SQL"
description: "Allows users to query RBAC roles defined on organizations linked to the Claude Enterprise parent."
---

# Table: anthropic_compliance_role - Query Anthropic RBAC Roles using SQL

RBAC roles are custom roles defined per linked organization, independent of built-in membership levels.

## Table Usage Guide

The `anthropic_compliance_role` table iterates all linked organizations automatically. Requires a Compliance Access Key with the `read:compliance_org_data` scope.

## Examples

### List roles across organizations
Review the custom roles defined on each linked organization.

```sql+postgres
select
  organization_uuid,
  id,
  name,
  description,
  created_at
from
  anthropic_compliance_role;
```

```sql+sqlite
select
  organization_uuid,
  id,
  name,
  description,
  created_at
from
  anthropic_compliance_role;
```
