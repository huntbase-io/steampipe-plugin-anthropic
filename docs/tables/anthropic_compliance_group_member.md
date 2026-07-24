---
title: "Steampipe Table: anthropic_compliance_group_member - Query Anthropic Group Members using SQL"
description: "Allows users to query membership of each group in the Claude Enterprise parent organization."
---

# Table: anthropic_compliance_group_member - Query Anthropic Group Members using SQL

Lists the members of every group, keyed by user ID and email.

## Table Usage Guide

The `anthropic_compliance_group_member` table iterates all groups automatically. Requires a Compliance Access Key with the `read:compliance_org_data` and `read:compliance_user_data` scopes.

## Examples

### List members per group
Walk group membership end to end.

```sql+postgres
select
  g.name as group_name,
  gm.email,
  gm.user_id
from
  anthropic_compliance_group_member gm
  join anthropic_compliance_group g on g.id = gm.group_id;
```

```sql+sqlite
select
  g.name as group_name,
  gm.email,
  gm.user_id
from
  anthropic_compliance_group_member gm
  join anthropic_compliance_group g on g.id = gm.group_id;
```
