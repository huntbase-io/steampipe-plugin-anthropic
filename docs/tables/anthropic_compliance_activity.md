---
title: "Steampipe Table: anthropic_compliance_activity - Query Anthropic Compliance Activity Feed using SQL"
description: "Allows users to query per-event audit records for the organization from the Anthropic Compliance API."
---

# Table: anthropic_compliance_activity - Query Anthropic Compliance Activity Feed using SQL

The Activity Feed records per-event audit activity (logins, chat creation, admin changes, RBAC changes) across a Claude Enterprise organization.

## Table Usage Guide

The `anthropic_compliance_activity` table provides audit events with actor details. Accepts a Compliance Access Key (`compliance_api_key`) or an Admin API key (`admin_api_key`). Event-type-specific fields are in the `event` JSON column.

## Examples

### List recent activity
Review the latest audit events.

```sql+postgres
select
  created_at,
  type,
  actor_email,
  actor_ip_address
from
  anthropic_compliance_activity
limit 100;
```

```sql+sqlite
select
  created_at,
  type,
  actor_email,
  actor_ip_address
from
  anthropic_compliance_activity
limit 100;
```

### Count events by type
Understand which activities are most common.

```sql+postgres
select
  type,
  count(*)
from
  anthropic_compliance_activity
group by
  type
order by
  count desc;
```

```sql+sqlite
select
  type,
  count(*)
from
  anthropic_compliance_activity
group by
  type;
```

### Trace chat creation events with chat IDs
Extract event-specific fields from the raw payload.

```sql+postgres
select
  created_at,
  actor_email,
  event ->> 'claude_chat_id' as chat_id
from
  anthropic_compliance_activity
where
  type = 'claude_chat_created';
```

```sql+sqlite
select
  created_at,
  actor_email,
  event ->> 'claude_chat_id' as chat_id
from
  anthropic_compliance_activity
where
  type = 'claude_chat_created';
```

### Audit admin-surface changes
Track API key creation, RBAC changes, and SSO configuration edits.

```sql+postgres
select
  created_at,
  type,
  actor_email,
  actor_ip_address
from
  anthropic_compliance_activity
where
  type like 'admin_api_key%'
  or type like 'rbac_%'
  or type like 'org_sso%'
order by
  created_at desc;
```

```sql+sqlite
select
  created_at,
  type,
  actor_email,
  actor_ip_address
from
  anthropic_compliance_activity
where
  type like 'admin_api_key%'
  or type like 'rbac_%'
  or type like 'org_sso%'
order by
  created_at desc;
```

### Find users active from many distinct IP addresses
A high distinct-IP count can indicate session sharing or credential theft.

```sql+postgres
select
  actor_email,
  count(distinct actor_ip_address) as distinct_ips
from
  anthropic_compliance_activity
where
  actor_email is not null
group by
  actor_email
having
  count(distinct actor_ip_address) > 3
order by
  distinct_ips desc;
```

```sql+sqlite
select
  actor_email,
  count(distinct actor_ip_address) as distinct_ips
from
  anthropic_compliance_activity
where
  actor_email is not null
group by
  actor_email
having
  count(distinct actor_ip_address) > 3
order by
  distinct_ips desc;
```
