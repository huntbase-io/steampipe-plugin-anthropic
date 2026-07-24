---
title: "Steampipe Table: anthropic_compliance_organization_setting - Query Anthropic Effective Organization Settings using SQL"
description: "Allows users to query the effective data-privacy, security, and capability settings in force for each linked organization."
---

# Table: anthropic_compliance_organization_setting - Query Anthropic Effective Organization Settings using SQL

Returns the enforced settings for each linked organization (retention windows, content redaction, SSO enforcement, IP allowlist, session controls) — the resolved state after policy and dependency rules apply.

## Table Usage Guide

The `anthropic_compliance_organization_setting` table streams one row per setting per organization. Requires a Compliance Access Key with the `read:compliance_org_data` scope; the settings endpoint is enabled per parent organization. A missing setting row means "not controllable by this organization's administrators", not "off".

## Examples

### Attest retention and redaction settings
Verify controls match your documented baseline.

```sql+postgres
select
  o.name as organization,
  s.name as setting,
  s.value
from
  anthropic_compliance_organization_setting s
  join anthropic_compliance_organization o on o.uuid = s.organization_uuid
where
  s.name in ('data_retention_periods', 'content_redaction_enabled');
```

```sql+sqlite
select
  o.name as organization,
  s.name as setting,
  s.value
from
  anthropic_compliance_organization_setting s
  join anthropic_compliance_organization o on o.uuid = s.organization_uuid
where
  s.name in ('data_retention_periods', 'content_redaction_enabled');
```

### List every effective setting for one organization
Dump the full resolved configuration.

```sql+postgres
select
  name,
  type,
  value
from
  anthropic_compliance_organization_setting
where
  organization_uuid = '91012d09-e48b-438e-a489-1bebfd8fa6f9';
```

```sql+sqlite
select
  name,
  type,
  value
from
  anthropic_compliance_organization_setting
where
  organization_uuid = '91012d09-e48b-438e-a489-1bebfd8fa6f9';
```
