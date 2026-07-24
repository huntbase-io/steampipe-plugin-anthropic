---
title: "Steampipe Table: anthropic_compliance_organization - Query Anthropic Linked Organizations using SQL"
description: "Allows users to query organizations linked to the Claude Enterprise parent organization."
---

# Table: anthropic_compliance_organization - Query Anthropic Linked Organizations using SQL

A Claude Enterprise tenant has one parent organization with linked claude.ai and Claude Console organizations. This table lists them.

## Table Usage Guide

The `anthropic_compliance_organization` table enumerates linked organizations. Requires a Compliance Access Key with the `read:compliance_org_data` scope. The `uuid` column joins to `organization_uuid` in activity, chat, and project tables.

## Examples

### List linked organizations
Enumerate every organization under the parent.

```sql+postgres
select
  uuid,
  name,
  created_at
from
  anthropic_compliance_organization;
```

```sql+sqlite
select
  uuid,
  name,
  created_at
from
  anthropic_compliance_organization;
```
